package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	_ "golang.org/x/net/bpf"
	"math"
	"math/big"
	"stp_dao_v2/consts"
	_ "stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
	_ "time"
)

func (svc *Service) DoScheduledTask() {
	defer time.AfterFunc(time.Duration(1)*time.Second, svc.DoScheduledTask)
	svc.scheduledTask()
}

func (svc *Service) scheduledTask() {

	for indexScan := range svc.scanInfo {
		for indexUrl := range svc.scanInfo[indexScan].ChainId {

			url := svc.scanInfo[indexScan].ScanUrl[indexUrl]
			chainId := svc.scanInfo[indexScan].ChainId[indexUrl]

			needEvent, currentBlockNum, haveFirstBlock, errNeed := needSaveEvent(chainId)
			if errNeed != nil {
				oo.LogW("needSaveEvent failed. err:%v chainId:%d", errNeed, chainId)
				continue
			}
			if !haveFirstBlock || currentBlockNum == 0 {
				oo.LogD("needSaveEvent failed. chainId:%d", chainId)
				continue
			}

			var latestBlockNum int
			resBlock, err := utils.QueryLatestBlock(url)
			if err != nil {
				oo.LogW("QueryLatestBlock failed. err: %v\n", err)
				continue
			}
			latestBlockNum = utils.Hex2Dec(resBlock.Result)

			latestBlockNum = int(math.Min(float64(latestBlockNum-svc.appConfig.DelayedBlockNumber), float64(currentBlockNum+svc.appConfig.BlockNumberPerReq)))
			for ; currentBlockNum <= latestBlockNum; currentBlockNum++ {

				var blockData = make([]map[string]interface{}, 0)
				currentBlock := fmt.Sprintf("0x%x", currentBlockNum)
				res, errCb := utils.ScanBlock(currentBlock, url)
				if errCb != nil {
					oo.LogW("ScanBlock failed. currentBlock id: %d. chainId:%d. err: %v\n", currentBlockNum, chainId, errCb)
					break
				}

				if len(res.Result) != 0 {
					for i := range res.Result {
						if !res.Result[i].Removed {
							var topic0, topic1, topic2, topic3 string
							if res.Result[i].Topics != nil && len(res.Result[i].Topics) == 1 {
								topic0 = res.Result[i].Topics[0]
								topic1 = "0x"
								topic2 = "0x"
								topic3 = "0x"
							} else if res.Result[i].Topics != nil && len(res.Result[i].Topics) == 2 {
								topic0 = res.Result[i].Topics[0]
								topic1 = res.Result[i].Topics[1]
								topic2 = "0x"
								topic3 = "0x"
							} else if res.Result[i].Topics != nil && len(res.Result[i].Topics) == 3 {
								topic0 = res.Result[i].Topics[0]
								topic1 = res.Result[i].Topics[1]
								topic2 = res.Result[i].Topics[2]
								topic3 = "0x"
							} else if res.Result[i].Topics != nil && len(res.Result[i].Topics) == 4 {
								topic0 = res.Result[i].Topics[0]
								topic1 = res.Result[i].Topics[1]
								topic2 = res.Result[i].Topics[2]
								topic3 = res.Result[i].Topics[3]
							}
							eventType := consts.EventTypes(strings.TrimPrefix(topic0, "0x"))

							for indexNeed := range needEvent {
								if eventType == needEvent[indexNeed].EventType &&
									res.Result[i].Address == strings.ToLower(needEvent[indexNeed].Address) &&
									needEvent[indexNeed].LastBlockNumber <= currentBlockNum {

									resTime, errTime := utils.QueryTimesTamp(currentBlock, url)
									if errTime != nil {
										oo.LogW("QueryTimesTamp failed. currentBlock id: %d. chainId:%s. err: %v\n", currentBlockNum, chainId, errTime)
									}

									var b = make(map[string]interface{})
									b["event_type"] = eventType
									b["address"] = res.Result[i].Address
									b["topic0"] = topic0
									b["topic1"] = topic1
									b["topic2"] = topic2
									b["topic3"] = topic3
									b["data"] = res.Result[i].Data
									b["block_number"] = res.Result[i].BlockNumber
									b["time_stamp"] = resTime.Result.Timestamp
									b["gas_price"] = "0x"
									b["gas_used"] = resTime.Result.GasUsed
									b["log_index"] = res.Result[i].LogIndex
									b["transaction_hash"] = res.Result[i].TransactionHash
									b["transaction_index"] = res.Result[i].TransactionIndex
									b["chain_id"] = chainId

									blockData = append(blockData, b)
									break
								}
							}
						}
					}
				}
				save(blockData, currentBlockNum, chainId)
			}
		}
	}
}

