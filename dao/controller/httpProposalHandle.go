package controller

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"time"
)

// @Summary query proposal list
// @Tags proposal
// @version 0.0.1
// @description query proposal list
// @Produce json
// @Param daoAddress query string true "Dao address"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResProposalsListPage
// @Router /stpdao/v2/proposal/list [get]
func httpProposalsList(c *gin.Context) {
	daoAddressParam := c.Query("daoAddress")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offestParam, _ := strconv.Atoi(offset)

	var listEntities []models.EventHistoricalModel
	sqler := oo.NewSqler().Table(consts.TbNameEventHistorical).
		Where("event_type", consts.EvCreateProposal).
		Where("address", daoAddressParam)
	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Limit(countParam).Offset(offestParam).Select()
		err = oo.SqlSelect(sqlStr, &listEntities)
	}
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResProposalsList, 0)
	for index := range listEntities {
		proposalId := utils.FixTo0xString(listEntities[index].Topic1)
		proposer := utils.FixTo0xString(listEntities[index].Topic2)
		startTime, _ := utils.Hex2Int64(utils.FixTo0xString(listEntities[index].Data[2:66]))
		endTime, _ := utils.Hex2Int64(utils.FixTo0xString(listEntities[index].Data[66:130]))

		var counts int
		sqlCancel := oo.NewSqler().Table(consts.TbNameEventHistorical).
			Where("event_type", consts.EvCancelProposal).
			Where("topic1", listEntities[index].Topic1).Count()
		err = oo.SqlGet(sqlCancel, &counts)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusOK, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		status := "Open"
		if time.Now().Unix() > endTime || counts >= 1 {
			status = "Closed"
		}

		data = append(data, models.ResProposalsList{
			DaoAddress: daoAddressParam,
			ProposalId: proposalId,
			Proposer:   proposer,
			StartTime:  startTime,
			EndTime:    endTime,
			Status:     status,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResProposalsListPage{
			List:  data,
			Total: total,
		},
	})

}
