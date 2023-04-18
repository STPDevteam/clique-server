package tasks

import (
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"time"
)

func OnceTaskForTeamSpaces() {
	daoArr, err := db.SelectTbDaoModel()
	if err != nil {
		oo.LogW("SQL err:%v", err)
		return
	}

	oo.LogD("start init team spaces...")

	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL errTx: %v", errTx)
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	for i := range daoArr {
		ls := daoArr[i]

		countSpace, err := o.Count(consts.TbTeamSpaces, o.W("chain_id", ls.ChainId), o.W("dao_address", ls.DaoAddress))
		if err != nil {
			errTx = err
			oo.LogW("SQL errTx:%v", errTx)
			return
		}
		if countSpace == 0 {
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = ls.ChainId
			v["dao_address"] = ls.DaoAddress
			v["creator"] = ls.Creator
			v["title"] = "General"
			v["last_edit_time"] = time.Now().Unix()
			v["last_edit_by"] = ls.Creator
			v["access"] = "public"
			m = append(m, v)
			_, errTx = o.InsertTx(tx, consts.TbTeamSpaces, m)
			if errTx != nil {
				oo.LogW("SQL errTx:%v", errTx)
				return
			}
		}

		countJob, err := o.Count(consts.TbJobs, o.W("chain_id", ls.ChainId), o.W("dao_address", ls.DaoAddress))
		if err != nil {
			errTx = err
			oo.LogW("OnceTaskForTeamSpaces errTx:%v", errTx)
			return
		}
		if countJob == 0 {
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["chain_id"] = ls.ChainId
			v["dao_address"] = ls.DaoAddress
			v["account"] = ls.Creator
			v["job"] = consts.Jobs_A_superAdmin
			m = append(m, v)
			_, errTx = o.InsertTx(tx, consts.TbJobs, m)
			if errTx != nil {
				oo.LogW("SQL errTx:%v", errTx)
				return
			}
		}
	}

	oo.LogD("team spaces ended...")
}
