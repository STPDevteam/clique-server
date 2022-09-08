package controller

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"strconv"
	"time"
)

// @Summary query votes list
// @Tags votes
// @version 0.0.1
// @description query votes list
// @Produce json
// @Param chainId query int false "chainId"
// @Param daoAddress query string false "Dao address"
// @Param status query int false "status:Soon:1,Open:2,Closed:3"
// @Param title query int false "title"
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
	titleParam := c.Query("title")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var listEntities []models.ActivityModel
	sqler := oo.NewSqler().Table(consts.TbNameActivity)
	if chainIdParam != 0 && daoAddressParam != "" {
		sqler = sqler.Where("chain_id", chainIdParam).Where("dao_address", daoAddressParam)
	}
	if titleParam != "" {
		sqler = sqler.Where("title", titleParam)
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
		data = append(data, models.ResActivityList{
			Title:        dataIndex.Title,
			ChainId:      dataIndex.ChainId,
			DaoAddress:   dataIndex.DaoAddress,
			Creator:      dataIndex.Creator,
			ActivityId:   dataIndex.ActivityId,
			TokenAddress: dataIndex.TokenAddress,
			Amount:       dataIndex.Amount,
			StartTime:    dataIndex.StartTime,
			EndTime:      dataIndex.EndTime,
			Price:        dataIndex.Price,
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
