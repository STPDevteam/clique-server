package handlers

import (
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/errs"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"strings"
)

// @Summary create sbt
// @Tags sbt
// @version 0.0.1
// @description create task, need superAdmin, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqSBTCreate true "request"
// @Success 200 {object} models.ResSBTCreate
// @Router /stpdao/v2/sbt/create [post]
func CreateSBT(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	var params models.ReqSBTCreate
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	if !IsSuperAdmin(params.ChainId, params.DaoAddress, user.Account) {
		handleError(c, errs.ErrNoPermission)
		return
	}

	daoData, err := db.GetTbDao(o.W("chain_id", params.ChainId), o.W("dao_address", params.DaoAddress))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		return
	}

	tx, errTx := oo.NewSqlxTx()
	if handleErrorIfExists(c, errTx, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		return
	}
	defer oo.CloseSqlxTx(tx, &errTx)

	mSBT := []map[string]any{{
		"chain_id":       params.ChainId,
		"dao_address":    params.DaoAddress,
		"token_chain_id": params.TokenChainId,
		"file_url":       params.FileUrl,
		"item_name":      params.ItemName,
		"introduction":   params.Introduction,
		"total_supply":   params.TotalSupply,
		"start_time":     params.StartTime,
		"end_time":       params.EndTime,
		"way":            params.Way,
		"whitelist":      params.Whitelist,
	}}
	res, err := o.InsertTx(tx, consts.TbSBT, mSBT)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		errTx = err
		return
	}

	sbtId, err := res.LastInsertId()
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		errTx = err
		return
	}

	scanTaskData, err := db.GetTbScanTaskModel(o.W("event_type", "Deployment"), o.W("chain_id", params.TokenChainId))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		errTx = err
		return
	}

	message := fmt.Sprintf(
		"%s%s%s%s%s%s%s",
		strings.TrimPrefix(user.Account, "0x"),
		fmt.Sprintf("%064x", params.TokenChainId),
		strings.TrimPrefix(scanTaskData.Address, "0x"),
		fmt.Sprintf("%064x", sbtId),
		params.ItemName,
		params.ItemName,
		params.FileUrl,
	)
	signature, err := utils.SignMessage(message, viper.GetString("app.sign_message_pri_key"))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SignMessage err: %v", err)
		errTx = err
		return
	}

	var meta = models.MetaData{
		Description: fmt.Sprintf("The SBT for %s dao only represents a status symbol, cannot be transferred, and has no financial attributes.", daoData.Handle),
		ExternalUrl: "https://www.myclique.io/daos",
		Image:       params.FileUrl,
		Name:        params.ItemName,
	}
	metaStr, err := json.Marshal(meta)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("json.Marshal err: %v", err)
		errTx = err
		return
	}

	resp := models.ResSBTCreate{
		Signature: signature,
		Meta:      string(metaStr),
	}

	jsonData(c, resp)
}

// @Summary sbt list
// @Tags sbt
// @version 0.0.1
// @description sbt list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param chainId query int true "token chainId"
// @Param status query string false "status:soon;active;ended"
// @Success 200 {object} models.ResSBTList
// @Router /stpdao/v2/sbt/list [get]
func SBTList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	chainId := c.Query("chainId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	chainIdParam, _ := strconv.ParseInt(chainId, 10, 64)
	statusParam := c.Query("status")

	var wChain, wStatus [][]interface{}
	if chainIdParam != 0 {
		wChain = o.W("token_chain_id", chainIdParam)
	}
	if statusParam != "" {
		wStatus = o.W("status", statusParam)
	}

	order := fmt.Sprintf("create_time ASC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbSBT(order, page, wChain, wStatus)
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}

// @Summary sbt detail
// @Tags sbt
// @version 0.0.1
// @description sbt detail
// @Produce json
// @Success 200 {object} models.ResSBTDetail
// @Router /stpdao/v2/sbt/detail/:sbtId [get]
func SBTDetail(c *gin.Context) {
	sbtIdParam := c.Param("sbtId")

	sbtData, err := db.GetTbSBT(o.W("id", sbtIdParam))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		return
	}

	dao, err := db.GetTbDao(o.W("chain_id", sbtData.ChainId), o.W("dao_address", sbtData.DaoAddress))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		return
	}

	resp := models.ResSBTDetail{
		ChainId:      sbtData.ChainId,
		DaoAddress:   sbtData.DaoAddress,
		DaoName:      dao.DaoName,
		DaoLogo:      dao.DaoLogo,
		TokenChainId: sbtData.TokenChainId,
		TokenAddress: sbtData.TokenAddress,
		FileUrl:      sbtData.FileUrl,
		ItemName:     sbtData.ItemName,
		Introduction: sbtData.Introduction,
		TotalSupply:  sbtData.TotalSupply,
		Way:          sbtData.Way,
		StartTime:    sbtData.StartTime,
		EndTime:      sbtData.EndTime,
		Status:       sbtData.Status,
	}

	jsonData(c, resp)
}

