package handlers

import (
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/models"
)

func PageTbTask(table, order string, page ReqPagination, w ...[][]interface{}) (list []models.ResTaskList, total int64, err error) {
	var data []db.TbTask
	sqler := o.DBPre(table, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResTaskList, 0)
	for i := range data {
		ls := data[i]
		list = append(list, models.ResTaskList{
			ChainId:       ls.ChainId,
			DaoAddress:    ls.DaoAddress,
			TaskName:      ls.TaskName,
			Content:       ls.Content,
			Deadline:      ls.Deadline,
			Priority:      ls.Priority,
			AssignAccount: ls.AssignAccount,
			ProposalId:    ls.ProposalId,
			Reward:        ls.Reward,
			Status:        ls.Status,
			Weight:        ls.Weight,
		})
	}
	return list, total, nil
}
