package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"math/big"
	"stp_dao_v2/consts"
	"stp_dao_v2/errs"
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
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	var lastId int64
	var signature string
	tx, errTx := oo.NewSqlxTx()
	if handleErrorIfExists(c, errTx, errs.ErrServer) {
		return
	}
	defer func() {
		oo.CloseSqlxTx(tx, &errTx)
		if errTx == nil {
			jsonData(c, models.ResCreateSale{
				SaleId:    lastId,
				Signature: signature,
			})
		}
	}()

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["chain_id"] = params.ChainId
	v["creator"] = params.Creator
	v["sale_way"] = params.SaleWay
	v["sale_token"] = params.SaleToken
	v["sale_amount"] = params.SaleAmount
	v["sale_price"] = params.SalePrice
	//v["original_price"] =
	v["receive_token"] = params.ReceiveToken
	v["limit_min"] = params.LimitMin
	v["limit_max"] = params.LimitMax
	v["start_time"] = params.StartTime
	v["end_time"] = params.EndTime
	v["white_list"] = params.WhiteList
	v["about"] = params.About
	m = append(m, v)
	sqlIns := oo.NewSqler().Table(consts.TbNameSwap).Insert(m)
	res, errTx := oo.SqlxTxExec(tx, sqlIns)
	if handleErrorIfExists(c, errTx, errs.ErrServer) {
		oo.LogW("SQL err: %v", errTx)
		return
	}

	lastId, errTx = res.LastInsertId()
	if handleErrorIfExists(c, errTx, errs.ErrServer) {
		oo.LogW("LastInsertId err: %v", errTx)
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
	if handleErrorIfExists(c, errTx, errs.ErrServer) {
		oo.LogW("SignMessage err: %v", errTx)
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
// @Success 200 {object} models.
// @Router /stpdao/v2/swap/list [get]
func swapList(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	limitParam, _ := strconv.Atoi(limit)
	offsetParam, _ := strconv.Atoi(offset)

	var swapArr []models.TbSwap
	sqler := oo.NewSqler().Table(consts.TbNameSwap).Where("status", "!=", consts.StatusPending)

	var total int64
	sqlCopy := *sqler
	sqlSel := sqlCopy.Count()
	err := oo.SqlGet(sqlSel, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlSel = sqlCopy.Order("status DESC,create_time DESC").Limit(limitParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlSel, &swapArr)
	}
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("SQL err: %v, str: %s", err, sqlSel)
		return
	}

	var data = make([]models.ResSwapList, 0)
	for index := range swapArr {
		ls := swapArr[index]

		data = append(data, models.ResSwapList{
			SaleId: ls.Id,
		})
	}

	page := ReqPagination{
		Offset: offsetParam,
		Limit:  limitParam,
	}

	jsonPagination(c, data, total, page)
}
