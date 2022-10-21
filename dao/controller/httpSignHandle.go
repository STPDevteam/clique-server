package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"regexp"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"strings"
	"time"
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
		Where("topic3", daoAddress).
		Where("chain_id", params.ChainId).Select()
	err = oo.SqlSelect(sqler, &createDaoEntity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.1",
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
			Message: "Something went wrong, Please try again later.2",
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
			if chainId == int(tokenChainId) {
				url = svc.scanInfo[indexScan].ScanUrl[indexUrl]
				break
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
		data := fmt.Sprintf("%s%s", paramsDataPrefix, strings.TrimPrefix(params.Account, "0x"))
		res, errQb := utils.QueryMethodEthCall(tokenAddress, data, url)
		if errQb != nil || res.Result == nil || res.Result == "0x" {
			oo.LogW("DoPost err: %v, tokenAddress:%v, data:%v, url:%v", errQb, tokenAddress, data, url)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.3",
			})
			return
		}
		resBalance = strings.TrimPrefix(res.Result.(string), "0x")

	} else if params.SignType == "1" {
		//var voteEntity models.EventHistoricalModel
		//proposalId := utils.FixTo0x64String(fmt.Sprintf(`%x`, params.ProposalId))
		//sqlVote := oo.NewSqler().Table(consts.TbNameEventHistorical).
		//	Where("event_type", consts.EvCreateProposal).
		//	Where("address", params.DaoAddress).
		//	Where("chain_id", params.ChainId).
		//	Where("topic1", proposalId).Select()
		//err = oo.SqlGet(sqlVote, &voteEntity)
		//if err != nil {
		//	oo.LogW("SQL err: %v", err)
		//	c.JSON(http.StatusInternalServerError, models.Response{
		//		Code:    500,
		//		Message: "Something went wrong, Please try again later.4",
		//	})
		//	return
		//}

		var blockNumber string
		sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("chain_id", params.ChainId).
			Where("dao_address", params.DaoAddress).Where("proposal_id", params.ProposalId).Select("block_number")
		err = oo.SqlGet(sqlSel, &blockNumber)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.4",
			})
			return
		}

		var success = false
		key := fmt.Sprintf(`%d-%s-%s-%s`, tokenChainId, tokenAddress, params.Account, blockNumber)
		cacheBalance, ok := svc.mCache.Get(key)
		if !ok {
			for _, testChainId := range svc.appConfig.ArchiveBalanceSign {
				if tokenChainId == int64(testChainId) {
					url = svc.getArchiveNode(tokenChainId)
					if url == "" {
						c.JSON(http.StatusInternalServerError, models.Response{
							Code:    500,
							Message: "Unsupported token.",
						})
						return
					}
					data := fmt.Sprintf("%s%s", paramsDataPrefix, strings.TrimPrefix(params.Account, "0x"))
					var tag string
					for _, testnet := range svc.appConfig.TestnetBalanceSign {
						if testnet == testChainId {
							tag = "latest"
							break
						}
						tag = blockNumber
					}
					res, errQb := utils.QueryMethodEthCallByTag(tokenAddress, data, url, tag)
					if errQb != nil || res.Result == nil || res.Result == "0x" {
						oo.LogW("DoPost err: %v, tokenAddress:%v, data:%v, url:%v, tag:%v", errQb, tokenAddress, data, url, tag)
						c.JSON(http.StatusInternalServerError, models.Response{
							Code:    500,
							Message: "Something went wrong, Please try again later.5",
						})
						return
					}

					resBalance = strings.TrimPrefix(res.Result.(string), "0x")
					success = true
					break
				}
			}
			if !success {
				blockNumberDec, _ := strconv.ParseInt(blockNumber, 16, 64)
				if !ok {
					res, errQ := utils.QuerySpecifyBalance(tokenAddress, params.Account, url, blockNumberDec)
					if errQ != nil {
						oo.LogW("DoPost err: %v, tokenAddress:%v, data:%v, url:%v, blockNumber:%v", errQ, tokenAddress, params.Account, url, blockNumber)
						c.JSON(http.StatusInternalServerError, models.Response{
							Code:    500,
							Message: "Something went wrong, Please try again later.6",
						})
						return
					}
					decBalance, _ := new(big.Int).SetString(res.Result.Value, 10)
					resBalance = fmt.Sprintf("%064s", fmt.Sprintf("%x", decBalance))
				}
			}

			svc.mCache.Set(key, resBalance, time.Duration(72)*time.Hour)
		} else {
			resBalance = cacheBalance.(string)
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

// @Summary sign lock dao handle
// @Tags sign
// @version 0.0.1
// @description sign lock dao handle
// @Produce json
// @Param request body models.SignDaoHandleParam true "request"
// @Success 200 {object} models.ResSignDaoHandleData
// @Router /stpdao/v2/sign/lock/handle [post]
func (svc *Service) httpLockDaoHandleSign(c *gin.Context) {
	var params models.SignDaoHandleParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}
	params.Handle = strings.Replace(params.Handle, " ", "", -1)
	r, _ := regexp.Compile("^[0-9a-z_]*$")
	if !r.MatchString(params.Handle) {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Handle Invalid parameters.",
		})
		return
	}

	resAccount := strings.TrimPrefix(params.Account, "0x")
	resChainId := fmt.Sprintf("%064x", params.ChainId)
	resHandle := utils.Keccak256(params.Handle)

	var resBlock string
	var latestBlockNum int
	for indexScan := range svc.scanInfo {
		for indexUrl := range svc.scanInfo[indexScan].ChainId {
			url := svc.scanInfo[indexScan].ScanUrl[indexUrl]
			chainId := svc.scanInfo[indexScan].ChainId[indexUrl]
			if chainId == params.ChainId {
				res, err := utils.QueryLatestBlock(url)
				if err != nil || res.Result.(string) == "" {
					oo.LogW("QueryLatestBlock err: %v", err)
					c.JSON(http.StatusInternalServerError, models.Response{
						Code:    500,
						Message: "Something went wrong, Please try again later.",
					})
					return
				}
				resResultBlock, _ := utils.Hex2Dec(res.Result.(string))
				latestBlockNum = resResultBlock + svc.scanInfo[indexScan].HandleLockBlock[indexUrl]
				resBlock = fmt.Sprintf("%064x", latestBlockNum)
				break
			}
		}
	}

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameHandleLock).Where("handle", params.Handle).Count()
	err = oo.SqlGet(sqlSel, &count)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	if count != 0 {
		var handleEntity []models.HandleLockModel
		sqlSel = oo.NewSqler().Table(consts.TbNameHandleLock).
			Where("handle", params.Handle).Where("account", params.Account).Where("chain_id", params.ChainId).
			Where("lock_block", "!=", consts.MaxIntUnsigned).Select()
		err = oo.SqlSelect(sqlSel, &handleEntity)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if len(handleEntity) == 0 {
			c.JSON(http.StatusOK, models.Response{
				Code:    200,
				Message: "handle locked or existed",
			})
			return
		} else {
			resOldAccount := strings.TrimPrefix(handleEntity[0].Account, "0x")
			resOldChainId := fmt.Sprintf("%064x", handleEntity[0].ChainId)
			resOldBlock := fmt.Sprintf("%064x", handleEntity[0].LockBlock)

			messageOld := fmt.Sprintf("%s%s%s%s", resOldAccount, resOldChainId, resOldBlock, handleEntity[0].HandleKeccak)
			signatureOld, err := utils.SignMessage(messageOld, svc.appConfig.SignMessagePriKey)
			if err != nil {
				oo.LogW("SignMessage err: %v", err)
				c.JSON(http.StatusInternalServerError, models.Response{
					Code:    500,
					Message: "Signature err",
				})
				return
			}
			signatureOld = fmt.Sprintf("0x%s", signatureOld)
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusOK,
				Message: "ok",
				Data: models.ResSignDaoHandleData{
					Signature:    signatureOld,
					Account:      handleEntity[0].Account,
					ChainId:      handleEntity[0].ChainId,
					LockBlockNum: handleEntity[0].LockBlock,
				},
			})
			return
		}

	} else {
		var m = make([]map[string]interface{}, 0)
		var v = make(map[string]interface{})
		v["handle"] = params.Handle
		v["handle_keccak"] = resHandle
		v["lock_block"] = latestBlockNum
		v["chain_id"] = params.ChainId
		v["account"] = params.Account
		m = append(m, v)
		sqlIns := oo.NewSqler().Table(consts.TbNameHandleLock).Insert(m)
		err = oo.SqlExec(sqlIns)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
	}

	message := fmt.Sprintf("%s%s%s%s", resAccount, resChainId, resBlock, resHandle)
	signature, err := utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if err != nil {
		oo.LogW("SignMessage err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Signature err",
		})
		return
	}
	signature = fmt.Sprintf("0x%s", signature)

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSignDaoHandleData{
			Signature:    signature,
			Account:      params.Account,
			ChainId:      params.ChainId,
			LockBlockNum: latestBlockNum,
		},
	})

}

