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
// @description create task, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqCreateTask true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/create [post]
func CreateTask(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqCreateTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	spacesData, err := db.GetTbTeamSpaces(o.W("id", params.SpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	_, ok = IsAboveAdmin(spacesData.ChainId, spacesData.DaoAddress, user.Account)
	if !ok {
		handleError(c, errs.ErrUnAuthorized)
		return
	}

	var weight float64
	sqlSel := o.Sqler(consts.TbTask, o.W("spaces_id", params.SpacesId),
		o.W("status", consts.Task_status_A_notStarted)).Max("weight")
	err = oo.SqlGet(sqlSel, &weight)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["spaces_id"] = params.SpacesId
	v["task_name"] = params.TaskName
	v["content"] = params.Content
	v["deadline"] = params.Deadline
	v["priority"] = params.Priority
	v["assign_account"] = params.AssignAccount
	v["proposal_id"] = params.ProposalId
	v["reward"] = params.Reward
	v["status"] = consts.Task_status_A_notStarted
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
// @description update task, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqUpdateTask true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/update [post]
func UpdateTask(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqUpdateTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	spacesData, err := db.GetTbTeamSpaces(o.W("id", params.SpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	_, ok1 := IsAboveAdmin(spacesData.ChainId, spacesData.DaoAddress, user.Account)
	ok2 := IsTaskAssign(params.TaskId, user.Account)
	if !ok1 && !ok2 {
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
	err = o.Update(consts.TbTask, v, o.W("id", params.TaskId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary remove task to trash
// @Tags task
// @version 0.0.1
// @description remove task to trash, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqRemoveTask true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/task/remove [post]
func TaskRemoveToTrash(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqRemoveTask
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	spacesData, err := db.GetTbTeamSpaces(o.W("id", params.SpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	_, ok = IsAboveAdmin(spacesData.ChainId, spacesData.DaoAddress, user.Account)
	if !ok {
		handleError(c, errs.ErrUnAuthorized)
		return
	}

	var v = make(map[string]interface{})
	v["is_trash"] = 1
	for _, val := range params.TaskId {
		err = o.Update(consts.TbTask, v, o.W("id", val))
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
// @Param spacesId query int true "spacesId"
// @Param status query string false "status:A_notStarted;B_inProgress;C_done;D_notStatus"
// @Param priority query string false "priority:A_low;B_medium;C_high"
// @Success 200 {object} models.ResTaskList
// @Router /stpdao/v2/task/list [get]
func TaskList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	spacesId := c.Query("spacesId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	spacesIdParam, _ := strconv.Atoi(spacesId)
	statusParam := c.Query("status")
	priorityParam := c.Query("priority")

	var wStatus, wPriority [][]interface{}
	if statusParam != "" {
		wStatus = o.W("status", statusParam)
	}
	if priorityParam != "" {
		wPriority = o.W("priority", priorityParam)
	}

	order := fmt.Sprintf("weight ASC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbTask(order, page, o.W("spaces_id", spacesIdParam), o.W("is_trash", 0), wStatus, wPriority)
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
		SpacesId:       task.SpacesId,
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
