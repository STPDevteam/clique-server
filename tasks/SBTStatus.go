package tasks

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"time"
)

func UpdateSBTStatus() {
	defer time.AfterFunc(time.Duration(60)*time.Second, UpdateSBTStatus)

	sbtArr, err := db.SelectTbSBT(o.W(fmt.Sprintf(`status IN ('%s','%s')`, consts.StatusSoon, consts.StatusActive)))
	if err != nil {
		oo.LogW("SQL err: %v", err)
		return
	}

	for index := range sbtArr {
		ls := sbtArr[index]

		var status string
		if ls.StartTime <= time.Now().Unix() {
			status = consts.StatusActive
		}
		if ls.EndTime <= time.Now().Unix() {
			status = consts.StatusEnded
		}
		if status == "" {
			continue
		}

		var v = make(map[string]interface{})
		v["status"] = status
		err = o.Update(consts.TbSBT, v, o.W("id", ls.Id))
		if err != nil {
			oo.LogW("SQL err: %v", err)
			return
		}

	}

}
