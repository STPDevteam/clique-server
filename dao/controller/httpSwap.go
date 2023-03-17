package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"math/big"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"strings"
	"time"
)

// @Summary create sale
// @Tags swap
// @version 0.0.1
// @description create sale
// @Produce json
// @Param request body models.ReqCreateSale true "request"
// @Success 200 {object} models.ResCreateSale
// @Router /stpdao/v2/swap/create [post]
func (svc *Service) createSwap(c *gin.Context) {
	var params models.ReqCreateSale
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var salePriceData models.TbSwapToken
	sqlSel := oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", params.ChainId).Where("token_address", params.SaleToken).Select()
	err = oo.SqlGet(sqlSel, &salePriceData)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var receivePriceData models.TbSwapToken
	sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", params.ChainId).Where("token_address", params.ReceiveToken).Select()
	err = oo.SqlGet(sqlSel, &receivePriceData)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var discount string
	if params.SaleWay == "discount" {
		outAmountD, err1 := decimal.NewFromString(params.SaleAmount)
		salePriceD, err2 := decimal.NewFromString(params.SalePrice)
		if err1 != nil || err2 != nil {
			oo.LogW("decimal.NewFromString err1: %v, err2: %v", err1, err2)
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusBadRequest,
				Message: "Invalid parameters.",
			})
			return
		}

		e := decimal.NewFromInt(10).Pow(decimal.NewFromInt(receivePriceData.Decimals))
		inAmountD := outAmountD.Mul(salePriceD).Div(e)

		modelTo := inAmountD.Mul(decimal.NewFromFloat(receivePriceData.Price))
		modelFrom := outAmountD.Mul(decimal.NewFromFloat(salePriceData.Price))

		dis := modelTo.Div(modelFrom)
		if !dis.LessThanOrEqual(decimal.NewFromFloat(0.9)) {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("At least 10%% off is required, now is %s", discount),
			})
			return
		}
		discount = dis.String()
	}

	var lastId int64
	var signature string
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	defer func() {
		oo.CloseSqlxTx(tx, &errTx)
		if errTx == nil {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusOK,
				Message: "ok",
				Data: models.ResCreateSale{
					SaleId:    lastId,
					Signature: signature,
				},
			})
		}
	}()

	var whiteStr string
	for i := range params.WhiteList {
		if whiteStr == "" {
			whiteStr = params.WhiteList[i]
		} else {
			whiteStr = fmt.Sprintf("%s,%s", whiteStr, params.WhiteList[i])
		}
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["creator"] = params.Creator
	v["sale_way"] = params.SaleWay
	v["sale_token"] = params.SaleToken
	v["sale_token_img"] = salePriceData.Img
	v["sale_amount"] = params.SaleAmount
	v["sale_price"] = params.SalePrice
	v["original_discount"] = discount
	v["receive_token"] = params.ReceiveToken
	v["receive_token_img"] = receivePriceData.Img
	v["limit_min"] = params.LimitMin
	v["limit_max"] = params.LimitMax
	v["start_time"] = params.StartTime
	v["end_time"] = params.EndTime
	v["white_list"] = whiteStr
	v["about"] = params.About
	m = append(m, v)
	sqlIns := oo.NewSqler().Table(consts.TbNameSwap).Insert(m)
	res, errTx := oo.SqlxTxExec(tx, sqlIns)
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	lastId, errTx = res.LastInsertId()
	if errTx != nil {
		oo.LogW("SQL err: %v", errTx)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	saleAmount, _ := new(big.Int).SetString(params.SaleAmount, 10)
	salePrice, _ := new(big.Int).SetString(params.SalePrice, 10)
	message := fmt.Sprintf(
		"%s%s%s%s%s%s%s%s%s%s",
		strings.TrimPrefix(params.Creator, "0x"),
		fmt.Sprintf("%064x", lastId),
		strings.TrimPrefix(params.SaleToken, "0x"),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", saleAmount)),
		strings.TrimPrefix(params.ReceiveToken, "0x"),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", salePrice)),
		fmt.Sprintf("%064x", params.LimitMin),
		fmt.Sprintf("%064x", params.LimitMax),
		fmt.Sprintf("%064x", params.StartTime),
		fmt.Sprintf("%064x", params.EndTime),
	)
	oo.LogW("create sale sign message: %s", message)
	signature, errTx = utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if errTx != nil {
		oo.LogW("utils.SignMessage err: %v", errTx)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
}

// @Summary purchased sale
// @Tags swap
// @version 0.0.1
// @description purchased sale
// @Produce json
// @Param request body models.ReqPurchased true "request"
// @Success 200 {object} models.ResPurchased
// @Router /stpdao/v2/swap/purchased [post]
func (svc *Service) purchasedSwap(c *gin.Context) {
	var params models.ReqPurchased
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var swapData models.TbSwap
	sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", params.SaleId).Select()
	err = oo.SqlGet(sqlSel, &swapData)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var isWhite bool
	if swapData.WhiteList != "" {
		arr := strings.Split(swapData.WhiteList, ",")
		for i := range arr {
			if strings.ToLower(arr[i]) == strings.ToLower(params.Account) {
				isWhite = true
				break
			}
		}
	}
	if !isWhite {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Not whitelist.",
		})
		return
	}

	if swapData.Status != consts.StatusNormal || swapData.EndTime <= time.Now().Unix() {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Swap is ended.",
		})
		return
	}

	buyAmountD, err1 := decimal.NewFromString(params.BuyAmount)
	saleAmountD, err2 := decimal.NewFromString(swapData.SaleAmount)
	soleAmountD, err3 := decimal.NewFromString(swapData.SoldAmount)
	limitMinD, err4 := decimal.NewFromString(swapData.LimitMin)
	limitMaxD, err5 := decimal.NewFromString(swapData.LimitMax)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if !buyAmountD.LessThanOrEqual(saleAmountD.Sub(soleAmountD)) || !buyAmountD.LessThanOrEqual(limitMaxD) || !buyAmountD.GreaterThanOrEqual(limitMinD) {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Not enough balance.",
		})
		return
	}

	buyAmount, _ := new(big.Int).SetString(params.BuyAmount, 10)
	message := fmt.Sprintf(
		"%s%s%s",
		strings.TrimPrefix(params.Account, "0x"),
		fmt.Sprintf("%064x", params.SaleId),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", buyAmount)),
	)
	oo.LogW("purchased sign message: %s", message)
	signature, err := utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if err != nil {
		oo.LogW("utils.SignMessage err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResPurchased{
			Signature: signature,
		},
	})
}

