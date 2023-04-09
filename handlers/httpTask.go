package handlers

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"stp_dao_v2/consts"
	"stp_dao_v2/db/o"
	"stp_dao_v2/errs"
	"stp_dao_v2/models"
	"strconv"
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

	//if !checkAdminForTaskCreate(params.Sign) {
	//	oo.LogD("SignData err not auth")
	//	handleError(c, errs.ErrUnAuthorized)
	//	return
	//}

	var weight float64
	sqlSel := oo.NewSqler().Table(consts.TbTask).Where("chain_id", params.Sign.ChainId).Where("dao_address", params.Sign.DaoAddress).
		Where("status", "A_notStarted").Max("weight")
	err := oo.SqlGet(sqlSel, &weight)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.Sign.ChainId
	v["dao_address"] = params.Sign.DaoAddress
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

// @Summary task list
// @Tags task
// @version 0.0.1
// @description task list
// @Produce json
// @Param offset query  int true "offset,page"
// @Param limit query  int true "limit,page"
// @Success 200 {object} models.
// @Router /stpdao/v2/task/list [get]
func TaskList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)

	order := fmt.Sprintf("weight ASC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbTask(consts.TbTask, order, page, o.W("is_trash", 0))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}
