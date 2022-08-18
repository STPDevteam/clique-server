package controller

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
)

// @Summary query votes list
// @Tags votes
// @version 0.0.1
// @description query votes list
// @Produce json
// @Param proposalId query string true "proposalId"
// @Param daoAddress query string true "Dao address"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResVotesListPage
// @Router /stpdao/v2/votes/list [get]
func httpVotesList(c *gin.Context) {
	daoAddressParam := c.Query("daoAddress")
	proposalIdParam := c.Query("proposalId")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var total uint64
	var listEntities []models.EventHistoricalModel
	proposalIdParam64 := utils.FixTo0x64String(proposalIdParam)
	sqler := oo.NewSqler().Table(consts.TbNameEventHistorical).
		Where("event_type", consts.EvVote).
		Where("address", daoAddressParam).
		Where("topic1", proposalIdParam64)
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

	var data = make([]models.ResVotesList, 0)
	for index := range listEntities {
		optionIndex := utils.FixTo0x40String(listEntities[index].Topic3)
		voter := utils.FixTo0x40String(listEntities[index].Topic2)
		amount, _ := utils.Hex2BigInt(listEntities[index].Data[:66])

		data = append(data, models.ResVotesList{
			OptionIndex: optionIndex,
			Voter:       voter,
			Amount:      amount.String(),
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResVotesListPage{
			List:  data,
			Total: total,
		},
	})

}
