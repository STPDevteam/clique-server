package handlers

import (
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
)

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
	if !daoData.Approve {
		handleError(c, errs.NewError(403, "Dao not approve."))
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

	message := fmt.Sprintf("%d", sbtId)
	signature, err := utils.SignMessage(message, viper.GetString("app.sign_message_pri_key"))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SignMessage err: %v", err)
		errTx = err
		return
	}

	jsonData(c, signature)
}

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
