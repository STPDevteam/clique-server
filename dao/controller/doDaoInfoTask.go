package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"math"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
)

func (svc *Service) updateDaoInfoTask() {
	defer time.AfterFunc(time.Duration(60)*time.Second, svc.updateDaoInfoTask)
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
						res, errQ := utils.QueryMethodEthCall(entities[index].DaoAddress, data, url)
						if errQ != nil {
							oo.LogW("QueryDaoInfo failed. chainId:%d. err: %v\n", chainId, errQ)
							return
						}
						val, ok := res.Result.(string)
						if !ok {
							oo.LogW("QueryDaoInfo failed. chainId:%d. err: %v\n", chainId, errQ)
							return
						}

						if val != "" {
							var outputParameters []string
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")
							outputParameters = append(outputParameters, "string")

							daoInfo, errDe := utils.Decode(outputParameters, strings.TrimPrefix(val, "0x"))
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

	category := daoInfo[2]

	var v = make(map[string]interface{})
	v["dao_logo"] = daoInfo[7]
	v["dao_name"] = daoInfo[0]
	v["handle"] = daoInfo[1]
	v["description"] = daoInfo[3]
	v["twitter"] = daoInfo[4].(string)[:int(math.Min(float64(len(daoInfo[4].(string))), 256))]
	v["github"] = daoInfo[5].(string)[:int(math.Min(float64(len(daoInfo[5].(string))), 256))]
	v["discord"] = daoInfo[6].(string)[:int(math.Min(float64(len(daoInfo[6].(string))), 256))]
	v["website"] = daoInfo[8].(string)[:int(math.Min(float64(len(daoInfo[8].(string))), 256))]
	v["update_bool"] = 0
	sqlIns := oo.NewSqler().Table(consts.TbNameDao).Where("dao_address", daoAddress).Where("chain_id", chainId).Update(v)
	_, errTx = oo.SqlxTxExec(tx, sqlIns)
	if errTx != nil {
		oo.LogW("SQL failed. err: %v\n", errTx)
		return
	}

	var daoId int
	sqlSelDId := oo.NewSqler().Table(consts.TbNameDao).Where("dao_address", daoAddress).Where("chain_id", chainId).Select("id")
	errTx = oo.SqlGet(sqlSelDId, &daoId)
	if errTx != nil {
		oo.LogW("SQL err: %v\n", errTx)
		return
	}

	/* event setting: maybe delete the category first */
	sqlDel := oo.NewSqler().Table(consts.TbNameDaoCategory).Where("dao_id", daoId).Delete()
	_, errTx = oo.SqlxTxExec(tx, sqlDel)
	if errTx != nil {
		oo.LogW("SQL err: %v\n", errTx)
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
		errTx = oo.SqlExec(sqlUP)
		if errTx != nil {
			oo.LogW("SQL err: %v\n", errTx)
			return
		}

		var categoryId int
		sqlSelCId := oo.NewSqler().Table(consts.TbNameCategory).Where("category_name", categoryName).Select("id")
		errTx = oo.SqlGet(sqlSelCId, &categoryId)
		if errTx != nil {
			oo.LogW("SQL err: %v\n", errTx)
			return
		}

		sqlInsCategory := fmt.Sprintf(`INSERT INTO %s (dao_id,category_id) VALUES (%d,%d)`,
			consts.TbNameDaoCategory,
			daoId,
			categoryId,
		)
		_, errTx = oo.SqlxTxExec(tx, sqlInsCategory)
		if errTx != nil {
			oo.LogW("SQL failed. err: %v\n", errTx)
			return
		}
	}
}
