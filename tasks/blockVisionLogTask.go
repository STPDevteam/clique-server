package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/spf13/viper"
	_ "golang.org/x/net/bpf"
	"math"
	"math/big"
	"stp_dao_v2/consts"
	_ "stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/utils"
	"strconv"
	"strings"
	"time"
	_ "time"
)

func ScheduledTask() {
	defer time.AfterFunc(time.Duration(1)*time.Second, ScheduledTask)

	for indexUrl := range viper.GetIntSlice("scan.chain_id") {

		url := viper.GetStringSlice("scan.scan_url")[indexUrl]
		chainId := viper.GetIntSlice("scan.chain_id")[indexUrl]
		delayedBlockNum := viper.GetIntSlice("scan.delayed_block_number")[indexUrl]
		blockNumberPerReq := viper.GetInt64("scan.block_number_per_req")

		needEvent, currentBlockNum, haveFirstBlock, errNeed := needSaveEvent(chainId)
		if errNeed != nil {
			//oo.LogW("needSaveEvent failed. err:%v chainId:%d", errNeed, chainId)
			continue
		}
		if !haveFirstBlock || currentBlockNum == 0 {
			//oo.LogD("needSaveEvent failed. chainId:%d", chainId)
			continue
		}

		resBlock, err1 := utils.QueryLatestBlock(url)
		if err1 != nil {
			oo.LogW("utils.QueryLatestBlock failed. err1: %v, chainId: %d", err1, chainId)
			continue
		}
		s, ok := resBlock.Result.(string)
		if !ok {
			continue
		}
		latestBlockNum, err2 := utils.Hex2Int64(s)
		if err2 != nil {
			oo.LogW("utils.QueryLatestBlock failed. err2: %v, chainId: %d", err2, chainId)
			continue
		}

		latestBlockNum = int64(math.Min(float64(latestBlockNum-int64(delayedBlockNum)), float64(currentBlockNum+blockNumberPerReq)))
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
						eventType := EventTypes(strings.TrimPrefix(topic0, "0x"))

						for indexNeed := range needEvent {
							if eventType == needEvent[indexNeed].EventType &&
								strings.ToLower(res.Result[i].Address) == strings.ToLower(needEvent[indexNeed].Address) &&
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

								var eventCount int
								sqlSel := oo.NewSqler().Table(consts.TbNameEventHistorical).
									Where("chain_id", chainId).Where("transaction_hash", res.Result[i].TransactionHash).
									Where("log_index", res.Result[i].LogIndex).Count()
								err := oo.SqlGet(sqlSel, &eventCount)
								if err != nil {
									oo.LogW("query event_historical_data SQL err: %v", err)
									return
								}
								if eventCount != 0 {
									oo.LogW("event already exists")
									continue
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

func needSaveEvent(chainId int) ([]db.TbScanTaskModel, int64, bool, error) {
	var (
		err       error
		needEvent []db.TbScanTaskModel
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
		min = int(int64(math.Min(float64(needEvent[i].LastBlockNumber), float64(min))))
	}

	return needEvent, int64(min), true, nil
}

func save(blockData []map[string]interface{}, currentBlockNum int64, chainId int, url string) {
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer func() {
		oo.CloseSqlxTx(tx, &errTx)
		sqlDel := oo.NewSqler().Table(consts.TbNameHandleLock).
			Where("lock_block", "<", currentBlockNum).
			Where("chain_id", chainId).Delete()
		err := oo.SqlExec(sqlDel)
		if err != nil {
			oo.LogW("SQL err: %v", err)
		}
	}()

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

			errTx = ownTokensImgSave(blockData[i]["address"].(string), tokenAddress, url, chainId, tx)
			if errTx != nil {
				oo.LogW("ownTokensImgSave func err: %v", errTx)
				return
			}

			//for notification
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = tokenAddress
			values["types"] = consts.TypesNameReserveToken
			values["activity_id"] = 0
			values["dao_logo"] = ""
			values["dao_name"] = ""
			values["activity_name"] = ""
			values["start_time"] = 0
			values["update_bool"] = 0
			notificationData = append(notificationData, values)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			var result sql.Result
			result, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			res, err := utils.ReserveToken(blockData[i]["transaction_hash"].(string), url)
			if err != nil {
				oo.LogW("ReserveToken failed. err: %v\n", err)
				errTx = err
				return
			}
			resTokenAccount, err := utils.DecodeDistribution(res.Result.Input)
			if err != nil {
				oo.LogW("DecodeDistribution failed. err: %v\n", err)
				errTx = err
				return
			}

			var m = make([]map[string]interface{}, 0)
			for index := range resTokenAccount {
				var v = make(map[string]interface{})
				v["notification_id"], _ = result.LastInsertId()
				v["account"] = resTokenAccount[index].Recipient
				v["already_read"] = 0
				v["notification_time"] = resTokenAccount[index].LockDate
				m = append(m, v)
			}
			sqlIns = oo.NewSqler().Table(consts.TbNameNotificationAccount).InsertBatch(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v\n", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = blockData[i]["message_sender"].(string)
			record["types"] = consts.EvCreateERC20
			record["chain_id"] = chainId
			record["address"] = tokenAddress
			record["activity_id"] = 0
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvClaimReserve {
			// save account record
			account := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			tokenAddress := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = account
			record["types"] = consts.EvClaimReserve
			record["chain_id"] = chainId
			record["address"] = tokenAddress
			record["activity_id"] = 0
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateDao {
			daoAddress := utils.FixTo0x40String(blockData[i]["topic3"].(string))
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

			// save dao
			creatorAddress := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			tokenAddress := utils.FixTo0x40String(blockData[i]["data"].(string)[66:130])
			tokenChainId, _ := utils.Hex2Int64(blockData[i]["data"].(string)[:66])

			var daoMap = make([]map[string]interface{}, 0)
			var daoValues = make(map[string]interface{})
			daoValues["dao_logo"] = ""
			daoValues["dao_name"] = ""
			daoValues["dao_address"] = daoAddress
			daoValues["creator"] = creatorAddress
			daoValues["handle"] = ""
			daoValues["description"] = ""
			daoValues["chain_id"] = chainId
			daoValues["token_chain_id"] = tokenChainId
			daoValues["token_address"] = tokenAddress
			daoValues["proposal_threshold"] = 0
			daoValues["voting_quorum"] = 0
			daoValues["voting_period"] = 0
			daoValues["voting_type"] = ""
			daoValues["twitter"] = ""
			daoValues["github"] = ""
			daoValues["discord"] = ""
			daoValues["website"] = ""
			daoValues["update_bool"] = 0
			daoValues["approve"] = 0 //approve 0
			daoMap = append(daoMap, daoValues)
			sqlInsDao := oo.NewSqler().Table(consts.TbNameDao).Insert(daoMap)
			_, errTx = oo.SqlxTxExec(tx, sqlInsDao)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save superAdmin
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

			// team space jobs
			var mJob = make([]map[string]interface{}, 0)
			var vJob = make(map[string]interface{})
			vJob["chain_id"] = chainId
			vJob["dao_address"] = daoAddress
			vJob["account"] = creatorAddress
			vJob["job"] = consts.Jobs_A_superAdmin
			mJob = append(mJob, vJob)
			_, errTx = o.InsertTx(tx, consts.TbJobs, mJob)
			if errTx != nil {
				oo.LogW("SQL err:%v", errTx)
				return
			}

			var mGeneral = make([]map[string]interface{}, 0)
			var vGeneral = make(map[string]interface{})
			vGeneral["chain_id"] = chainId
			vGeneral["dao_address"] = daoAddress
			vGeneral["creator"] = creatorAddress
			vGeneral["title"] = "General"
			vGeneral["last_edit_time"] = time.Now().Unix()
			vGeneral["last_edit_by"] = creatorAddress
			vGeneral["access"] = "public"
			mGeneral = append(mGeneral, vGeneral)
			_, errTx = o.InsertTx(tx, consts.TbTeamSpaces, mGeneral)
			if errTx != nil {
				oo.LogW("SQL err:%v", errTx)
				return
			}

			// save member
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

			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = daoAddress
			v["account"] = creatorAddress
			v["operate"] = "created"
			v["signature"] = blockData[i]["transaction_hash"].(string)
			v["message"] = ""
			v["timestamp"] = time.Now().Unix()
			m = append(m, v)
			sqlIns = oo.NewSqler().Table(consts.TbNameAccountSign).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// update handle lock
			handleKeccak := strings.TrimPrefix(blockData[i]["topic1"].(string), "0x")
			var daoHandle = make(map[string]interface{})
			daoHandle["lock_block"] = consts.MaxIntUnsigned
			sqlUpDaoHandle := oo.NewSqler().Table(consts.TbNameHandleLock).Where("handle_keccak", handleKeccak).Update(daoHandle)
			_, errTx = oo.SqlxTxExec(tx, sqlUpDaoHandle)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = creatorAddress
			record["types"] = consts.EvCreateDao
			record["chain_id"] = chainId
			record["address"] = daoAddress
			record["activity_id"] = 0
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvSetting {
			sqlUp = fmt.Sprintf(`UPDATE %s SET update_bool=%t WHERE dao_address='%s' AND chain_id=%d`,
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
			proposalId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			startTime, _ := utils.Hex2Dec(blockData[i]["data"].(string)[:66])
			endTime, _ := utils.Hex2Dec(blockData[i]["data"].(string)[66:130])

			proposalTitle, err := proposalInfo(daoAddress, blockData[i]["topic1"].(string), url)
			if err != nil {
				oo.LogW("proposalInfo func err: %v", err)
				errTx = err
				return
			}

			var daoEntity db.TbDaoModel
			sqlSel := oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", chainId).Where("dao_address", daoAddress).Select()
			errTx = oo.SqlGet(sqlSel, &daoEntity)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			var blockNumber string
			if chainId == daoEntity.TokenChainId {
				blockNumber = blockData[i]["block_number"].(string)
			} else {
				for indexUrl := range viper.GetIntSlice("scan.chain_id") {
					if viper.GetIntSlice("scan.chain_id")[indexUrl] == daoEntity.TokenChainId {
						timestamp, _ := utils.Hex2Dec(blockData[i]["time_stamp"].(string))
						if daoEntity.TokenChainId == consts.KlaytnTestnet1001 || daoEntity.TokenChainId == consts.KlaytnMainnet8217 ||
							daoEntity.TokenChainId == consts.BSCTestnet97 || daoEntity.TokenChainId == consts.BSCMainnet56 {
							var duration = 0
						Loop:
							var klaytnUrl = fmt.Sprintf(viper.GetStringSlice("scan.query_block_number_url")[indexUrl], daoEntity.TokenChainId, timestamp, duration)
							resKB, errK := utils.GetKlaytnBlock(klaytnUrl)
							if errK != nil {
								oo.LogW("GetKlaytnBlock failed err: %v url:%s tokenChainId:%v timestamp:%v duration:%v", errK, viper.GetStringSlice("scan.query_block_number_url")[indexUrl], daoEntity.TokenChainId, timestamp, duration)
								errTx = errK
								return
							}
							if duration > 1000 {
								oo.LogW("query block number failed,duration > 1000. timestamp: %v, chainId:%d", timestamp, daoEntity.TokenChainId)
								errTx = errors.New("query block number failed,duration > 1000")
								return
							}
							if len(resKB.Data) == 0 {
								duration += 100
								goto Loop
							}
							blockNumber = resKB.Data[0]
						} else {
							urlGetBlock := fmt.Sprintf(viper.GetStringSlice("scan.query_block_number_url")[indexUrl], timestamp)
							res, errG := utils.GetBlockNumberFromTimestamp(urlGetBlock)
							if errG != nil {
								oo.LogW("GetBlockNumberFromTimestamp err: %v", errTx)
								errTx = errG
								return
							}
							if res.Status != "1" || strings.ToLower(res.Message) != "ok" {
								errTx = errors.New("getBlockNumberFromTimestamp failed")
								return
							}
							blockDec, errA := strconv.Atoi(res.Result)
							if errA != nil {
								errTx = errA
								return
							}
							blockNumber = fmt.Sprintf("0x%x", blockDec)
						}
					}
				}

			}

			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = daoAddress
			v["proposal_id"] = proposalId
			v["title"] = proposalTitle[:int(math.Min(float64(len(proposalTitle)), 500))]
			v["id_v1"] = 0
			v["content_v1"] = ""
			v["proposer"] = proposer
			v["start_time"] = startTime
			v["end_time"] = endTime
			v["version"] = "v2"
			v["block_number"] = blockNumber
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameProposal).Insert(m)
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
			values["types"] = consts.TypesNameNewProposal
			values["activity_id"] = proposalId
			values["dao_logo"] = ""
			values["dao_name"] = ""
			values["activity_name"] = proposalTitle[:int(math.Min(float64(len(proposalTitle)), 500))]
			values["start_time"] = startTime
			values["update_bool"] = 1
			notificationData = append(notificationData, values)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// for dao order with proposal total
			//var totalProposal int
			//sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("deprecated", 0).
			//	Where("chain_id", chainId).Where("dao_address", daoAddress).Count()
			//errTx = oo.SqlGet(sqlSel, &totalProposal)
			//if errTx != nil {
			//	oo.LogW("SQL err: %v", errTx)
			//	return
			//}
			var weight = make(map[string]interface{})
			weight["weight"] = daoEntity.Weight.Int64 + 1
			sqlUp = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", chainId).Where("dao_address", daoAddress).Update(weight)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v ", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = proposer
			record["types"] = consts.EvCreateProposal
			record["chain_id"] = chainId
			record["address"] = daoAddress
			record["activity_id"] = proposalId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCancelProposal {
			proposalId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			endTime, _ := utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			daoAddress := blockData[i]["address"].(string)
			sqlUp = fmt.Sprintf(`UPDATE %s SET end_time=%d WHERE proposal_id=%d AND chain_id=%d AND dao_address='%s'`,
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

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = blockData[i]["message_sender"].(string)
			record["types"] = consts.EvCancelProposal
			record["chain_id"] = chainId
			record["address"] = daoAddress
			record["activity_id"] = proposalId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// total proposals
			dao, err := db.GetTbDao(o.W("chain_id", chainId), o.W("dao_address", daoAddress))
			if err != nil {
				oo.LogW("SQL err: %v", err)
				errTx = err
				return
			}

			var vDao = make(map[string]interface{})
			vDao["total_proposals"] = dao.TotalProposals + 1
			_, errTx = o.UpdateTx(tx, consts.TbNameDao, vDao, o.W("id", dao.Id))
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvVote {
			proposalId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			voter := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			daoAddress := blockData[i]["address"].(string)

			amount, _ := utils.Hex2BigInt(blockData[i]["data"].(string)[:66])
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["dao_address"] = daoAddress
			v["proposal_id"] = proposalId
			v["voter"] = voter
			v["option_index"], _ = utils.Hex2Dec(blockData[i]["topic3"].(string))
			v["amount"] = amount.String()
			v["nonce"] = 0
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameVote).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = voter
			record["types"] = consts.EvVote
			record["chain_id"] = chainId
			record["address"] = daoAddress
			record["activity_id"] = proposalId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvAdmin {
			daoAddress := blockData[i]["address"].(string)
			account := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			enable, _ := utils.Hex2Dec(blockData[i]["data"].(string))

			// jobs admin, first delete, if add, insert
			errTx = o.Delete(consts.TbJobs, o.W("chain_id", chainId), o.W("dao_address", daoAddress),
				o.W("account", account), o.W("job", "!=", consts.Jobs_A_superAdmin))
			if errTx != nil {
				oo.LogW("SQL err:%v", errTx)
				return
			}

			if enable == 0 {
				errTx = o.Delete(consts.TbNameAdmin, o.W("chain_id", chainId), o.W("dao_address", daoAddress),
					o.W("account", account), o.W("account_level", consts.LevelAdmin))
				if errTx != nil {
					oo.LogW("SQL err:%v", errTx)
					return
				}
			} else if enable == 1 {
				// jobs admin, first delete, if add, insert
				count, err := o.Count(consts.TbJobs, o.W("chain_id", chainId), o.W("dao_address", daoAddress),
					o.W("account", account), o.W("job", consts.Jobs_A_superAdmin))
				if err != nil {
					oo.LogW("SQL err:%v", errTx)
					errTx = err
					return
				}
				if count == 0 {
					// jobs admin
					var mJobAdmin = make([]map[string]interface{}, 0)
					var vJobAdmin = make(map[string]interface{})
					vJobAdmin["chain_id"] = chainId
					vJobAdmin["dao_address"] = daoAddress
					vJobAdmin["account"] = account
					vJobAdmin["job"] = consts.Jobs_B_admin
					mJobAdmin = append(mJobAdmin, vJobAdmin)
					_, errTx = o.InsertTx(tx, consts.TbJobs, mJobAdmin)
					if errTx != nil {
						oo.LogW("SQL err:%v", errTx)
						return
					}

					var mAdmin = make([]map[string]interface{}, 0)
					var vAdmin = make(map[string]interface{})
					vAdmin["chain_id"] = chainId
					vAdmin["dao_address"] = daoAddress
					vAdmin["account"] = account
					vAdmin["account_level"] = consts.LevelAdmin
					mAdmin = append(mAdmin, vAdmin)
					errTx = o.Insert(consts.TbNameAdmin, mAdmin)
					if errTx != nil {
						oo.LogW("SQL err:%v", errTx)
						return
					}
				}
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = blockData[i]["message_sender"].(string)
			record["types"] = consts.EvAdmin
			record["chain_id"] = chainId
			record["address"] = daoAddress
			record["activity_id"] = 0
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvOwnershipTransferred {
			daoAddress := blockData[i]["address"].(string)
			previousOwner := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			newOwner := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			sqlUpSuperAdmin := fmt.Sprintf(`UPDATE %s SET account='%s' WHERE dao_address='%s' AND chain_id=%d AND account='%s' AND account_level='%s'`,
				consts.TbNameAdmin,
				newOwner,
				daoAddress,
				chainId,
				previousOwner,
				consts.LevelSuperAdmin,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUpSuperAdmin)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// jobs superAdmin
			var vJobsSuperAdmin = make(map[string]interface{})
			vJobsSuperAdmin["account"] = newOwner
			errTx = o.Update(consts.TbJobs, vJobsSuperAdmin, o.W("chain_id", chainId),
				o.W("dao_address", daoAddress),
				o.W("account", previousOwner))
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

			// save account record
			var contractAddress string
			sqlSel := oo.NewSqler().Table(consts.TbNameScanTask).Where("event_type", consts.EvCreateDao).
				Where("chain_id", chainId).Select("address")
			errTx = oo.SqlGet(sqlSel, &contractAddress)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			if previousOwner != consts.ZeroAddress0x40 && previousOwner != strings.ToLower(contractAddress) {
				var recordInsert = make([]map[string]interface{}, 0)
				var record = make(map[string]interface{})
				record["creator"] = previousOwner
				record["types"] = consts.EvOwnershipTransferred
				record["chain_id"] = chainId
				record["address"] = daoAddress
				record["activity_id"] = 0
				record["avatar"] = ""
				record["dao_name"] = ""
				record["titles"] = ""
				record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
				record["update_bool"] = 1
				recordInsert = append(recordInsert, record)
				sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
				_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
					return
				}
			}

			// save member
			var contractAddr string
			sqlSel = oo.NewSqler().Table(consts.TbNameScanTask).Where("chain_id", chainId).Where("event_type", consts.EvCreateDao).Select("address")
			errTx = oo.SqlGet(sqlSel, &contractAddr)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			if strings.ToLower(contractAddr) != strings.ToLower(newOwner) {
				var count int
				sqlSel = oo.NewSqler().Table(consts.TbNameMember).Where("chain_id", chainId).Where("dao_address", daoAddress).
					Where("account", newOwner).Count()
				errTx = oo.SqlGet(sqlSel, &count)
				if errTx != nil {
					oo.LogW("SQL err: %v", errTx)
					return
				}
				if count == 0 {
					var m = make([]map[string]interface{}, 0)
					var v = make(map[string]interface{})
					v["chain_id"] = chainId
					v["dao_address"] = daoAddress
					v["account"] = newOwner
					v["join_switch"] = 1
					m = append(m, v)
					sqlInsMember := oo.NewSqler().Table(consts.TbNameMember).Insert(m)
					_, errTx = oo.SqlxTxExec(tx, sqlInsMember)
					if errTx != nil {
						oo.LogW("SQL err: %v", errTx)
						return
					}
				} else {
					var oldDataCount int
					sqlSel = oo.NewSqler().Table(consts.TbNameMember).Where("chain_id", chainId).Where("dao_address", daoAddress).
						Where("account", newOwner).Where("join_switch", 1).Count()
					errTx = oo.SqlGet(sqlSel, &oldDataCount)
					if errTx != nil {
						oo.LogW("SQL err: %v", errTx)
						return
					}
					if oldDataCount != 1 {
						var update = make(map[string]interface{})
						update["join_switch"] = 1
						sqlUp = oo.NewSqler().Table(consts.TbNameMember).Where("chain_id", chainId).Where("dao_address", daoAddress).
							Where("account", newOwner).Update(update)
						_, errTx = oo.SqlxTxExec(tx, sqlUp)
						if errTx != nil {
							oo.LogW("SQL err: %v", errTx)
							return
						}
					}
				}
			}
		}

		if blockData[i]["event_type"] == consts.EvCreateAirdrop {
			creator := utils.FixTo0x40String(blockData[i]["topic1"].(string))
			airdropId, _ := utils.Hex2Dec(blockData[i]["topic2"].(string))
			tokenAddress := utils.FixTo0x40String(blockData[i]["data"].(string)[2:66])
			stakingAmount, _ := utils.Hex2BigInt(fmt.Sprintf("0x%s", blockData[i]["data"].(string)[66:130]))
			airdropStartTime, _ := utils.Hex2Dec(blockData[i]["data"].(string)[130:194])
			airdropEndTime, _ := utils.Hex2Dec(blockData[i]["data"].(string)[194:258])

			var airdropEntity []db.TbAirdropModel
			sqlSel := oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", airdropId).Select()
			errTx = oo.SqlSelect(sqlSel, &airdropEntity)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			if len(airdropEntity) == 0 {
				oo.LogW("there is no airdrop data in the background")
				errTx = errors.New("there is no airdrop data in the background")
				return
			}

			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["types"] = consts.TypesNameAirdrop
			v["chain_id"] = airdropEntity[0].ChainId
			v["dao_address"] = airdropEntity[0].DaoAddress
			v["creator"] = creator
			v["activity_id"] = airdropId
			v["token_chain_id"] = airdropEntity[0].TokenChainId
			v["token_address"] = tokenAddress
			v["staking_amount"] = stakingAmount.String()
			v["airdrop_amount"] = 0
			v["merkle_root"] = ""
			v["start_time"] = airdropEntity[0].StartTime
			v["end_time"] = airdropEntity[0].EndTime
			v["airdrop_start_time"] = airdropStartTime
			v["airdrop_end_time"] = airdropEndTime
			v["publish_time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			v["price"] = ""
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameActivity).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = creator
			record["types"] = consts.EvCreateAirdrop
			record["chain_id"] = airdropEntity[0].ChainId
			record["address"] = airdropEntity[0].DaoAddress
			record["activity_id"] = airdropId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvSettleAirdrop {
			airdropId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			airdropAmount, _ := utils.Hex2BigInt(blockData[i]["data"].(string)[:66])
			merkleRoot := fmt.Sprintf("0x%s", blockData[i]["data"].(string)[66:130])

			var prepareAddress string
			sqlSel := oo.NewSqler().Table(consts.TbNameAirdropPrepare).Where("airdrop_id", airdropId).
				Where("root", merkleRoot).Select("prepare_address")
			errTx = oo.SqlGet(sqlSel, &prepareAddress)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			var v = make(map[string]interface{})
			v["airdrop_address"] = prepareAddress
			sqlUp = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", airdropId).Update(v)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var set = make(map[string]interface{})
			set["airdrop_amount"] = airdropAmount
			set["merkle_root"] = merkleRoot
			sqlUp = oo.NewSqler().Table(consts.TbNameActivity).Where("activity_id", airdropId).Update(set)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var airdropEntity []db.TbAirdropModel
			sqlSel = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", airdropId).Select()
			errTx = oo.SqlSelect(sqlSel, &airdropEntity)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			if len(airdropEntity) == 0 {
				oo.LogW("there is no airdrop data in the background")
				errTx = errors.New("there is no airdrop data in the background")
				return
			}

			//for notification
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = airdropEntity[0].ChainId
			values["dao_address"] = airdropEntity[0].DaoAddress
			values["types"] = consts.TypesNameAirdrop
			values["activity_id"] = airdropId
			values["dao_logo"] = ""
			values["dao_name"] = ""
			values["activity_name"] = airdropEntity[0].Title[:int(math.Min(float64(len(airdropEntity[0].Title)), 500))]
			values["start_time"] = airdropEntity[0].AirdropStartTime
			values["update_bool"] = 1
			notificationData = append(notificationData, values)
			sqlIns := oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = blockData[i]["message_sender"].(string)
			record["types"] = consts.EvSettleAirdrop
			record["chain_id"] = airdropEntity[0].ChainId
			record["address"] = airdropEntity[0].DaoAddress
			record["activity_id"] = airdropId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvClaimed {
			airdropId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			amount, _ := utils.Hex2BigInt(fmt.Sprintf("0x%s", blockData[i]["data"].(string)[130:194]))

			var airdropEntity []db.TbAirdropModel
			sqlSel := oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", airdropId).Select()
			errTx = oo.SqlSelect(sqlSel, &airdropEntity)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			if len(airdropEntity) == 0 {
				oo.LogW("there is no airdrop data in the background")
				errTx = errors.New("there is no airdrop data in the background")
				return
			}

			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = airdropEntity[0].ChainId
			v["dao_address"] = airdropEntity[0].DaoAddress
			v["airdrop_id"] = airdropId
			v["index_id"], _ = utils.Hex2Dec(blockData[i]["data"].(string)[:66])
			v["account"] = utils.FixTo0x40String(blockData[i]["data"].(string)[66:130])
			v["amount"] = amount.String()
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameClaimed).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			// save account record
			var recordInsert = make([]map[string]interface{}, 0)
			var record = make(map[string]interface{})
			record["creator"] = blockData[i]["message_sender"].(string)
			record["types"] = consts.EvClaimed
			record["chain_id"] = airdropEntity[0].ChainId
			record["address"] = airdropEntity[0].DaoAddress
			record["activity_id"] = airdropId
			record["avatar"] = ""
			record["dao_name"] = ""
			record["titles"] = ""
			record["time"], _ = utils.Hex2Dec(blockData[i]["time_stamp"].(string))
			record["update_bool"] = 1
			recordInsert = append(recordInsert, record)
			sqlInsRecord := oo.NewSqler().Table(consts.TbNameAccountRecord).Insert(recordInsert)
			_, errTx = oo.SqlxTxExec(tx, sqlInsRecord)
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

			var entityTo []db.TbHolderDataModel
			sqlTo := oo.NewSqler().Table(consts.TbNameHolderData).
				Where("token_address", tokenAddress).
				Where("holder_address", to).
				Where("chain_id", chainId).Select()
			errTx = oo.SqlSelect(sqlTo, &entityTo)
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
			var toAmount = new(big.Int).Add(amount, toBaseAmount)
			sqlInsTo := fmt.Sprintf(`REPLACE INTO %s (token_address,holder_address,balance,chain_id) VALUES ('%s','%s','%s',%d)`,
				consts.TbNameHolderData,
				tokenAddress,
				to,
				toAmount.String(),
				chainId,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlInsTo)
			if errTx != nil {
				oo.LogW("hash: %v", blockData[i]["transaction_hash"].(string))
				oo.LogW("sqlInsTo sqler: %v", sqlInsTo)
				oo.LogW("SQL err: %v", errTx)
				return
			}

			if blockData[i]["topic1"].(string) != consts.ZeroAddress0x64 {
				var entityFrom []db.TbHolderDataModel
				sqlFrom := oo.NewSqler().Table(consts.TbNameHolderData).
					Where("token_address", tokenAddress).
					Where("holder_address", from).
					Where("chain_id", chainId).Select()
				errTx = oo.SqlSelect(sqlFrom, &entityFrom)
				if errTx != nil {
					oo.LogW("hash: %v", blockData[i]["transaction_hash"].(string))
					oo.LogW("sqlFrom sqler: %v", sqlFrom)
					oo.LogW("SQL err: %v", errTx)
					return
				}
				fromBaseAmount, _ := utils.Dec2BigInt(entityFrom[0].Balance)
				var fromAmount = new(big.Int).Sub(fromBaseAmount, amount)
				sqlInsFrom := fmt.Sprintf(`UPDATE %s SET balance='%s' WHERE token_address='%s' AND holder_address='%s' AND chain_id=%d`,
					consts.TbNameHolderData,
					fromAmount.String(),
					tokenAddress,
					from,
					chainId,
				)
				_, errTx = oo.SqlxTxExec(tx, sqlInsFrom)
				if errTx != nil {
					oo.LogW("hash: %v", blockData[i]["transaction_hash"].(string))
					oo.LogW("sqlInsFrom sqler: %v", sqlInsFrom)
					oo.LogW("SQL err: %v", errTx)
					return
				}
			}

			if blockData[i]["topic1"].(string) == consts.ZeroAddress0x64 {
				var entityZero []db.TbHolderDataModel
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
				var zeroAmount = new(big.Int).Add(amount, zeroBaseAmount)
				sqlInsZero := fmt.Sprintf(`REPLACE INTO %s (token_address,holder_address,balance,chain_id) VALUES ('%s','%s','%s',%d)`,
					consts.TbNameHolderData,
					tokenAddress,
					consts.ZeroAddress0x40,
					zeroAmount.String(),
					chainId,
				)
				_, errTx = oo.SqlxTxExec(tx, sqlInsZero)
				if errTx != nil {
					oo.LogW("hash: %v", blockData[i]["transaction_hash"].(string))
					oo.LogW("sqlInsZero sqler: %v", sqlInsZero)
					oo.LogW("SQL err: %v", errTx)
					return
				}
			}
		}

		//CreatedSale(uint256 indexed saleId, address indexed saleToken, address indexed receiveToken, uint256 saleAmount, uint256 pricePer, uint256 limitMin, uint256 limitMax, uint256 startTime, uint256 endTime)
		if blockData[i]["event_type"] == consts.EvCreatedSale {
			saleId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			saleToken := utils.FixTo0x40String(blockData[i]["topic2"].(string))
			receiveToken := utils.FixTo0x40String(blockData[i]["topic3"].(string))

			decode, err := utils.Decode([]string{"uint256", "uint256", "uint256", "uint256", "uint256", "uint256"}, blockData[i]["data"].(string)[2:])
			if err != nil {
				oo.LogW("Decode err: %v", err)
				errTx = err
				return
			}

			var v = make(map[string]interface{})
			v["chain_id"] = chainId
			v["creator"] = blockData[i]["message_sender"]
			v["sale_token"] = saleToken
			v["sale_amount"] = decode[0]
			v["sale_price"] = decode[1]
			v["receive_token"] = receiveToken
			v["limit_min"] = decode[2]
			v["limit_max"] = decode[3]
			v["start_time"] = decode[4]
			v["end_time"] = decode[5]
			v["status"] = consts.StatusSoon
			sqlUp = oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Update(v)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//for notification
			var swapData db.TbSwap
			sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Select()
			errTx = oo.SqlGet(sqlSel, &swapData)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = saleToken
			values["types"] = consts.TypesNamePublicSaleCreated
			values["activity_id"] = swapData.Id
			values["dao_logo"] = swapData.SaleTokenImg
			values["dao_name"] = ""
			values["activity_name"] = swapData.Title
			values["start_time"] = swapData.StartTime
			values["update_bool"] = 0
			notificationData = append(notificationData, values)
			sqlIns := oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			var result sql.Result
			result, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var m = make([]map[string]interface{}, 0)
			var val = make(map[string]interface{})
			val["notification_id"], _ = result.LastInsertId()
			val["account"] = blockData[i]["message_sender"]
			val["already_read"] = 0
			val["notification_time"] = time.Now().Unix()
			m = append(m, val)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotificationAccount).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		//Purchased(uint256 indexed saleId, uint256 indexed buyAmount)
		if blockData[i]["event_type"] == consts.EvPurchased {
			saleId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))
			buyAmount, _ := utils.Hex2BigInt(blockData[i]["topic2"].(string))
			payAmount, _ := utils.Hex2BigInt(blockData[i]["topic3"].(string))

			var swapData db.TbSwap
			sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Select()
			errTx = oo.SqlGet(sqlSel, &swapData)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			totalAmount, err := utils.StringAddBigInt(swapData.SoldAmount, buyAmount)
			if err != nil {
				oo.LogW("utils.StringAddBigInt err: %v", err)
				errTx = err
				return
			}

			var v = make(map[string]interface{})
			v["sold_amount"] = totalAmount
			if totalAmount == swapData.SaleAmount {
				v["status"] = consts.StatusEnded
			}
			sqlUp = oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Update(v)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var mTran = make([]map[string]interface{}, 0)
			var vTran = make(map[string]interface{})
			vTran["sale_id"] = saleId
			vTran["buyer"] = blockData[i]["message_sender"]
			vTran["buy_amount"] = buyAmount
			vTran["pay_amount"] = payAmount
			vTran["time"], _ = utils.Hex2Int64(blockData[i]["time_stamp"].(string))
			vTran["chain_id"] = chainId
			vTran["buy_token"] = swapData.SaleToken
			vTran["pay_token"] = swapData.ReceiveToken
			mTran = append(mTran, vTran)
			sqlIns := oo.NewSqler().Table(consts.TbNameSwapTransaction).Insert(mTran)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//for notification
			var t = time.Now().Unix()
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = swapData.SaleToken
			values["types"] = consts.TypesNamePublicSalePurchased
			values["activity_id"] = swapData.Id
			values["dao_logo"] = swapData.SaleTokenImg
			values["dao_name"] = blockData[i]["message_sender"]
			values["activity_name"] = swapData.Title
			values["start_time"] = t
			values["update_bool"] = 0
			notificationData = append(notificationData, values)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			var result sql.Result
			result, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var m = make([]map[string]interface{}, 0)
			var val1 = make(map[string]interface{})
			val1["notification_id"], _ = result.LastInsertId()
			val1["account"] = blockData[i]["message_sender"]
			val1["already_read"] = 0
			val1["notification_time"] = t
			m = append(m, val1)
			var val2 = make(map[string]interface{})
			val2["notification_id"], _ = result.LastInsertId()
			val2["account"] = swapData.Creator
			val2["already_read"] = 0
			val2["notification_time"] = t
			m = append(m, val2)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotificationAccount).InsertBatch(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}

		if blockData[i]["event_type"] == consts.EvCancelSale {
			saleId, _ := utils.Hex2Dec(blockData[i]["topic1"].(string))

			var v = make(map[string]interface{})
			v["status"] = consts.StatusCancel
			sqlUp = oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Update(v)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			//for notification
			var swapData db.TbSwap
			sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleId).Select()
			errTx = oo.SqlGet(sqlSel, &swapData)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
			var t = time.Now().Unix()
			var notificationData = make([]map[string]interface{}, 0)
			var values = make(map[string]interface{})
			values["chain_id"] = chainId
			values["dao_address"] = swapData.SaleToken
			values["types"] = consts.TypesNamePublicSaleCanceled
			values["activity_id"] = swapData.Id
			values["dao_logo"] = swapData.SaleTokenImg
			values["dao_name"] = ""
			values["activity_name"] = swapData.Title
			values["start_time"] = t
			values["update_bool"] = 0
			notificationData = append(notificationData, values)
			sqlIns := oo.NewSqler().Table(consts.TbNameNotification).Insert(notificationData)
			var result sql.Result
			result, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			var m = make([]map[string]interface{}, 0)
			var val = make(map[string]interface{})
			val["notification_id"], _ = result.LastInsertId()
			val["account"] = swapData.Creator
			val["already_read"] = 0
			val["notification_time"] = t
			m = append(m, val)
			sqlIns = oo.NewSqler().Table(consts.TbNameNotificationAccount).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}
		}
	}

}

func proposalInfo(daoAddress, proposalId, url string) (string, error) {
	// proposals(uint256 proposalId)
	const dataPrefix = "0x013cf08b"
	data := fmt.Sprintf("%s%s", dataPrefix, strings.TrimPrefix(proposalId, "0x"))
	res, err := utils.QueryMethodEthCall(daoAddress, data, url)
	if err != nil {
		oo.LogW("QueryMethodEthCall err: %v", err)
		return "", err
	}

	var outputParameters []string
	outputParameters = append(outputParameters, "bool")
	outputParameters = append(outputParameters, "address")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "uint256")
	outputParameters = append(outputParameters, "uint256")
	outputParameters = append(outputParameters, "uint256")
	outputParameters = append(outputParameters, "uint8")

	proposal, err := utils.Decode(outputParameters, strings.TrimPrefix(res.Result.(string), "0x"))
	if err != nil {
		oo.LogW("Decode failed. err: %v\n", err)
		return "", err
	}
	return proposal[2].(string), nil
}