// @Summary sale list
// @Tags swap
// @version 0.0.1
// @description sale list
// @Produce json
// @Param saleId query int false "saleId"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResSwapListPage
// @Router /stpdao/v2/swap/list [get]
func swapList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	saleId := c.Query("saleId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	saleIdParam, _ := strconv.Atoi(saleId)

	var swapArr []models.TbSwap
	sqler := oo.NewSqler().Table(consts.TbNameSwap).Where("status", "!=", consts.StatusPending)

	if saleId != "" {
		sqler.Where("id", saleIdParam)
	}

	var total int64
	sqlCopy := *sqler
	sqlSel := sqlCopy.Count()
	err := oo.SqlGet(sqlSel, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlSel = sqlCopy.Order("status DESC,create_time DESC").Limit(limitParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlSel, &swapArr)
	}
	if err != nil {
		oo.LogW("SQL err: %v, str: %s", err, sqlSel)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResSwapList, 0)
	for index := range swapArr {
		ls := swapArr[index]

		var status string
		if ls.EndTime < time.Now().Unix() {
			status = consts.StatusEnded
		} else {
			status = ls.Status
		}

		data = append(data, models.ResSwapList{
			SaleId:           ls.Id,
			ChainId:          ls.ChainId,
			Creator:          ls.Creator,
			SaleToken:        ls.SaleToken,
			SaleTokenImg:     ls.SaleTokenImg,
			SaleAmount:       ls.SaleAmount,
			SalePrice:        ls.SalePrice,
			ReceiveToken:     ls.ReceiveToken,
			ReceiveTokenImg:  ls.ReceiveTokenImg,
			LimitMin:         ls.LimitMin,
			LimitMax:         ls.LimitMax,
			StartTime:        ls.StartTime,
			EndTime:          ls.EndTime,
			Status:           status,
			About:            ls.About,
			OriginalDiscount: ls.OriginalDiscount,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSwapListPage{
			List:  data,
			Total: total,
		},
	})
}

// @Summary transactions list
// @Tags swap
// @version 0.0.1
// @description transactions list
// @Produce json
// @Param saleId query int true "saleId"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResSwapTransactionListPage
// @Router /stpdao/v2/swap/transactions [get]
func swapTransactionsList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	saleId := c.Query("saleId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	saleIdParam, _ := strconv.Atoi(saleId)

	var transactionArr []models.TbSwapTransaction
	sqler := oo.NewSqler().Table(consts.TbNameSwapTransaction).Where("sale_id", saleIdParam)

	var total int64
	sqlCopy := *sqler
	sqlSel := sqlCopy.Count()
	err := oo.SqlGet(sqlSel, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlSel = sqlCopy.Order("create_time DESC").Limit(limitParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlSel, &transactionArr)
	}
	if err != nil {
		oo.LogW("SQL err: %v, str: %s", err, sqlSel)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResSwapTransactionList, 0)
	for index := range transactionArr {
		ls := transactionArr[index]

		var buyToken models.TbSwapToken
		sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", ls.ChainId).Where("token_address", ls.BuyToken).Select()
		err = oo.SqlGet(sqlSel, &buyToken)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var payToken models.TbSwapToken
		sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", ls.ChainId).Where("token_address", ls.PayToken).Select()
		err = oo.SqlGet(sqlSel, &payToken)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		data = append(data, models.ResSwapTransactionList{
			SaleId:       ls.Id,
			Buyer:        ls.Buyer,
			BuyAmount:    ls.BuyAmount,
			PayAmount:    ls.PayAmount,
			Time:         ls.Time,
			BuyTokenName: buyToken.TokenName,
			PayTokenName: payToken.TokenName,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSwapTransactionListPage{
			List:  data,
			Total: total,
		},
	})
}
