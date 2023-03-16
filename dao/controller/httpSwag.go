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
)

// @Summary create sale
// @Tags swap
// @version 0.0.1
// @description create sale
// @Produce json
// @Param request body models.ReqCreateSale true "request"
// @Success 200 {object} models.ResCreateSale
// @Router /stpdao/v2/swag/create [post]
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

		inAmountD := outAmountD.Mul(salePriceD).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(receivePriceData.Decimals)))

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
		discount = dis.String()[:10]
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
	v["white_list"] = params.WhiteList
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
			Status:           ls.Status,
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
