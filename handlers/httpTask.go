package handlers

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
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
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/create [post]
func CreateTask(c *gin.Context) {
	var params models.ReqCreateTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !checkAdminForTaskCreate(params.Sign) {
		oo.LogD("SignData err not auth")
		handleError(c, errs.ErrUnAuthorized)
		return
	}

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

// @Summary update task
// @Tags task
// @version 0.0.1
// @description update task
// @Produce json
// @Param request body models.ReqUpdateTask true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/update [post]
func UpdateTask(c *gin.Context) {
	var params models.ReqUpdateTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !checkAdminForTaskCreate(params.Sign) {
		oo.LogD("SignData err not auth")
		handleError(c, errs.ErrUnAuthorized)
		return
	}

	var v = make(map[string]interface{})
	v["task_name"] = params.TaskName
	v["content"] = params.Content
	v["deadline"] = params.Deadline
	v["priority"] = params.Priority
	v["assign_account"] = params.AssignAccount
	v["proposal_id"] = params.ProposalId
	v["reward"] = params.Reward
	v["status"] = params.Status
	v["weight"] = params.Weight
	err := o.Update(consts.TbTask, v, o.W("id", params.TaskId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary remove task to trash
// @Tags task
// @version 0.0.1
// @description remove task to trash
// @Produce json
// @Param request body models.ReqRemoveTask true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/remove [post]
func TaskRemoveToTrash(c *gin.Context) {
	var params models.ReqRemoveTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !checkAdminForTaskCreate(params.Sign) {
		oo.LogD("SignData err not auth")
		handleError(c, errs.ErrUnAuthorized)
		return
	}

	var v = make(map[string]interface{})
	v["is_trash"] = 1
	for _, val := range params.TaskId {
		err := o.Update(consts.TbTask, v, o.W("id", val))
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}
	}

	jsonSuccess(c)
}

// @Summary task list
// @Tags task
// @version 0.0.1
// @description task list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Param status query string false "status:A_notStarted;B_inProgress;C_done;D_notStatus"
// @Success 200 {object} models.ResTaskList
// @Router /stpdao/v2/task/list [get]
func TaskList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	chainId := c.Query("chainId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")
	statusParam := c.Query("status")

	var wStatus [][]interface{}
	if statusParam != "" {
		wStatus = o.W("status", statusParam)
	}

	order := fmt.Sprintf("weight ASC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbTask(order, page,
		o.W("chain_id", chainIdParam), o.W("dao_address", daoAddressParam), o.W("is_trash", 0), wStatus)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}

// @Summary task detail
// @Tags task
// @version 0.0.1
// @description task detail
// @Produce json
// @Success 200 {object} models.ResTaskDetail
// @Router /stpdao/v2/task/detail/:taskId [get]
func TaskDetail(c *gin.Context) {
	taskId := c.Param("taskId")
	taskIdParam, _ := strconv.Atoi(taskId)

	task, err := db.GetTbTask(o.W("id", taskIdParam))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	var avatar, nickname string
	if task.AssignAccount != "" {
		account, err := db.GetTbAccountModel(o.W("account", task.AssignAccount))
		if handleErrorIfExistsExceptNoRows(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}
		avatar = account.AccountLogo.String
		nickname = account.Nickname.String
	}

	data := models.ResTaskDetail{
		TaskId:         task.Id,
		ChainId:        task.ChainId,
		DaoAddress:     task.DaoAddress,
		TaskName:       task.TaskName,
		Content:        task.Content,
		Deadline:       task.Deadline,
		Priority:       task.Priority,
		AssignAccount:  task.AssignAccount,
		AssignAvatar:   avatar,
		AssignNickname: nickname,
		ProposalId:     task.ProposalId,
		Reward:         task.Reward,
		Status:         task.Status,
		Weight:         task.Weight,
	}

	jsonData(c, data)
}
