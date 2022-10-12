package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
)

// @Summary overview total
// @Tags overview
// @version 0.0.1
// @description overview total
// @Produce json
// @Success 200 {object} models.ResOverview
// @Router /stpdao/v2/overview/total [get]
func httpRecordTotal(c *gin.Context) {

	var totalDao int
	sqlSel := oo.NewSqler().Table(consts.TbNameDao).Count()
	err := oo.SqlGet(sqlSel, &totalDao)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var totalApproveDao int
	sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("approve", 1).Count()
	err = oo.SqlGet(sqlSel, &totalApproveDao)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var totalVote int
	sqlSel = oo.NewSqler().Table(consts.TbNameVote).Count()
	err = oo.SqlGet(sqlSel, &totalVote)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var totalAddressVote int
	sqlSel = fmt.Sprintf(`SELECT count(DISTINCT voter) as count FROM tb_vote`)
	err = oo.SqlGet(sqlSel, &totalAddressVote)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResOverview{
			TotalDao:         totalDao,
			TotalApproveDao:  totalApproveDao,
			TotalVote:        totalVote,
			TotalAddressVote: totalAddressVote,
		},
	})

}
