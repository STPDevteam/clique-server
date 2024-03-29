package handlers

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"math/big"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
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
func CreateSwap(c *gin.Context) {
	var params models.ReqCreateSale
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var tbSysConfig db.TbSysConfig
	sqlSel := oo.NewSqler().Table(consts.TbSysConfig).Where("cfg_name", "cfg_swap_creator_white_list").Select()
	err = oo.SqlGet(sqlSel, &tbSysConfig)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var isApproved bool
	if tbSysConfig.CfgIsEnabled {
		sliceVal := strings.Split(tbSysConfig.CfgVal, ",")
		for _, v := range sliceVal {
			if strings.ToLower(v) == strings.ToLower(params.Creator) {
				isApproved = true
				break
			}
		}
		if !isApproved {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusBadRequest,
				Message: "Not creator whitelist.",
			})
			return
		}
	}

	var salePriceData db.TbSwapToken
	sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", params.ChainId).Where("token_address", params.SaleToken).Select()
	err = oo.SqlGet(sqlSel, &salePriceData)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var receivePriceData db.TbSwapToken
	sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", params.ChainId).Where("token_address", params.ReceiveToken).Select()
	err = oo.SqlGet(sqlSel, &receivePriceData)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var discount string
	//if params.SaleWay == "discount" {
	outAmountD, err1 := decimal.NewFromString(params.SaleAmount)
	salePriceD, err2 := decimal.NewFromString(params.SalePrice)
	limitMinD, err3 := decimal.NewFromString(params.LimitMin)
	if err1 != nil || err2 != nil || err3 != nil {
		oo.LogW("decimal.NewFromString err1: %v, err2: %v", err1, err2)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}
	if outAmountD.LessThan(limitMinD) {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "limit min more than the sale amount.",
		})
		return
	}

	e := decimal.NewFromInt(10).Pow(decimal.NewFromInt(receivePriceData.Decimals))
	inAmountD := outAmountD.Mul(salePriceD).Div(e)

	modelTo := inAmountD.Mul(decimal.NewFromFloat(receivePriceData.Price))
	modelFrom := outAmountD.Mul(decimal.NewFromFloat(salePriceData.Price))

	dis := modelTo.Div(modelFrom)
	//if !dis.LessThanOrEqual(decimal.NewFromFloat(2)) {
	//	oo.LogW(fmt.Sprintf("Cannot increase the price by more than 2 times, now is %s", dis.String()))
	//	c.JSON(http.StatusOK, models.Response{
	//		Code:    http.StatusBadRequest,
	//		Message: fmt.Sprintf("Cannot increase the price by more than 2 times"),
	//	})
	//	return
	//}
	discount = dis.String()
	//}

	var lastId int64
	var signature string
	tx, errTx := oo.NewSqlxTx()
	if errTx != nil {
		oo.LogW("SQL err: %v", err)
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
					Signature: fmt.Sprintf("0x%s", signature),
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

	discountFloat, err := strconv.ParseFloat(discount, 64)
	if err != nil {
		oo.LogW("strconv.ParseFloat err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "discount calculation failed.",
		})
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["title"] = params.Title
	v["creator"] = params.Creator
	v["sale_way"] = params.SaleWay
	v["sale_token"] = params.SaleToken
	v["sale_token_img"] = salePriceData.Img
	v["sale_amount"] = params.SaleAmount
	v["sale_price"] = params.SalePrice
	v["original_discount"] = fmt.Sprintf("%.4f", discountFloat)
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
	limitMin, _ := new(big.Int).SetString(params.LimitMin, 10)
	limitMax, _ := new(big.Int).SetString(params.LimitMax, 10)
	message := fmt.Sprintf(
		"%s%s%s%s%s%s%s%s%s%s",
		strings.TrimPrefix(params.Creator, "0x"),
		fmt.Sprintf("%064x", lastId),
		strings.TrimPrefix(params.SaleToken, "0x"),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", saleAmount)),
		strings.TrimPrefix(params.ReceiveToken, "0x"),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", salePrice)),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", limitMin)),
		fmt.Sprintf("%064s", fmt.Sprintf("%x", limitMax)),
		fmt.Sprintf("%064x", params.StartTime),
		fmt.Sprintf("%064x", params.EndTime),
	)
	oo.LogW("create sale sign message: %s", message)
	signature, errTx = utils.SignMessage(message, viper.GetString("app.sign_message_pri_key"))
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
func PurchasedSwap(c *gin.Context) {
	var params models.ReqPurchased
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("params err:%v", params)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var swapData db.TbSwap
	sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", params.SaleId).Select()
	err = oo.SqlGet(sqlSel, &swapData)
	if err != nil {
		oo.LogW("SQL err: %v", err)
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
		if !isWhite {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusBadRequest,
				Message: "Not whitelist.",
			})
			return
		}
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
		oo.LogW("decimal.NewFromString err1: %v,err2: %v,err3: %v,err4: %v,err5: %v", err1, err2, err3, err4, err5)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if !buyAmountD.LessThanOrEqual(saleAmountD.Sub(soleAmountD)) {
		oo.LogW("buyAmountD: %s, saleAmountD: %s, soleAmountD: %s, limitMinD: %s, limitMaxD: %s", buyAmountD.String(), saleAmountD.String(), soleAmountD.String(), limitMinD.String(), limitMaxD.String())
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "not enough stock for swap.",
		})
		return
	}
	if !buyAmountD.LessThanOrEqual(limitMaxD) || !buyAmountD.GreaterThanOrEqual(limitMinD) {
		oo.LogW("buyAmountD: %s, saleAmountD: %s, soleAmountD: %s, limitMinD: %s, limitMaxD: %s", buyAmountD.String(), saleAmountD.String(), soleAmountD.String(), limitMinD.String(), limitMaxD.String())
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Not in sales limit.",
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
	signature, err := utils.SignMessage(message, viper.GetString("app.sign_message_pri_key"))
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
			Signature: fmt.Sprintf("0x%s", signature),
		},
	})
}