// @Summary get myself can claim sbt
// @Tags sbt
// @version 0.0.1
// @description get myself can claim sbt, request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Success 200 {object} models.ResSBTClaimInfo
// @Router /stpdao/v2/sbt/claim/:sbtId [get]
func SBTCanClaim(c *gin.Context) {
	var ok bool
	var user *db.TbAccountModel
	user, ok = parseJWTCache(c)
	if !ok {
		return
	}

	sbtId := c.Param("sbtId")
	sbtIdParam, _ := strconv.ParseInt(sbtId, 10, 64)

	sbtData, err := db.GetTbSBT(o.W("id", sbtIdParam))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v", err)
		return
	}

	var canClaim bool
	if sbtData.Way == consts.SBT_WAY_anyone {
		canClaim = true

	} else if sbtData.Way == consts.SBT_WAY_joined {
		count, err := o.Count(consts.TbJobs, o.W("chain_id", sbtData.ChainId),
			o.W("dao_address", sbtData.DaoAddress), o.W("account", user.Account))
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SQL err: %v", err)
			return
		}
		if count > 0 {
			canClaim = true
		}

	} else if sbtData.Way == consts.SBT_WAY_whitelist {
		var data models.JsonWhitelist
		err = json.Unmarshal([]byte(sbtData.WhiteList), &data)
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("json.Unmarshal err: %v", err)
			return
		}

		for _, val := range data.Account {
			if strings.EqualFold(val, user.Account) {
				canClaim = true
			}
		}
	}

	var signature string
	if canClaim {
		scanTaskData, err1 := db.GetTbScanTaskModel(o.W("event_type", "Deployment"), o.W("chain_id", sbtData.TokenChainId))
		if handleErrorIfExists(c, err1, errs.ErrServer) {
			oo.LogW("SQL err: %v", err1)
			return
		}

		message := fmt.Sprintf(
			"%s%s%s%s",
			fmt.Sprintf("%064x", sbtIdParam),
			fmt.Sprintf("%064x", sbtData.TokenChainId),
			strings.TrimPrefix(scanTaskData.Address, "0x"),
			strings.TrimPrefix(user.Account, "0x"),
		)
		signature, err = utils.SignMessage(message, viper.GetString("app.sign_message_pri_key"))
		if handleErrorIfExists(c, err, errs.ErrServer) {
			oo.LogW("SignMessage err: %v", err)
			return
		}
	}

	resp := models.ResSBTClaimInfo{
		CanClaim:  canClaim,
		Signature: signature,
	}

	jsonData(c, resp)
}

// @Summary sbt claim list
// @Tags sbt
// @version 0.0.1
// @description sbt claim list
// @Produce json
// @Param offset query int true "offset,page"
// @Param limit query int true "limit,page"
// @Param sbtId query int true "sbtId"
// @Success 200 {object} models.ResSBTClaimList
// @Router /stpdao/v2/sbt/claim/list [get]
func SBTClaimList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	sbtIdParam := c.Query("sbtId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)

	// consts.TbSBTClaim AS a	consts.TbNameAccount AS b
	order := fmt.Sprintf("a.create_time DESC")
	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}
	list, total, err := PageTbSBTClaim(order, page, o.W("sbt_id", sbtIdParam))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err:%v", err)
		return
	}

	jsonPagination(c, list, total, page)
}
