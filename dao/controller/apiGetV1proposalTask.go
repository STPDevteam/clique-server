package controller

import (
	"errors"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"math"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
)

func (svc *Service) getV1Proposal() {
	defer time.AfterFunc(time.Duration(60)*time.Second, svc.getV1Proposal)

	var v1StartId int
	sqlSel := oo.NewSqler().Table(consts.TbNameProposalV1).Min("start_id_v1")
	err := oo.SqlGet(sqlSel, &v1StartId)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		return
	}
	if v1StartId == -1 {
		var updateStartId = make(map[string]interface{})
		updateStartId["start_id_v1"] = 0
		sqlUp := oo.NewSqler().Table(consts.TbNameProposalV1).Where("start_id_v1 = -1").Update(updateStartId)
		err = oo.SqlExec(sqlUp)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

		v1StartId = 0
	}

	var v1Url = fmt.Sprintf(svc.appConfig.ApiV1ProposalUrl, v1StartId)
	res, err := utils.GetV1ProposalHistory(v1Url)
	if err != nil {
		oo.LogW("GetV1LastBlockNumber failed err: %v", err)
		return
	}

	var updateId int
	if res.Data != nil {
		for _, data := range res.Data {
			var countExist int
			sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("id_v1", data.Id).Count()
			err = oo.SqlGet(sqlSel, &countExist)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}

			if countExist == 0 {

				const dataPre = "0x3656de21"
				var proposalId = strings.TrimPrefix(data.Topic1, "0x")
				dataParam := fmt.Sprintf("%s%s", dataPre, proposalId)
				resContent, err := utils.QueryMethodEthCallByTag(data.Address, dataParam, svc.appConfig.ApiV1ProposalContentUrl, "latest")
				if err != nil {
					oo.LogW("QueryMethodEthCallByTag failed err: %v", err)
					return
				}
				val, ok := resContent.Result.(string)
				if !ok {
					oo.LogW(".(string) failed.")
					return
				}
				var outputParameters []string
				outputParameters = append(outputParameters, "uint256")
				outputParameters = append(outputParameters, "uint256")
				outputParameters = append(outputParameters, "address")
				outputParameters = append(outputParameters, "string")
				outputParameters = append(outputParameters, "string")
				decode, err := utils.Decode(outputParameters, val[66:])
				if err != nil {
					oo.LogW("Decode failed. err: %v", err)
					return
				}

				err = saveV1Proposal(decode, data)
				if err != nil {
					oo.LogW("saveV1Proposal failed err: %v", err)
					return
				}
			}
			updateId = data.Id
		}

		if updateId != 0 {
			var updateStartId = make(map[string]interface{})
			updateStartId["start_id_v1"] = updateId
			sqlUp := oo.NewSqler().Table(consts.TbNameProposalV1).Where("start_id_v1 >= 0").Update(updateStartId)
			err = oo.SqlExec(sqlUp)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}
		}

	}

}

func saveV1Proposal(decode []interface{}, data models.V1ProposalData) error {
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	var title string
	title, ok := decode[3].(string)
	if !ok {
		oo.LogW(".(string) failed.")
		errTx = errors.New(".(string) failed")
		return errTx
	}
	var content = decode[4]
	startTime, _ := utils.Hex2Dec(data.Data[194:258])
	endTime, _ := utils.Hex2Dec(data.Data[258:322])
	proposalIdDec, _ := utils.Hex2Dec(data.Topic1)
	proposalCreator := utils.FixTo0x40String(data.Topic2)

	var entities []models.ProposalV1Model
	sqlSel := oo.NewSqler().Table(consts.TbNameProposalV1).Where("voting_v1", data.Address).Select()
	errTx = oo.SqlSelect(sqlSel, &entities)
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
		return errTx
	}

	if len(entities) != 0 {

		var count int
		sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", entities[0].ChainId).
			Where("dao_address", entities[0].DaoAddress).Count()
		errTx = oo.SqlGet(sqlSel, &count)
		if errTx != nil {
			oo.LogW("SQL err: %v", errTx)
			return errTx
		}

		if count == 1 {
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = entities[0].ChainId
			v["dao_address"] = entities[0].DaoAddress
			v["proposal_id"] = proposalIdDec
			v["title"] = title[:int(math.Min(float64(len(title)), 500))]
			v["id_v1"] = data.Id
			v["content_v1"] = content
			v["proposer"] = proposalCreator
			v["start_time"] = startTime
			v["end_time"] = endTime
			v["version"] = "v1"
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameProposal).Insert(m)
			_, errTx = oo.SqlxTxExec(tx, sqlIns)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return errTx
			}

			// for dao order with proposal total
			var totalProposal int
			sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("chain_id", entities[0].ChainId).Where("dao_address", entities[0].DaoAddress).Count()
			errTx = oo.SqlGet(sqlSel, &totalProposal)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return errTx
			}
			var weight = make(map[string]interface{})
			weight["weight"] = totalProposal + 1
			sqlUp := oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", entities[0].ChainId).Where("dao_address", entities[0].DaoAddress).Update(weight)
			_, errTx = oo.SqlxTxExec(tx, sqlUp)
			if errTx != nil {
				oo.LogW("SQL err: %v", errTx)
				return errTx
			}
		}

	}

	return nil
}
