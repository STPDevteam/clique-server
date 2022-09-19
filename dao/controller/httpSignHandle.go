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
		Where("topic2", daoAddress).
		Where("chain_id", params.ChainId).Select()
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
	resChainId := chainIdAndTokenAddress[:64]

	var nonceEntity []models.NonceModel
	sqlSel := oo.NewSqler().Table(consts.TbNameNonce).
		Where("chain_id", params.ChainId).
		Where("account", params.Account).Select()
	err = oo.SqlSelect(sqlSel, &nonceEntity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var resNonce string
	if len(nonceEntity) == 0 {
		resNonce = "0000000000000000000000000000000000000000000000000000000000000000"
	} else {
		resNonce = fmt.Sprintf("%064s", strings.TrimPrefix(strconv.FormatInt(int64(nonceEntity[0].Nonce), 16), "0x"))
	}

	tokenAddress := utils.FixTo0x40String(chainIdAndTokenAddress[64:128])
	tokenChainId, _ := strconv.ParseInt(chainIdAndTokenAddress[:64], 16, 64)

	var url string
	for indexScan := range svc.scanInfo {
		for indexUrl := range svc.scanInfo[indexScan].ChainId {
			chainId := svc.scanInfo[indexScan].ChainId[indexUrl]
			if chainId == params.ChainId {
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

	var resBalance string
	const paramsDataPrefix = "0x70a08231000000000000000000000000"
	if params.SignType == "0" {
		//res, errQ := utils.QueryBalance(tokenAddress, params.Account, url)
		//if errQ != nil {
		//	oo.LogW("DoPost err: %v", errQ)
		//	c.JSON(http.StatusInternalServerError, models.Response{
		//		Code:    500,
		//		Message: "Something went wrong, Please try again later.",
		//	})
		//	return
		//}
		//balance = res.Result.Value
		data := fmt.Sprintf("%s%s", paramsDataPrefix, strings.TrimPrefix(params.Account, "0x"))
		res, errQb := utils.QueryMethodEthCall(tokenAddress, data, url)
		if errQb != nil || res.Result == nil || res.Result == "0x" {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		resBalance = strings.TrimPrefix(res.Result.(string), "0x")

	} else if params.SignType == "1" {

		for _, mainChainId := range svc.appConfig.MainnetBalanceSign {
			if params.ChainId == mainChainId {
				var VoteEntity []models.EventHistoricalModel
				proposalId := utils.FixTo0x64String(fmt.Sprintf(`%x`, params.ProposalId))
				sqlVote := oo.NewSqler().Table(consts.TbNameEventHistorical).
					Where("event_type", consts.EvCreateProposal).
					Where("address", params.DaoAddress).
					Where("chain_id", params.ChainId).
					Where("topic1", proposalId).Select()
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
				decBalance, _ := new(big.Int).SetString(res.Result.Value, 10)
				resBalance = fmt.Sprintf("%064s", fmt.Sprintf("%x", decBalance))
			}
		}

		for _, testChainId := range svc.appConfig.TestnetBalanceSign {
			if params.ChainId == testChainId {
				data := fmt.Sprintf("%s%s", paramsDataPrefix, strings.TrimPrefix(params.Account, "0x"))
				res, errQb := utils.QueryMethodEthCall(tokenAddress, data, url)
				if errQb != nil || res.Result == nil || res.Result == "0x" {
					c.JSON(http.StatusInternalServerError, models.Response{
						Code:    500,
						Message: "Something went wrong, Please try again later.",
					})
					return
				}
				resBalance = strings.TrimPrefix(res.Result.(string), "0x")
			}
		}

	}
	balance, _ := utils.Hex2BigInt(fmt.Sprintf("0x%s", resBalance))

	resTokenAddress := strings.TrimPrefix(tokenAddress, "0x")
	resSignType := fmt.Sprintf("%064s", params.SignType)
	resAccount := strings.TrimPrefix(params.Account, "0x")

	message := fmt.Sprintf("%s%s%s%s%s%s", resAccount, resNonce, resChainId, resTokenAddress, resBalance, resSignType)
	signature, err := utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "Signature err",
		})
		return
	}
	signature = fmt.Sprintf("0x%s", signature)

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSignCreateData{
			Account:      params.Account,
			TokenChainId: tokenChainId,
			TokenAddress: tokenAddress,
			Balance:      balance.String(),
			Signature:    signature,
		},
	})
}

func (svc *Service) httpDaoHandleSign(c *gin.Context) {

}
