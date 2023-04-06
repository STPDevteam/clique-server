package tasks

import (
	oo "github.com/Anna2024/liboo"
	"math"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"time"
)

var sleepTime int64

func UpdateSwapStatus() {
	var minTime int64
	defer func() {
		sleepTime = minTime - time.Now().Unix()
		if sleepTime > 60 {
			sleepTime = 60
		} else {
			sleepTime = 1
		}
		time.AfterFunc(time.Duration(sleepTime)*time.Second, UpdateSwapStatus)
	}()
	var min = consts.MaxValue

	var swapArr []db.TbSwap
	sqlSel := oo.NewSqler().Table(consts.TbNameSwap).OrWhere("status", consts.StatusSoon).OrWhere("status", consts.StatusNormal).Select()
	err := oo.SqlSelect(sqlSel, &swapArr)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		return
	}

	if len(swapArr) > 0 {
		for index := range swapArr {
			ls := swapArr[index]

			var status string
			if ls.StartTime <= time.Now().Unix() {
				status = consts.StatusNormal
			}
			if ls.EndTime <= time.Now().Unix() {
				status = consts.StatusEnded
			}
			if status == "" {
				continue
			}

			var v = make(map[string]interface{})
			v["status"] = status
			sqlUpd := oo.NewSqler().Table(consts.TbNameSwap).Where("id", ls.Id).Update(v)
			err := oo.SqlExec(sqlUpd)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				return
			}

			min = int(math.Min(float64(ls.StartTime), float64(min)))
		}
	}
	minTime = int64(min)
}
