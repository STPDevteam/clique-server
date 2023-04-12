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

// @Summary jobs apply
// @Tags jobs
// @version 0.0.1
// @description jobs apply
// @Produce json
// @Param request body models.ReqJobsApply true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/jobs/apply [post]
func JobsApply(c *gin.Context) {
	var params models.ReqJobsApply
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	//if !checkLogin(&params.Sign) {
	//	oo.LogD("SignData err not auth")
	//	handleError(c, errs.ErrUnAuthorized)
	//	return
	//}

	countJobs, err := o.Count(consts.TbJobs, o.W("chain_id", params.ChainId), o.W("dao_address", params.DaoAddress),
		o.W("account", params.Sign.Account), o.W("job", params.ApplyRole))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}
	if countJobs > 0 {
		handleError(c, errs.NewError(400, "You have successfully applied."))
		return
	}

	countJobsApply, err := o.Count(consts.TbJobsApply, o.W("chain_id", params.ChainId), o.W("dao_address", params.DaoAddress),
		o.W("account", params.Sign.Account), o.W("status", consts.Jobs_Status_InApplication))
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
	v["account"] = params.Sign.Account
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
// @Router /stpdao/v2/jobs/list [get]
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

// @Summary jobs list
// @Tags jobs
// @version 0.0.1
// @description jobs list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Success 200 {object} models.
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
