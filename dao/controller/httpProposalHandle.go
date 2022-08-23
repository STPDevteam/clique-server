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
// @Param status query int false "status:Soon:1,Open:2,Closed:3"
// @Param daoAddress query string true "Dao address"
// @Param chainId query int true "chainId"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResProposalsListPage
// @Router /stpdao/v2/proposal/list [get]
func httpProposalsList(c *gin.Context) {
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)
	status := c.Query("status")
	statusParam, _ := strconv.Atoi(status)

	var listEntities []models.ProposalModel
	sqler := oo.NewSqler().Table(consts.TbNameProposal).
		Where("chain_id", chainIdParam).
		Where("dao_address", daoAddressParam)
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
		sqlStr = sqlCopy.Order("proposal_id DESC").Limit(countParam).Offset(offsetParam).Select()
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
		data = append(data, models.ResProposalsList{
			ChainId:    chainIdParam,
			DaoAddress: daoAddressParam,
			ProposalId: listEntities[index].ProposalId,
			Proposer:   listEntities[index].Proposer,
			StartTime:  listEntities[index].StartTime,
			EndTime:    listEntities[index].EndTime,
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

// @Summary query proposal snapshot
// @Tags proposal
// @version 0.0.1
// @description query proposal snapshot
// @Produce json
// @Param chainId query int true "chainId"
// @Param daoAddress query string true "daoAddress"
// @Param proposalId query string true "proposalId"
// @Success 200 {object} models.ResProposalContent
// @Router /stpdao/v2/proposal/snapshot [get]
func httpQuerySnapshot(c *gin.Context) {
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	proposalId := c.Query("proposalId")
	proposalIdParam, _ := strconv.Atoi(proposalId)
	daoAddressParam := c.Query("daoAddress")

	var blockNumber string
	proposalId0x64 := utils.FixTo0x64String(strconv.FormatInt(int64(proposalIdParam), 16))
	sqlSel := oo.NewSqler().Table(consts.TbNameEventHistorical).
		Where("event_type", consts.EvCreateProposal).
		Where("address", daoAddressParam).
		Where("chain_id", chainIdParam).
		Where("topic1", proposalId0x64).Select("block_number")
	err := oo.SqlGet(sqlSel, &blockNumber)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	snapshot, _ := utils.Hex2Int64(blockNumber)
	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResSnapshot{
			ChainId:    chainIdParam,
			DaoAddress: daoAddressParam,
			ProposalId: proposalIdParam,
			Snapshot:   snapshot,
		},
	})
}
