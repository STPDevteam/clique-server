package controller

import (
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
)

func updateOnlineData() {

	var count int
	sqlStr := oo.NewSqler().Table(consts.TbNameVote).Count()
	err := oo.SqlGet(sqlStr, &count)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		return
	}

	if count == 0 {
		var entities []models.EventHistoricalModel
		sqler := oo.NewSqler().Table(consts.TbNameEventHistorical).Where("event_type", consts.EvVote).Select()
		err = oo.SqlSelect(sqler, &entities)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

		var m = make([]map[string]interface{}, 0)
		for index := range entities {
			amount, _ := utils.Hex2BigInt(entities[index].Data[:66])

			var v = make(map[string]interface{})
			v["chain_id"] = entities[index].ChainId
			v["dao_address"] = entities[index].Address
			v["proposal_id"] = utils.Hex2Dec(entities[index].Topic1)
			v["voter"] = utils.FixTo0x40String(entities[index].Topic2)
			v["option_index"] = utils.Hex2Dec(entities[index].Topic3)
			v["amount"] = amount.String()
			v["nonce"] = utils.Hex2Dec(entities[index].Data[66:130])
			m = append(m, v)
		}
		sqlIns := oo.NewSqler().Table(consts.TbNameVote).InsertBatch(m)
		err = oo.SqlExec(sqlIns)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}
	}

}
