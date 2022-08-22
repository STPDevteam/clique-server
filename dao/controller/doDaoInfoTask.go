package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
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
							saveDaoInfoAndCategory(daoInfo, entities[index].DaoAddress, entities[index].TokenAddress, chainId)
						}
					}
				}
			}
		}
	}
}

func saveDaoInfoAndCategory(daoInfo []interface{}, daoAddress, tokenAddress string, chainId int) {
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	//id := 1
	//t := "0xdf5e0e81dff6faf3a7e52ba697820c5e32d806a8"
	tokensImgTask(chainId, tokenAddress)

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
