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

func (svc *Service) scheduledTask() {
	defer time.AfterFunc(time.Duration(1)*time.Second, svc.scheduledTask)

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
			if err != nil || resBlock.Result.(string) == "" {
				oo.LogW("QueryLatestBlock failed. err: %v\n", err)
				continue
			}
			latestBlockNum = utils.Hex2Dec(resBlock.Result.(string))

			latestBlockNum = int(math.Min(float64(latestBlockNum-svc.appConfig.DelayedBlockNumber), float64(currentBlockNum+svc.appConfig.BlockNumberPerReq)))
			for ; currentBlockNum <= latestBlockNum; currentBlockNum++ {

				var blockData = make([]map[string]interface{}, 0)
				currentBlock := fmt.Sprintf("0x%x", currentBlockNum)
				res, errCb := utils.ScanBlock(currentBlock, url)
				if errCb != nil {
					oo.LogW("ScanBlock failed. currentBlock id: %d. chainId:%d. err: %v\n", currentBlockNum, chainId, errCb)
					return
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
										return
									}

									resFrom, errFrom := utils.GetTransactionByHashFrom(res.Result[i].TransactionHash, url)
									if errFrom != nil {
										oo.LogW("GetTransactionByHashFrom failed. currentBlock id: %d. chainId:%s. err: %v\n", currentBlockNum, chainId, errFrom)
										return
									}

									var b = make(map[string]interface{})
									b["message_sender"] = resFrom.Result.From
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
				save(blockData, currentBlockNum, chainId, url)
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

func save(blockData []map[string]interface{}, currentBlockNum, chainId int, url string) {
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
			tokenAddress := utils.FixTo0x40String(blockData[i]["data"].(string))
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

			errTx = ownTokensImgSave(blockData[i]["address"].(string), tokenAddress, url, chainId)
			if errTx != nil {
				oo.LogW("ownTokensImgSave err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateDao {
			daoAddress := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			var addEvent = make([]map[string]interface{}, 0)
			// Add CreateDao Event Type
			eventType := []string{consts.EvCreateProposal, consts.EvVote, consts.EvCancelProposal, consts.EvAdmin, consts.EvSetting, consts.EvOwnershipTransferred, consts.EvCreateAirdrop, consts.EvClaimed}
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
			creatorAddress := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			tokenAddress := utils.FixTo0x40String(blockData[i]["data"].(string)[66:130])
			tokenChainId, _ := utils.Hex2Int64(blockData[i]["data"].(string)[:66])
			sqlInsDao := fmt.Sprintf(`INSERT INTO %s (dao_logo,dao_name,dao_address,creator,handle,description,chain_id,token_chain_id,token_address,proposal_threshold,voting_quorum,voting_period,voting_type,twitter,github,discord,update_bool) VALUES ('%s','%s','%s','%s','%s','%s',%d,%d,'%s',%d,%d,%d,'%s','%s','%s','%s',%t)`,
				consts.TbNameDao,
				"",
				"",
				daoAddress,
				creatorAddress,
				"",
				"",
				chainId,
				tokenChainId,
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
			sqlInsAdmin := fmt.Sprintf(`INSERT INTO %s (dao_address,chain_id,account,account_level) VALUES ('%s',%d,'%s','%s')`,
				consts.TbNameAdmin,
				daoAddress,
				chainId,
				creatorAddress,
				consts.LevelSuperAdmin,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsAdmin)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//save member
			sqlInsMember := fmt.Sprintf(`INSERT INTO %s (dao_address,chain_id,account,join_switch) VALUES ('%s',%d,'%s',%d)`,
				consts.TbNameMember,
				daoAddress,
				chainId,
				creatorAddress,
				1,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsMember)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvSetting {
			sqlUp := fmt.Sprintf(`UPDATE %s SET update_bool=%t WHERE dao_address='%s' AND chain_id=%d`,
				consts.TbNameDao,
				true,
				blockData[i]["address"].(string),
				chainId,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateProposal {
			proposer := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			daoAddress := blockData[i]["address"].(string)
			proposalId := utils.Hex2Dec(blockData[i]["topic1"].(string))
			startTime := utils.Hex2Dec(blockData[i]["data"].(string)[66:130])
			endTime := utils.Hex2Dec(blockData[i]["data"].(string)[130:194])
			sqlIns := fmt.Sprintf(`INSERT INTO %s (account,nonce,chain_id) VALUES ('%s',%d,%d) ON DUPLICATE KEY UPDATE nonce=nonce+1`,
				consts.TbNameNonce,
				proposer,
				1,
				chainId,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = daoAddress
			v["proposal_id"] = proposalId
			v["proposer"] = proposer
			v["start_time"] = startTime
			v["end_time"] = endTime
			m = append(m, v)
			sqlIns = oo.NewSqler().Table(consts.TbNameProposal).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//for notification
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = daoAddress
			values["types"] = consts.TypesNameProposal
			values["activity_id"] = proposalId
			values["token_address"] = ""
			values["dao_logo"] = ""
			values["notification_time"] = startTime
			values["update_bool"] = 1
			notificationData = append(notificationData, values)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

		}

		if blockData[i]["event_type"] == consts.EvCancelProposal {
			proposalId := utils.Hex2Dec(blockData[i]["topic1"].(string))
			endTime := utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			daoAddress := blockData[i]["address"].(string)
			sqlUp := fmt.Sprintf(`UPDATE %s SET end_time=%d WHERE proposal_id=%d AND chain_id=%d AND dao_address='%s'`,
				consts.TbNameProposal,
				endTime,
				proposalId,
				chainId,
				daoAddress,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvVote {
			voter := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			nonce := utils.Hex2Dec(blockData[i]["data"].(string)[66:130])
			sqlUpdate := fmt.Sprintf(`INSERT INTO %s (account,nonce,chain_id) VALUES ('%s',%d,%d) ON DUPLICATE KEY UPDATE nonce=%d`,
				consts.TbNameNonce,
				voter,
				nonce+1,
				chainId,
				nonce+1,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUpdate)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			amount, _ := utils.Hex2BigInt(blockData[i]["data"].(string)[:66])
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = blockData[i]["address"].(string)
			v["proposal_id"] = utils.Hex2Dec(blockData[i]["topic1"].(string))
			v["voter"] = voter
			v["option_index"] = utils.Hex2Dec(blockData[i]["topic3"].(string))
			v["amount"] = amount.String()
			v["nonce"] = nonce
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameVote).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvAdmin {
			daoAddress := blockData[i]["address"].(string)
			account := utils.FixTo0x40String(blockData[i]["topic1"].(string))
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
			previousOwner := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			newOwner := utils.FixTo0x40String(blockData[i]["topic2"].(string))
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

		if blockData[i]["event_type"] == consts.EvCreateAirdrop {
			amount, _ := utils.Hex2BigInt(fmt.Sprintf("0x%s", blockData[i]["data"].(string)[66:130]))
			daoAddress := blockData[i]["address"].(string)
			activityId := utils.Hex2Dec(blockData[i]["topic2"].(string))
			tokenAddress := utils.FixTo0x40String(blockData[i]["data"].(string)[2:66])
			startTime := utils.Hex2Dec(blockData[i]["data"].(string)[194:258])
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["types"] = consts.TypesNameAirdrop
			v["chain_id"] = chainId
			v["dao_address"] = daoAddress
			v["creator"] = utils.FixTo0x40String(blockData[i]["topic1"].(string))
			v["activity_id"] = activityId
			v["token_address"] = tokenAddress
			v["amount"] = amount.String()
			v["merkle_root"] = blockData[i]["data"].(string)[130:194]
			v["start_time"] = startTime
			v["end_time"] = utils.Hex2Dec(blockData[i]["data"].(string)[258:322])
			v["publish_time"] = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			v["price"] = ""
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameActivity).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//for notification
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = daoAddress
			values["types"] = consts.TypesNameAirdrop
			values["activity_id"] = activityId
			values["token_address"] = tokenAddress
			values["dao_logo"] = ""
			values["notification_time"] = startTime
			values["update_bool"] = 1
			notificationData = append(notificationData, values)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvClaimed {
			amount, _ := utils.Hex2BigInt(fmt.Sprintf("0x%s", blockData[i]["data"].(string)[130:194]))
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = blockData[i]["address"].(string)
			v["airdrop_id"] = utils.Hex2Dec(blockData[i]["topic1"].(string))
			v["index_id"] = utils.Hex2Dec(blockData[i]["data"].(string)[:66])
			v["account"] = utils.FixTo0x40String(blockData[i]["data"].(string)[66:130])
			v["amount"] = amount.String()
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameClaimed).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvTransfer {
			tokenAddress := blockData[i]["address"].(string)
			from := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			to := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			amount, _ := utils.Hex2BigInt(blockData[i]["data"].(string))

			var entityTo []models.HolderDataModel
			sqlTo := oo.NewSqler().Table(consts.TbNameHolderData).
				Where("token_address", tokenAddress).
				Where("holder_address", to).
				Where("chain_id", chainId).Select()
			errTx := oo.SqlSelect(sqlTo, &entityTo)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
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
				errTx = oo.SqlSelect(sqlFrom, &entityFrom)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
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
				errTx = oo.SqlSelect(sqlFrom, &entityZero)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
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
