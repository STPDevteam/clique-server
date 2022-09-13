package controller

import (
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"time"
)

func updateNotification() {
	defer time.AfterFunc(time.Duration(60)*time.Second, updateNotification)

	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	var entities []models.NotificationModel
	sqlSel := oo.NewSqler().Table(consts.TbNameNotification).Where("update_bool", 1).Select()
	err := oo.SqlSelect(sqlSel, &entities)
	if err != nil {
		oo.LogW("query SQL notification failed. err:%v", err)
		return
	}

	if len(entities) != 0 {
		for index := range entities {
			var daoLogo string
			sqlSel = oo.NewSqler().Table(consts.TbNameDao).
				Where("chain_id", entities[index].ChainId).
				Where("dao_address", entities[index].DaoAddress).Select("dao_logo")
			err = oo.SqlGet(sqlSel, &daoLogo)
			if err != nil {
				oo.LogW("query SQL dao_logo failed. err:%v", err)
				return
			}

			sqlUp := fmt.Sprintf(`UPDATE %s SET dao_logo='%s',update_bool=%t WHERE chain_id=%d AND dao_address='%s' AND update_bool=%t`,
				consts.TbNameNotification,
				daoLogo,
				false,
				entities[index].ChainId,
				entities[index].DaoAddress,
				true,
			)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return
			}

			if entities[index].Types == consts.TypesNameProposal {
				var accountMember []string
				sqlSel = oo.NewSqler().Table(consts.TbNameMember).
					Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].DaoAddress).
					Where("join_switch", 1).Select("account")
				err = oo.SqlSelect(sqlSel, &accountMember)
				if err != nil {
					oo.LogW("query SQL account failed. err:%v", err)
					return
				}
				var accountAdmin []string
				sqlSel = oo.NewSqler().Table(consts.TbNameAdmin).
					Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].DaoAddress).
					Where("account_level", consts.LevelAdmin).Select("account")
				err = oo.SqlSelect(sqlSel, &accountAdmin)
				if err != nil {
					oo.LogW("query SQL account failed. err:%v", err)
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
					err = oo.SqlGet(sqlSel, &count)
					if err != nil {
						oo.LogW("SQL err: %v\n", err)
						return
					}

					if count == 0 {
						sqlIns := fmt.Sprintf(`INSERT INTO %s (notification_id,account,already_read) VALUES (%d,'%s',%t)`,
							consts.TbNameNotificationAccount,
							entities[index].Id,
							account,
							false,
						)
						_, errTx = oo.SqlxTxExec(tx, sqlIns)
						if errTx != nil {
							oo.LogW("SQL err: %v\n", errTx)
							return
						}
					}
				}

			}

			if entities[index].Types == consts.TypesNameAirdrop {
				var entity []models.AirdropAddressModel
				sqlSel = oo.NewSqler().Table(consts.TbNameAirdropAddress).Where("id", entities[index].ActivityId).Select()
				err = oo.SqlSelect(sqlSel, &entity)
				if err != nil {
					oo.LogW("query SQL account failed. err:%v", err)
					return
				}

				var data models.AirdropAddressArray
				err = json.Unmarshal([]byte(entity[0].Content), &data)
				if err != nil {
					oo.LogW("Json Unmarshal err:%v", err)
					return
				}

				if len(data.Address) != 0 {
					var m = make([]map[string]interface{}, 0)
					for _, account := range data.Address {
						var v = make(map[string]interface{})
						v["notification_id"] = entities[index].Id
						v["account"] = account
						v["already_read"] = false
						m = append(m, v)
					}
					sqlIns := oo.NewSqler().Table(consts.TbNameNotificationAccount).InsertBatch(m)
					err = oo.SqlExec(sqlIns)
					if err != nil {
						oo.LogW("SQL err: %v", err)
						return
					}
				}
			}

			if entities[index].Types == consts.TypesNamePublicSale {

			}
		}
	}
}
