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
	"time"
)

// @Summary team spaces list
// @Tags spaces
// @version 0.0.1
// @description team spaces list
// @Produce json
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Success 200 {object} models.ResTeamSpacesList
// @Router /stpdao/v2/spaces/list [get]
func TeamSpacesList(c *gin.Context) {
	var user *db.TbAccountModel
	login := c.GetBool(consts.KEY_LOGIN)
	if login {
		var ok bool
		user, ok = parseJWTCache(c)
		if !ok {
			return
		}
	}

	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")

	var wAccess [][]interface{}
	if !login || !IsSuperAdmin(int64(chainIdParam), daoAddressParam, user.Account) {
		wAccess = o.W("access", "public")
	}

	order := fmt.Sprintf("create_time ASC")
	page := ReqPagination{
		Offset: 0,
		Limit:  10,
	}
	list, total, err := PageTbTeamSpaces(order, page,
		o.W("chain_id", chainIdParam), o.W("dao_address", daoAddressParam), o.W("is_trash", 0), wAccess)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}

// @Summary create team spaces
// @Tags spaces
// @version 0.0.1
// @description create team spaces, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqCreateTeamSpaces true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/spaces/create [post]
func CreateTeamSpaces(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqCreateTeamSpaces
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.NewError(401, "You are not super admin."))
		return
	}
	if len(params.Title) > 20 {
		handleError(c, errs.NewError(400, "Title length is greater than 20."))
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["dao_address"] = params.DaoAddress
	v["creator"] = user.Account
	v["title"] = params.Title
	v["last_edit_time"] = time.Now().Unix()
	v["last_edit_by"] = user.Account
	v["access"] = params.Access
	m = append(m, v)
	err := o.Insert(consts.TbTeamSpaces, m)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary update team spaces
// @Tags spaces
// @version 0.0.1
// @description update team spaces, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqUpdateTeamSpaces true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/spaces/update [post]
func UpdateTeamSpaces(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqUpdateTeamSpaces
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.NewError(401, "You are not super admin."))
		return
	}
	if len(params.Title) > 20 {
		handleError(c, errs.NewError(400, "Title length is greater than 20."))
		return
	}

	var v = make(map[string]interface{})
	v["title"] = params.Title
	v["url"] = params.Url
	v["last_edit_time"] = time.Now().Unix()
	v["last_edit_by"] = user.Account
	v["access"] = params.Access
	err := o.Update(consts.TbTeamSpaces, v, o.W("id", params.TeamSpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary remove team spaces
// @Tags spaces
// @version 0.0.1
// @description remove team spaces, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqRemoveTeamSpaces true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/spaces/remove [post]
func TeamSpacesRemoveToTrash(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqRemoveTeamSpaces
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.NewError(401, "You are not super admin."))
		return
	}

	var v = make(map[string]interface{})
	v["is_trash"] = 1
	err := o.Update(consts.TbTeamSpaces, v, o.W("id", params.TeamSpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}

// @Summary delete completely team spaces
// @Tags spaces
// @version 0.0.1
// @description delete completely team spaces, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqDeleteTeamSpaces true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/spaces/delete [post]
func DeleteTeamSpaces(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqDeleteTeamSpaces
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.NewError(401, "You are not super admin."))
		return
	}

	err := o.Delete(consts.TbTeamSpaces, o.W("id", params.TeamSpacesId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonSuccess(c)
}
