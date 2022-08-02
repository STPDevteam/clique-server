package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"strings"
)

// @Summary sign
// @Tags sign
// @version 0.0.1
// @description sign
// @Produce json
// @Param request body models.SignCreateDataParam true "request"
// @Success 200 {object} models.ResSignCreateData
// @Router /stpdao/v2/sign/create [post]
func (svc *Service) httpCreateSign(c *gin.Context) {
	var params models.SignCreateDataParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	daoAddress := utils.FixTo0x64String(params.DaoAddress)
	var createDaoEntity []models.EventHistoricalModel
	sqler := oo.NewSqler().Table(consts.TbNameEventHistorical).
		Where("event_type", consts.EvCreateDao).
		Where("topic2", daoAddress).Select()
	err = oo.SqlSelect(sqler, &createDaoEntity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	chainIdAndTokenAddress := strings.TrimPrefix(createDaoEntity[0].Data, "0x")
	ResChainId := chainIdAndTokenAddress[:64]

	var nonceEntity []models.NonceModel
	sqlSel := oo.NewSqler().Table(consts.TbNameNonce).Where("account", params.Account).Select()
	err = oo.SqlSelect(sqlSel, &nonceEntity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	resNonce := fmt.Sprintf("%064s", strings.TrimPrefix(strconv.FormatInt(int64(nonceEntity[0].Nonce), 16), "0x"))

	tokenAddress := utils.FixTo0xString(chainIdAndTokenAddress[64:128])
	tokenChainId, _ := strconv.ParseInt(chainIdAndTokenAddress[:64], 16, 64)
	var url string

	for indexScan := range svc.scanInfo {
		for indexUrl := range svc.scanInfo[indexScan].ChainId {

			chainId := svc.scanInfo[indexScan].ChainId[indexUrl]
			if int64(chainId) == tokenChainId {
				url = svc.scanInfo[indexScan].ScanUrl[indexUrl]
			}
		}
	}
	if url == "" || len(url) == 0 {
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Unsupported token.",
		})
		return
	}

	var balance string
	if params.SignType == "0" {
		res, errQ := utils.QueryBalance(tokenAddress, params.Account, url)
		if errQ != nil {
			oo.LogW("DoPost err: %v", errQ)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		balance = res.Result.Value

	} else if params.SignType == "1" {
		topic2 := fmt.Sprintf("0x%064s", strings.TrimPrefix(params.Account, "0x"))
		var VoteEntity []models.EventHistoricalModel
		sqlVote := oo.NewSqler().Table(consts.TbNameEventHistorical).
			Where("event_type", consts.EvCreateProposal).
			Where("address", daoAddress).
			Where("topic2", topic2).Select()
		err = oo.SqlSelect(sqlVote, &VoteEntity)
		if err != nil || VoteEntity == nil || len(VoteEntity) == 0 {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		blockNumber, _ := strconv.ParseInt(VoteEntity[0].BlockNumber, 16, 64)
		res, errQ := utils.QuerySpecifyBalance(tokenAddress, params.Account, url, blockNumber)
		if errQ != nil {
			oo.LogW("DoPost err: %v", errQ)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		balance = res.Result.Value

	}
	decBalance, _ := new(big.Int).SetString(balance, 10)
	resBalance := fmt.Sprintf("%064s", fmt.Sprintf("%x", decBalance))

	resTokenAddress := strings.TrimPrefix(tokenAddress, "0x")
	resSignType := fmt.Sprintf("%064s", params.SignType)
	resAccount := strings.TrimPrefix(params.Account, "0x")

	message := fmt.Sprintf("%s%s%s%s%s%s", resAccount, resNonce, ResChainId, resTokenAddress, resBalance, resSignType)

	signature, err := utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "Signature err",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSignCreateData{
			TokenChainId: tokenChainId,
			TokenAddress: tokenAddress,
			Balance:      balance,
			Signature:    signature,
		},
	})
}