// @Summary sale list
// @Tags swap
// @version 0.0.1
// @description sale list
// @Produce json
// @Param status query string false "status: soon/normal/ended"
// @Param saleId query int false "saleId"
// @Param offset query  int true "offset,page"
// @Param limit query  int true "limit,page"
// @Success 200 {object} models.ResSwapListPage
// @Router /stpdao/v2/swap/list [get]
func SwapList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	saleId := c.Query("saleId")
	statusParam := c.Query("status")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	saleIdParam, _ := strconv.Atoi(saleId)

	var swapArr []db.TbSwap
	sqler := oo.NewSqler().Table(consts.TbNameSwap).Where("status", "!=", consts.StatusPending)

	if saleId != "" {
		sqler.Where("id", saleIdParam)
	}
	if statusParam != "" {
		sqler.Where("status", statusParam)
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

		createTime, _ := time.Parse("2006-01-02 15:04:05", ls.CreateTime)

		data = append(data, models.ResSwapList{
			SaleId:           ls.Id,
			SaleWay:          ls.SaleWay,
			Title:            ls.Title,
			CreateTime:       createTime.Unix(),
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
			SoldAmount:       ls.SoldAmount,
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
// @Param limit query  int true "limit,page"
// @Success 200 {object} models.ResSwapTransactionListPage
// @Router /stpdao/v2/swap/transactions [get]
func SwapTransactionsList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	saleId := c.Query("saleId")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)
	saleIdParam, _ := strconv.Atoi(saleId)

	var transactionArr []db.TbSwapTransaction
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

		var buyToken db.TbSwapToken
		sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", ls.ChainId).Where("token_address", ls.BuyToken).Select()
		err = oo.SqlGet(sqlSel, &buyToken)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var payToken db.TbSwapToken
		sqlSel = oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", ls.ChainId).Where("token_address", ls.PayToken).Select()
		err = oo.SqlGet(sqlSel, &payToken)
		if err != nil {
			oo.LogW("SQL err: %v", err)
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

// @Summary prices
// @Tags swap
// @version 0.0.1
// @description prices
// @Produce json
// @Param chainId query int true "chainId"
// @Param tokens query  string false "tokens, separate by comma, if nil, return all"
// @Success 200 {object} models.ResSwapPrices
// @Router /stpdao/v2/swap/prices [get]
func SwapPrices(c *gin.Context) {
	chainId := c.Query("chainId")
	tokensParam := c.Query("tokens")
	chainIdParam, _ := strconv.Atoi(chainId)

	var data = make([]models.ResSwapPrices, 0)
	if tokensParam == "" {
		var tokenList []db.TbSwapToken
		sqlSel := oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", chainIdParam).Select()
		err := oo.SqlSelect(sqlSel, &tokenList)
		if err != nil {
			oo.LogW("SQL err: %v, str: %s", err, sqlSel)
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		for i := range tokenList {
			ls := tokenList[i]

			updateAt, _ := time.Parse("2006-01-02 15:04:05", ls.UpdateTime)
			data = append(data, models.ResSwapPrices{
				ChainId:          ls.ChainId,
				TokenAddress:     ls.TokenAddress,
				Price:            ls.Price,
				Img:              ls.Img,
				UrlCoingecko:     ls.UrlCoingecko,
				UrlCoinmarketcap: ls.UrlCoinmarketcap,
				TokenName:        ls.TokenName,
				Symbol:           ls.Symbol,
				Decimals:         ls.Decimals,
				UpdateAt:         updateAt.Unix(),
			})
		}
	} else {
		tokens := strings.Split(tokensParam, ",")
		for _, token := range tokens {
			var tbToken db.TbSwapToken
			sqlSel := oo.NewSqler().Table(consts.TbNameSwapToken).Where("chain_id", chainIdParam).Where("token_address", token).Select()
			err := oo.SqlGet(sqlSel, &tbToken)
			if err != nil {
				oo.LogW("SQL err: %v, str: %s", err, sqlSel)
				c.JSON(http.StatusOK, models.Response{
					Code:    http.StatusInternalServerError,
					Message: "Something went wrong, Please try again later.",
				})
				return
			}

			updateAt, _ := time.Parse("2006-01-02 15:04:05", tbToken.UpdateTime)
			data = append(data, models.ResSwapPrices{
				ChainId:          tbToken.ChainId,
				TokenAddress:     tbToken.TokenAddress,
				Price:            tbToken.Price,
				Img:              tbToken.Img,
				UrlCoingecko:     tbToken.UrlCoingecko,
				UrlCoinmarketcap: tbToken.UrlCoinmarketcap,
				TokenName:        tbToken.TokenName,
				Symbol:           tbToken.Symbol,
				Decimals:         tbToken.Decimals,
				UpdateAt:         updateAt.Unix(),
			})
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})
}

// @Summary buy isWhite
// @Tags swap
// @version 0.0.1
// @description buy isWhite
// @Produce json
// @Param account query string true "account"
// @Param saleId query int true "saleId"
// @Success 200 {object} models.ResIsWhite
// @Router /stpdao/v2/swap/isWhite [get]
func SwapIsWhite(c *gin.Context) {
	saleId := c.Query("saleId")
	accountParam := c.Query("account")
	saleIdParam, _ := strconv.Atoi(saleId)

	var swapData db.TbSwap
	sqlSel := oo.NewSqler().Table(consts.TbNameSwap).Where("id", saleIdParam).Select()
	err := oo.SqlGet(sqlSel, &swapData)
	if err != nil {
		oo.LogW("SQL err: %v", err)
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
			if strings.ToLower(arr[i]) == strings.ToLower(accountParam) {
				isWhite = true
				break
			}
		}
	} else {
		isWhite = true
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResIsWhite{
			IsWhite: isWhite,
		},
	})
}

// @Summary IsCreatorWhite
// @Tags swap
// @version 0.0.1
// @description IsCreatorWhite
// @Produce json
// @Param account query string true "account"
// @Success 200 {object} models.ResIsWhite
// @Router /stpdao/v2/swap/isCreatorWhite [get]
func SwapIsCreatorWhite(c *gin.Context) {
	accountParam := c.Query("account")

	var tbSysConfig db.TbSysConfig
	sqlSel := oo.NewSqler().Table(consts.TbSysConfig).Where("cfg_name", "cfg_swap_creator_white_list").Select()
	err := oo.SqlGet(sqlSel, &tbSysConfig)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var isCreatorWhite bool
	if tbSysConfig.CfgIsEnabled {
		sliceVal := strings.Split(tbSysConfig.CfgVal, ",")
		for _, v := range sliceVal {
			if strings.ToLower(v) == strings.ToLower(accountParam) {
				isCreatorWhite = true
				break
			}
		}
	} else {
		isCreatorWhite = true
	}
	if tbSysConfig.CfgVal == "" {
		isCreatorWhite = true
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResIsWhite{
			IsWhite: isCreatorWhite,
		},
	})
}
