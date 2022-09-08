package controller

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"strconv"
)

// @Summary query votes list
// @Tags votes
// @version 0.0.1
// @description query votes list
// @Produce json
// @Param chainId query int true "chainId"
// @Param proposalId query string true "proposalId"
// @Param daoAddress query string true "Dao address"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResVotesListPage
// @Router /stpdao/v2/votes/list [get]
func httpVotesList(c *gin.Context) {
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	daoAddressParam := c.Query("daoAddress")
	proposalIdParam := c.Query("proposalId")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var total int
	var listEntities []models.VoteModel
	sqler := oo.NewSqler().Table(consts.TbNameVote).
		Where("chain_id", chainIdParam).
		Where("dao_address", daoAddressParam).
		Where("proposal_id", proposalIdParam)
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("create_time DESC").Limit(countParam).Offset(offsetParam).Select()
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

	var data = make([]models.ResVotesList, 0)
	for index := range listEntities {
		dataIndex := listEntities[index]
		data = append(data, models.ResVotesList{
			ProposalId:  dataIndex.ProposalId,
			Voter:       dataIndex.Voter,
			OptionIndex: dataIndex.OptionIndex,
			Amount:      dataIndex.Amount,
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
