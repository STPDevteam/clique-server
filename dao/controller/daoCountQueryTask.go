package controller

import (
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"time"
)

func daoCountTask() {
	defer time.AfterFunc(time.Duration(60)*time.Minute, daoCountTask)

	var daoEntity []models.DaoModel
	sqlSel := oo.NewSqler().Table(consts.TbNameDao).Select()
	err := oo.SqlSelect(sqlSel, &daoEntity)
	if err != nil {
		oo.LogW("SQL err:%v", err)
		return
	}

	for index := range daoEntity {
		ls := daoEntity[index]

		var totalProposals int64
		sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("deprecated", 0).
			Where("chain_id", ls.ChainId).Where("dao_address", ls.DaoAddress).Count()
		err = oo.SqlGet(sqlSel, &totalProposals)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

		var members int64
		sqlMembers := oo.NewSqler().Table(consts.TbNameMember).
			Where("dao_address", ls.DaoAddress).
			Where("chain_id", ls.ChainId).
			Where("join_switch", 1).Count()
		err = oo.SqlGet(sqlMembers, &members)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

		var v = make(map[string]interface{})
		v["members"] = members
		v["total_proposals"] = totalProposals
		sqlUp := oo.NewSqler().Table(consts.TbNameDao).Where("id", ls.Id).Update(v)
		err = oo.SqlExec(sqlUp)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

	}
}
