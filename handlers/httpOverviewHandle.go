package handlers

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
func HttpRecordTotal(c *gin.Context) {

	var totalDao int
	sqlSel := oo.NewSqler().Table(consts.TbNameDao).Where("deprecated", 0).Count()
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
	sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("deprecated", 0).Where("approve", 1).Count()
	err = oo.SqlGet(sqlSel, &totalApproveDao)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var totalAccount uint64
	sqlSel = fmt.Sprintf(`SELECT count(DISTINCT account) as count FROM %s where join_switch=1`, consts.TbNameMember)
	err = oo.SqlGet(sqlSel, &totalAccount)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	//totalAccount, ok := svc.mCache.Get(consts.CacheTokenHolders)
	//if !ok {
	//	c.JSON(http.StatusInternalServerError, models.Response{
	//		Code:    500,
	//		Message: "Something went wrong, Please try again later.",
	//	})
	//	return
	//}
	//val, ok2 := totalAccount.(uint64)
	//if !ok2 {
	//	c.JSON(http.StatusInternalServerError, models.Response{
	//		Code:    500,
	//		Message: "Something went wrong, Please try again later.",
	//	})
	//	return
	//}

	var totalProposal int
	sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("deprecated", 0).Count()
	err = oo.SqlGet(sqlSel, &totalProposal)
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
			TotalDao:        totalDao,
			TotalApproveDao: totalApproveDao,
			TotalAccount:    totalAccount,
			TotalProposal:   totalProposal,
		},
	})

}
