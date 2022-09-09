package controller

import (
	"encoding/json"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"strconv"
	"time"
)

// @Summary query activity list
// @Tags activity
// @version 0.0.1
// @description query activity list
// @Produce json
// @Param chainId query int false "chainId"
// @Param daoAddress query string false "Dao address"
// @Param status query int false "status:Soon:1,Open:2,Closed:3"
// @Param types query string false "types:Airdrop,PublicSale"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResActivityPage
// @Router /stpdao/v2/activity/list [get]
func httpActivity(c *gin.Context) {
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")
	status := c.Query("status")
	statusParam, _ := strconv.Atoi(status)
	typesParam := c.Query("types")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var listEntities []models.ActivityModel
	sqler := oo.NewSqler().Table(consts.TbNameActivity)
	if chainIdParam != 0 && daoAddressParam != "" {
		sqler = sqler.Where("chain_id", chainIdParam).Where("dao_address", daoAddressParam)
	}
	if typesParam != "" {
		sqler = sqler.Where("types", typesParam)
	}
	var now = time.Now().Unix()
	if statusParam == 1 {
		sqler = sqler.Where("start_time", ">=", now)
	}
	if statusParam == 2 {
		sqler = sqler.Where("end_time", ">=", now).Where("start_time", "<=", now)
	}
	if statusParam == 3 {
		sqler = sqler.Where("end_time", "<=", now)
	}

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("start_time DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &listEntities)
	}
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResActivityList, 0)
	for index := range listEntities {
		dataIndex := listEntities[index]

		var entity []models.AirdropAddressModel
		sqlSel := oo.NewSqler().Table(consts.TbNameAirdropAddress).Where("id", dataIndex.ActivityId).Select()
		err = oo.SqlSelect(sqlSel, &entity)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var content models.AirdropAddressArray
		err = json.Unmarshal([]byte(entity[0].Content), &content)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Json Unmarshal Failed.",
			})
			return
		}

		var claimedCount int
		sqlSel = oo.NewSqler().Table(consts.TbNameClaimed).Where("airdrop_id", dataIndex.ActivityId).Count()
		err = oo.SqlGet(sqlSel, &claimedCount)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var claimedPercentage float64
		if len(content.Address) == 0 {
			claimedPercentage = 0
		} else {
			claimedPercentage = float64(claimedCount) / float64(len(content.Address))
		}

		data = append(data, models.ResActivityList{
			Title:             entity[0].Title,
			Types:             dataIndex.Types,
			ChainId:           dataIndex.ChainId,
			DaoAddress:        dataIndex.DaoAddress,
			Creator:           dataIndex.Creator,
			ActivityId:        dataIndex.ActivityId,
			TokenAddress:      dataIndex.TokenAddress,
			Amount:            dataIndex.Amount,
			StartTime:         dataIndex.StartTime,
			EndTime:           dataIndex.EndTime,
			PublishTime:       dataIndex.PublishTime,
			Price:             dataIndex.Price,
			AirdropNumber:     len(content.Address),
			ClaimedPercentage: claimedPercentage,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResActivityPage{
			List:  data,
			Total: total,
		},
	})

}
