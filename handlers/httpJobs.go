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

// @Summary jobs apply
// @Tags jobs
// @version 0.0.1
// @description jobs apply, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqJobsApply true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/jobs/apply [post]
func JobsApply(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqJobsApply
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}
	if params.ApplyRole != consts.Jobs_B_admin && params.ApplyRole != consts.Jobs_C_member {
		handleError(c, errs.ErrParam)
		return
	}

	countJobs, err := o.Count(consts.TbJobs, o.W("chain_id", params.ChainId), o.W("dao_address", params.DaoAddress),
		o.W("account", user.Account), o.W("job", params.ApplyRole))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}
	if countJobs > 0 {
		handleError(c, errs.NewError(400, "You have successfully applied."))
		return
	}

	if params.ApplyRole == consts.Jobs_C_member {
		var m = make([]map[string]interface{}, 0)
		var v = make(map[string]interface{})
		v["chain_id"] = params.ChainId
		v["dao_address"] = params.DaoAddress
		v["account"] = user.Account
		v["job"] = consts.Jobs_C_member
		m = append(m, v)
		err = o.Insert(consts.TbJobs, m)
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}

		jsonSuccess(c)
		return
	}

	countJobsApply, err := o.Count(consts.TbJobsApply, o.W("chain_id", params.ChainId), o.W("dao_address", params.DaoAddress),
		o.W("account", user.Account), o.W("status", consts.Jobs_Status_InApplication))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}
	if countJobsApply > 0 {
		handleError(c, errs.NewError(400, "Application submitted."))
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["dao_address"] = params.DaoAddress
	v["account"] = user.Account
	v["apply_role"] = params.ApplyRole
	v["message"] = params.Message
	m = append(m, v)
	err = o.Insert(consts.TbJobsApply, m)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary jobs apply list
// @Tags jobs
// @version 0.0.1
// @description jobs apply list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Success 200 {object} models.ResJobsApplyList
// @Router /stpdao/v2/jobs/apply/list [get]
func JobsApplyList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	chainId := c.Query("chainId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")

	order := fmt.Sprintf("create_time DESC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbJobsApply(order, page, o.W("chain_id", chainIdParam), o.W("dao_address", daoAddressParam),
		o.W("status", consts.Jobs_Status_InApplication))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}

// @Summary jobs apply review
// @Tags jobs
// @version 0.0.1
// @description jobs apply review, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqJobsApplyReview true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/jobs/apply/review [post]
func JobsApplyReview(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqJobsApplyReview
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.NewError(401, "You are not super admin."))
		return
	}

	jobsApplyData, err := db.GetTbJobsApply(o.W("id", params.JobsApplyId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	var status string
	if params.IsPass {
		status = consts.Jobs_Status_Agree
	} else {
		status = consts.Jobs_Status_Reject
	}
	var val = make(map[string]interface{})
	val["status"] = status
	err = o.Update(consts.TbJobsApply, val, o.W("id", params.JobsApplyId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	if params.IsPass {
		jobData, err := db.GetTbJobs(
			o.W("chain_id", jobsApplyData.ChainId),
			o.W("dao_address", jobsApplyData.DaoAddress),
			o.W("account", jobsApplyData.Account))
		if handleErrorIfExistsExceptNoRows(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}

		var v = make(map[string]interface{})
		v["job"] = jobsApplyData.ApplyRole
		var e error
		if jobData.Job != "" {
			e = o.Update(consts.TbJobs, v, o.W("chain_id", jobsApplyData.ChainId),
				o.W("dao_address", jobsApplyData.DaoAddress),
				o.W("account", jobsApplyData.Account))
		} else {
			var m = make([]map[string]interface{}, 0)
			v["chain_id"] = jobsApplyData.ChainId
			v["dao_address"] = jobsApplyData.DaoAddress
			v["account"] = jobsApplyData.Account
			m = append(m, v)
			e = o.Insert(consts.TbJobs, m)
		}
		if handleErrorIfExists(c, e, errs.ErrServer) {
			oo.LogW("SQL err:%v", e)
			return
		}
	}

	jsonSuccess(c)
}

// @Summary jobs list
// @Tags jobs
// @version 0.0.1
// @description jobs list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Success 200 {object} models.ResJobsList
// @Router /stpdao/v2/jobs/list [get]
func JobsList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	chainId := c.Query("chainId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")

	order := fmt.Sprintf("job ASC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbJobs(order, page, o.W("chain_id", chainIdParam), o.W("dao_address", daoAddressParam))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}

// @Summary jobs alter
// @Tags jobs
// @version 0.0.1
// @description jobs alter, only superAdmin or admin, change admin/member to admin/member/noRole, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqJobsAlter true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/jobs/alter [post]
func JobsAlter(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqJobsAlter
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	role, ok := IsAboveAdmin(params.ChainId, params.DaoAddress, user.Account)
	if !ok {
		handleError(c, errs.NewError(401, "You are not admin."))
		return
	}

	jobData, err := db.GetTbJobs(o.W("id", params.JobId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	if jobData.Job == consts.Jobs_B_admin {
		if role != consts.Jobs_A_superAdmin {
			handleError(c, errs.ErrUnAuthorized)
			return
		}

		if params.ChangeTo == consts.Jobs_B_admin {
			var val = make(map[string]interface{})
			val["job"] = consts.Jobs_B_admin
			err = o.Update(consts.TbJobs, val, o.W("id", params.JobId))
			if handleErrorIfExists(c, err, errs.ErrServer) {
				oo.LogW("SQL err:%v", err)
				return
			}
		}
	} else if jobData.Job == consts.Jobs_C_member {
		if role != consts.Jobs_A_superAdmin && role != consts.Jobs_B_admin {
			handleError(c, errs.ErrUnAuthorized)
			return
		}
	}

	if params.ChangeTo == consts.Jobs_C_member {
		var val = make(map[string]interface{})
		val["job"] = consts.Jobs_C_member
		err = o.Update(consts.TbJobs, val, o.W("id", params.JobId))
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}
	} else if params.ChangeTo == consts.Jobs_noRole {
		err = o.Delete(consts.TbJobs, o.W("id", params.JobId))
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SQL err:%v", err)
			return
		}
	}

	jsonSuccess(c)
}

// @Summary jobs identity
// @Tags jobs
// @version 0.0.1
// @description jobs identity, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/jobs/identity [get]
func JobsIdentity(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")

	jobs, err := db.GetTbJobs(o.W("chain_id", chainIdParam),
		o.W("dao_address", daoAddressParam), o.W("account", user.Account))
	if handleErrorIfExistsExceptNoRows(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonData(c, jobs.Job)
}