// @Summary sign query dao handle
// @Tags sign
// @version 0.0.1
// @description sign query dao handle
// @Produce json
// @Param handle query string true "handle"
// @Param chainId query int true "chainId"
// @Param account query string true "account"
// @Success 200 {object} models.ResResult
// @Router /stpdao/v2/sign/query/handle [get]
func (svc *Service) httpQueryDaoHandle(c *gin.Context) {
	handleParam := c.Query("handle")
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	accountParam := c.Query("account")

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameHandleLock).Where("handle", handleParam).Count()
	err := oo.SqlGet(sqlSel, &count)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if count != 0 {
		var countOwn int
		sqlSel = oo.NewSqler().Table(consts.TbNameHandleLock).
			Where("account", accountParam).
			Where("chain_id", chainIdParam).
			Where("handle", handleParam).
			Where("lock_block", "!=", consts.MaxIntUnsigned).Count()
		err = oo.SqlGet(sqlSel, &countOwn)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if countOwn == 0 {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusOK,
				Message: "ok",
				Data: models.ResResult{
					Success: false,
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResResult{
			Success: true,
		},
	})

}

func (svc *Service) getArchiveNode(chainId int64) string {
	if chainId == 1 {
		return svc.appConfig.MainnetChainstackRPC
	}
	if chainId == 137 {
		return svc.appConfig.PolygonQuickNodeRPC
	}
	if chainId == 80001 {
		return "https://rpc.ankr.com/polygon_mumbai"
	}
	if chainId == 5 {
		return "https://rpc.ankr.com/eth_goerli"
	}
	return ""
}