func needSaveEvent(chainId int) ([]models.ScanTaskModel, int, bool, error) {
	var (
		err       error
		needEvent []models.ScanTaskModel
		min       = consts.MaxValue
	)
	sqler := oo.NewSqler().Table(consts.TbNameScanTask).Where("chain_id", chainId).Select()
	err = oo.SqlSelect(sqler, &needEvent)
	if err != nil {
		return nil, 0, false, err
	}
	if needEvent == nil || len(needEvent) == 0 {
		return nil, 0, false, nil
	}

	for i := range needEvent {
		min = int(math.Min(float64(needEvent[i].LastBlockNumber), float64(min)))
	}

	return needEvent, min, true, nil
}

func save(blockData []map[string]interface{}, currentBlockNum, chainId int) {
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	if len(blockData) != 0 {
		sqlIns := oo.NewSqler().Table(consts.TbNameEventHistorical).InsertBatch(blockData)
		_, errTx = oo.SqlxTxExec(tx, sqlIns)
		if errTx != nil {
			oo.LogW("SQL err: %v", errTx)
			return
		}
	}
	sqlUp := fmt.Sprintf(`UPDATE %s SET last_block_number=%d WHERE chain_id=%d AND last_block_number=%d`,
		consts.TbNameScanTask,
		currentBlockNum+1,
		chainId,
		currentBlockNum,
	)
	_, errTx = oo.SqlxTxExec(tx, sqlUp)
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
		return
	}

	for i := range blockData {
		if blockData[i]["event_type"] == consts.EvCreateERC20 {
			tokenAddress := utils.FixTo0xString(blockData[i]["data"].(string))
			var addEvent = make([]map[string]interface{}, 0)
			// Add CreateERC20 Event Type
			eventType := []string{consts.EvTransfer}
			for eventIndex := range eventType {
				var event = make(map[string]interface{})
				event["event_type"] = eventType[eventIndex]
				event["address"] = tokenAddress
				event["last_block_number"] = currentBlockNum
				event["rest_parameter"] = "0x"
				event["chain_id"] = chainId
				addEvent = append(addEvent, event)
			}
			sqlIns := oo.NewSqler().Table(consts.TbNameScanTask).InsertBatch(addEvent)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateDao {
			daoAddress := utils.FixTo0xString(blockData[i]["topic2"].(string))
			var addEvent = make([]map[string]interface{}, 0)
			// Add CreateDao Event Type
			eventType := []string{consts.EvCreateProposal, consts.EvVote, consts.EvCancelProposal, consts.EvAdmin, consts.EvSetting, consts.EvOwnershipTransferred}
			for eventIndex := range eventType {
				var event = make(map[string]interface{})
				event["event_type"] = eventType[eventIndex]
				event["address"] = daoAddress
				event["last_block_number"] = currentBlockNum
				event["rest_parameter"] = "0x"
				event["chain_id"] = chainId
				addEvent = append(addEvent, event)
			}
			sqlIns := oo.NewSqler().Table(consts.TbNameScanTask).InsertBatch(addEvent)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//save dao
			creatorAddress := utils.FixTo0xString(blockData[i]["topic1"].(string))
			tokenAddress := utils.FixTo0xString(blockData[i]["data"].(string)[66:130])
			sqlInsDao := fmt.Sprintf(`INSERT INTO %s (dao_logo,dao_name,dao_address,creator,handle,description,chain_id,token_address,proposal_threshold,voting_quorum,voting_period,voting_type,twitter,github,discord,update_bool) VALUES ('%s','%s','%s','%s','%s','%s',%d,'%s',%d,%d,%d,'%s','%s','%s','%s',%t)`,
				consts.TbNameDao,
				"",
				"",
				daoAddress,
				creatorAddress,
				"",
				"",
				chainId,
				tokenAddress,
				0,
				0,
				0,
				"",
				"",
				"",
				"",
				false,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsDao)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//save superAdmin
			sqlInsMember := fmt.Sprintf(`INSERT INTO %s (dao_address,chain_id,account,account_level) VALUES ('%s',%d,'%s','%s')`,
				consts.TbNameAdmin,
				daoAddress,
				chainId,
				creatorAddress,
				consts.LevelSuperAdmin,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsMember)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

		}

		if blockData[i]["event_type"] == consts.EvSetting {
			sqlUP := fmt.Sprintf(`UPDATE %s SET update_bool=%t WHERE dao_address='%s'`,
				consts.TbNameDao,
				true,
				blockData[i]["address"].(string),
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUP)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateProposal {
			proposer := utils.FixTo0xString(blockData[i]["topic2"].(string))
			sqlIns := fmt.Sprintf(`INSERT INTO %s (account,nonce) VALUES ('%s',%d) ON DUPLICATE KEY UPDATE nonce=nonce+1`,
				consts.TbNameNonce,
				proposer,
				1,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvAdmin {
			daoAddress := blockData[i]["address"].(string)
			account := utils.FixTo0xString(blockData[i]["topic1"].(string))
			enable := utils.Hex2Dec(blockData[i]["data"].(string))
			var accountLevel string
			if enable == 0 {
				accountLevel = consts.LevelNoRole
			} else if enable == 1 {
				accountLevel = consts.LevelAdmin
			}
			sqlIns := fmt.Sprintf(`REPLACE INTO %s (dao_address,chain_id,account,account_level) VALUES ('%s',%d,'%s','%s')`,
				consts.TbNameAdmin,
				daoAddress,
				chainId,
				account,
				accountLevel,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvOwnershipTransferred {
			daoAddress := blockData[i]["address"].(string)
			previousOwner := utils.FixTo0xString(blockData[i]["topic1"].(string))
			newOwner := utils.FixTo0xString(blockData[i]["topic2"].(string))
			sqlUpSuperAdmin := fmt.Sprintf(`UPDATE %s SET account='%s' WHERE dao_address='%s' AND chain_id=%d AND account='%s'`,
				consts.TbNameAdmin,
				newOwner,
				daoAddress,
				chainId,
				previousOwner,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUpSuperAdmin)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			sqlUpDaoCreator := fmt.Sprintf(`UPDATE %s SET creator='%s' WHERE dao_address='%s' AND chain_id=%d AND creator='%s'`,
				consts.TbNameDao,
				newOwner,
				daoAddress,
				chainId,
				previousOwner,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUpDaoCreator)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvTransfer {
			tokenAddress := blockData[i]["address"].(string)
			from := utils.FixTo0xString(blockData[i]["topic1"].(string))
			to := utils.FixTo0xString(blockData[i]["topic2"].(string))
			amount, _ := utils.Hex2BigInt(blockData[i]["data"].(string))

			var entityTo []models.HolderDataModel
			sqlTo := oo.NewSqler().Table(consts.TbNameHolderData).
				Where("token_address", tokenAddress).
				Where("holder_address", to).
				Where("chain_id", chainId).Select()
			err := oo.SqlSelect(sqlTo, &entityTo)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}
			var toBaseAmount = new(big.Int)
			if len(entityTo) == 0 {
				toBaseAmount, _ = utils.Dec2BigInt("0")
			} else {
				toBaseAmount, _ = utils.Dec2BigInt(entityTo[0].Balance)
			}
			amount.Add(amount, toBaseAmount)
			sqlInsTo := fmt.Sprintf(`REPLACE INTO %s (token_address,holder_address,balance,chain_id) VALUES ('%s','%s','%s',%d)`,
				consts.TbNameHolderData,
				tokenAddress,
				to,
				amount.String(),
				chainId,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsTo)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			if blockData[i]["topic1"].(string) != consts.ZeroAddress0x64 {
				var entityFrom []models.HolderDataModel
				sqlFrom := oo.NewSqler().Table(consts.TbNameHolderData).
					Where("token_address", tokenAddress).
					Where("holder_address", from).
					Where("chain_id", chainId).Select()
				err = oo.SqlSelect(sqlFrom, &entityFrom)
				if err != nil {
					oo.LogW("SQL err: %v", err)
					return
				}
				fromBaseAmount, _ := utils.Dec2BigInt(entityFrom[0].Balance)
				amount.Sub(fromBaseAmount, amount)
				sqlInsFrom := fmt.Sprintf(`UPDATE %s SET balance='%s' WHERE token_address='%s' AND holder_address='%s' AND chain_id=%d`,
					consts.TbNameHolderData,
					amount.String(),
					tokenAddress,
					from,
					chainId,
				)
				_, errTx = oo.SqlxTxExec(tx, sqlInsFrom)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
					return
				}
			}

			if blockData[i]["topic1"].(string) == consts.ZeroAddress0x64 {
				var entityZero []models.HolderDataModel
				sqlFrom := oo.NewSqler().Table(consts.TbNameHolderData).
					Where("token_address", tokenAddress).
					Where("holder_address", consts.ZeroAddress0x40).
					Where("chain_id", chainId).Select()
				err = oo.SqlSelect(sqlFrom, &entityZero)
				if err != nil {
					oo.LogW("SQL err: %v", err)
					return
				}
				var zeroBaseAmount = new(big.Int)
				if len(entityZero) == 0 {
					zeroBaseAmount, _ = utils.Dec2BigInt("0")
				} else {
					zeroBaseAmount, _ = utils.Dec2BigInt(entityZero[0].Balance)
				}
				amount.Add(amount, zeroBaseAmount)
				sqlInsZero := fmt.Sprintf(`REPLACE INTO %s (token_address,holder_address,balance,chain_id) VALUES ('%s','%s','%s',%d)`,
					consts.TbNameHolderData,
					tokenAddress,
					consts.ZeroAddress0x40,
					amount.String(),
					chainId,
				)
				_, errTx = oo.SqlxTxExec(tx, sqlInsZero)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
					return
				}
			}
		}
	}
}

func (svc *Service) DoUpdateDaoInfoTask() {
	defer time.AfterFunc(time.Duration(60)*time.Second, svc.DoUpdateDaoInfoTask)
	svc.updateDaoInfoTask()
}

func (svc *Service) updateDaoInfoTask() {
	const data = "0xafe926e8"

	var entities []models.DaoModel
	sqler := oo.NewSqler().Table(consts.TbNameDao).Where("update_bool", 1).Select()
	err := oo.SqlSelect(sqler, &entities)
	if err != nil {
		oo.LogW("query Dao failed. err:%v", err)
		return
	}

	if len(entities) != 0 {
		for index := range entities {

			for indexScan := range svc.scanInfo {
				for indexUrl := range svc.scanInfo[indexScan].ChainId {

					url := svc.scanInfo[indexScan].ScanUrl[indexUrl]
					chainId := svc.scanInfo[indexScan].ChainId[indexUrl]

					if chainId == entities[index].ChainId {
						res, errQ := utils.QueryDaoInfo(entities[index].DaoAddress, data, url)
						if errQ != nil {
							oo.LogW("QueryDaoInfo failed. chainId:%d. err: %v\n", chainId, errQ)
							return
						}

						if len(res.Result) != 0 {
							var outputParameters []string
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")

							daoInfo, errDe := utils.Decode(outputParameters, strings.TrimPrefix(res.Result, "0x"))
							if errDe != nil {
								oo.LogW("Decode failed. chainId:%d. err: %v\n", chainId, errDe)
								return
							}
							saveDaoInfoAndCategory(daoInfo, entities[index].DaoAddress, chainId)
						}
					}
				}
			}
		}
	}
}

func saveDaoInfoAndCategory(daoInfo []interface{}, daoAddress string, chainId int) {
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	daoName := daoInfo[0]
	handle := daoInfo[1]
	category := daoInfo[2]
	description := daoInfo[3]
	twitter := daoInfo[4]
	github := daoInfo[5]
	discord := daoInfo[6]
	daoLogo := daoInfo[7]

	sqlIns := fmt.Sprintf(`UPDATE %s SET dao_logo='%s',dao_name='%s',handle='%s',description='%s',twitter='%s',github='%s',discord='%s',update_bool=%t WHERE dao_address='%s' AND chain_id=%d`,
		consts.TbNameDao,
		daoLogo,
		daoName,
		handle,
		description,
		twitter,
		github,
		discord,
		false,
		daoAddress,
		chainId,
	)
	_, err := oo.SqlxTxExec(tx, sqlIns)
	if err != nil {
		oo.LogW("SQL failed. err: %v\n", err)
		return
	}

	var daoId int
	sqlSelDId := oo.NewSqler().Table(consts.TbNameDao).Where("dao_address", daoAddress).Where("chain_id", chainId).Select("id")
	err = oo.SqlGet(sqlSelDId, &daoId)
	if err != nil {
		oo.LogW("SQL failed. err: %v\n", err)
		return
	}

	categorySplit := strings.Split(category.(string), ",")
	for _, categoryName := range categorySplit {
		if categoryName == "" {
			continue
		}
		sqlUP := fmt.Sprintf(`INSERT INTO %s (category_name) VALUES ('%s') ON DUPLICATE KEY UPDATE category_name=category_name`,
			consts.TbNameCategory,
			categoryName,
		)
		err = oo.SqlExec(sqlUP)
		if err != nil {
			oo.LogW("SQL failed. err: %v\n", err)
			return
		}

		var categoryId int
		sqlSelCId := oo.NewSqler().Table(consts.TbNameCategory).Where("category_name", categoryName).Select("id")
		err = oo.SqlGet(sqlSelCId, &categoryId)
		if err != nil {
			oo.LogW("SQL failed. err: %v\n", err)
			return
		}

		sqlInsCategory := fmt.Sprintf(`INSERT INTO %s (dao_id,category_id) VALUES (%d,%d)`,
			consts.TbNameDaoCategory,
			daoId,
			categoryId,
		)
		_, err = oo.SqlxTxExec(tx, sqlInsCategory)
		if err != nil {
			oo.LogW("SQL failed. err: %v\n", err)
			return
		}
	}
}
