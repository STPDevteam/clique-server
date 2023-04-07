package handlers

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"stp_dao_v2/consts"
	"stp_dao_v2/db/o"
	"stp_dao_v2/errs"
	"stp_dao_v2/models"
)

// @Summary create task
// @Tags task
// @version 0.0.1
// @description create task
// @Produce json
// @Param request body models.ReqCreateTask true "request"
// @Success 200 {object} models.
// @Router /stpdao/v2/task/create [post]
func CreateTask(c *gin.Context) {
	var params models.ReqCreateTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !checkAdminForTaskCreate(params) {
		oo.LogD("SignData err not auth")
		handleError(c, errs.ErrUnAuthorized)
		return
	}

	var weight float64
	sqlSel := oo.NewSqler().Table(consts.TbTask).Where("chain_id", params.ChainId).Where("dao_address", params.DaoAddress).
		Where("status", "A_notStarted").Max("weight")
	err := oo.SqlGet(sqlSel, &weight)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["dao_address"] = params.DaoAddress
	v["task_name"] = params.TaskName
	v["content"] = params.Content
	v["deadline"] = params.Deadline
	v["priority"] = params.Priority
	v["assign_account"] = params.AssignAccount
	v["proposal_id"] = params.ProposalId
	v["reward"] = params.Reward
	v["status"] = "A_notStarted"
	v["weight"] = weight + 10
	m = append(m, v)
	err = o.Insert(consts.TbTask, m)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}
