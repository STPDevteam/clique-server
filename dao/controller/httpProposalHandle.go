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
	offsetParam, _ := strconv.Atoi(offset)

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
		sqlStr = sqlCopy.Limit(countParam).Offset(offsetParam).Select()
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
		proposalId := utils.FixTo0x40String(listEntities[index].Topic1)
		proposer := utils.FixTo0x40String(listEntities[index].Topic2)
		startTime, _ := utils.Hex2Int64(utils.FixTo0x40String(listEntities[index].Data[2:66]))
		endTime, _ := utils.Hex2Int64(utils.FixTo0x40String(listEntities[index].Data[66:130]))

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

// @Summary save proposal info
// @Tags proposal
// @version 0.0.1
// @description save proposal info
// @Produce json
// @Param request body models.ProposalInfoParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/proposal/save [post]
func httpSaveProposal(c *gin.Context) {
	var params models.ProposalInfoParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	uuid := utils.GenerateUuid()
	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["uuid"] = uuid
	v["content"] = params.Content
	m = append(m, v)
	sqlIns := oo.NewSqler().Table(consts.TbNameProposalInfo).Insert(m)
	err = oo.SqlExec(sqlIns)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResProposalUuid{
			Uuid: uuid,
		},
	})
}

// @Summary query proposal info
// @Tags proposal
// @version 0.0.1
// @description query proposal info
// @Produce json
// @Param uuid query string true "uuid"
// @Success 200 {object} models.ResProposalContent
// @Router /stpdao/v2/proposal/query [get]
func httpQueryProposal(c *gin.Context) {
	uuidParams := c.Query("uuid")

	var content string
	sqlSel := oo.NewSqler().Table(consts.TbNameProposalInfo).Where("uuid", uuidParams).Select("content")
	err := oo.SqlGet(sqlSel, &content)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResProposalContent{
			Content: content,
		},
	})

}
