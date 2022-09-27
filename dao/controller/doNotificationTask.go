package controller

import (
	"encoding/json"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"time"
)

func updateNotification() {
	defer time.AfterFunc(time.Duration(60)*time.Second, updateNotification)
	var nowTime = time.Now().Unix()

	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	var entities []models.NotificationModel
	sqlSel := oo.NewSqler().Table(consts.TbNameNotification).Where("update_bool", 1).Select()
	errTx = oo.SqlSelect(sqlSel, &entities)
	if errTx != nil {
		oo.LogW("query SQL notification failed. err:%v", errTx)
		return
	}

	if len(entities) != 0 {
		for index := range entities {
			var daoEntity []models.DaoModel
			sqlSel = oo.NewSqler().Table(consts.TbNameDao).
				Where("chain_id", entities[index].ChainId).
				Where("dao_address", entities[index].DaoAddress).Select()
			errTx = oo.SqlSelect(sqlSel, &daoEntity)
			if errTx != nil {
				oo.LogW("SQL err:%v", errTx)
				return
			}

			var sqlStr, activityName string
			if entities[index].Types == consts.TypesNameNewProposal {
				sqlStr = oo.NewSqler().Table(consts.TbNameProposal).
					Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].DaoAddress).
					Where("proposal_id", entities[index].ActivityId).Select("title")
			} else if entities[index].Types == consts.TypesNameAirdrop {
				sqlStr = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", entities[index].ActivityId).Select("title")
			}
			errTx = oo.SqlGet(sqlStr, &activityName)
			if errTx != nil {
				oo.LogW("SQL err:%v", errTx)
				return
			}

			var info = make(map[string]interface{})
			info["dao_logo"] = daoEntity[0].DaoLogo
			info["dao_name"] = daoEntity[0].DaoName
			info["activity_name"] = activityName
			info["update_bool"] = 0
			sqlUp := oo.NewSqler().Table(consts.TbNameNotification).
				Where("chain_id", entities[index].ChainId).
				Where("dao_address", entities[index].DaoAddress).Update(info)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			if entities[index].Types == consts.TypesNameNewProposal {
				var accountMember []string
				sqlSel = oo.NewSqler().Table(consts.TbNameMember).
					Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].DaoAddress).
					Where("join_switch", 1).Select("account")
				errTx = oo.SqlSelect(sqlSel, &accountMember)
				if errTx != nil {
					oo.LogW("SQL err:%v", errTx)
					return
				}
				var accountAdmin []string
				sqlSel = oo.NewSqler().Table(consts.TbNameAdmin).
					Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].DaoAddress).
					Where("account_level", consts.LevelAdmin).Select("account")
				errTx = oo.SqlSelect(sqlSel, &accountAdmin)
				if errTx != nil {
					oo.LogW("SQL err:%v", errTx)
					return
				}

				var accountArray = make([]string, 0)
				for _, member := range accountMember {
					accountArray = append(accountArray, member)
				}
				for _, admin := range accountAdmin {
					accountArray = append(accountArray, admin)
				}

				for _, account := range accountArray {
					var count int
					sqlSel = oo.NewSqler().Table(consts.TbNameNotificationAccount).
						Where("notification_id", entities[index].Id).
						Where("account", account).Count()
					errTx = oo.SqlGet(sqlSel, &count)
					if errTx != nil {
						oo.LogW("SQL err: %v\n", errTx)
						return
					}

					if count == 0 {
						var m = make([]map[string]interface{}, 0)
						var v = make(map[string]interface{})
						v["notification_id"] = entities[index].Id
						v["account"] = account
						v["already_read"] = 0
						v["notification_time"] = nowTime
						m = append(m, v)
						sqlIns := oo.NewSqler().Table(consts.TbNameNotificationAccount).Insert(m)
						_, errTx = oo.SqlxTxExec(tx, sqlIns)
						if errTx != nil {
							oo.LogW("SQL err: %v\n", errTx)
							return
						}
					}
				}

			}

			if entities[index].Types == consts.TypesNameAirdrop {
				var addressEntity []models.AirdropAddressModel
				sqlSel = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", entities[index].ActivityId).Select()
				errTx = oo.SqlSelect(sqlSel, &addressEntity)
				if errTx != nil {
					oo.LogW("SQL err:%v", errTx)
					return
				}

				var data models.AirdropAddressArray
				errTx = json.Unmarshal([]byte(addressEntity[0].Content), &data)
				if errTx != nil {
					oo.LogW("Json Unmarshal err:%v", errTx)
					return
				}

				if len(data.Address) != 0 {
					var m = make([]map[string]interface{}, 0)
					for _, account := range data.Address {
						var v = make(map[string]interface{})
						v["notification_id"] = entities[index].Id
						v["account"] = account
						v["already_read"] = 0
						v["notification_time"] = nowTime
						m = append(m, v)
					}
					sqlIns := oo.NewSqler().Table(consts.TbNameNotificationAccount).InsertBatch(m)
					_, errTx = oo.SqlxTxExec(tx, sqlIns)
					if errTx != nil {
						oo.LogW("SQL err: %v", errTx)
						return
					}
				}
			}

			if entities[index].Types == consts.TypesNamePublicSale {

			}
		}
	}
}
