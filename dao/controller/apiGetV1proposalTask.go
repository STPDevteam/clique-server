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

func (svc *Service) getV1Proposal() {
	defer time.AfterFunc(time.Duration(60)*time.Second, svc.getV1Proposal)

	var v1MaxId int
	sqlSel := oo.NewSqler().Table(consts.TbNameProposal).Max("id_v1")
	err := oo.SqlGet(sqlSel, &v1MaxId)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		return
	}
	if v1MaxId != 0 {
		v1MaxId = v1MaxId + 1
	}

	var v1Url = fmt.Sprintf(svc.appConfig.ApiV1ProposalUrl, v1MaxId)
	res, err := utils.GetV1ProposalHistory(v1Url)
	if err != nil {
		oo.LogW("GetV1LastBlockNumber failed error: %v", err)
		return
	}

	if res.Data != nil {
		for _, data := range res.Data {
			const dataPre = "0x3656de21"
			var proposalId = strings.TrimPrefix(data.Topic1, "0x")
			dataParam := fmt.Sprintf("%s%s", dataPre, proposalId)
			resContent, err := utils.QueryMethodEthCallByTag(data.Address, dataParam, svc.appConfig.ApiV1ProposalContentUrl, "latest")
			if err != nil {
				oo.LogW("QueryMethodEthCallByTag failed error: %v", err)
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

			var title string
			title, ok = decode[4].(string)
			if !ok {
				oo.LogW(".(string) failed.")
				return
			}
			var content = decode[4]
			startTime, _ := utils.Hex2Dec(data.Data[194:258])
			endTime, _ := utils.Hex2Dec(data.Data[258:322])
			proposalIdDec, _ := utils.Hex2Dec(data.Topic1)
			proposalCreator := utils.FixTo0x40String(data.Topic2)

			var entities []models.ProposalV1Model
			sqlSel = oo.NewSqler().Table(consts.TbNameProposalV1).Where("voting_v1", data.Address).Select()
			err = oo.SqlSelect(sqlSel, &entities)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}
			if entities == nil {
				continue
			}

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
			err = oo.SqlExec(sqlIns)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}
		}
	}

}
