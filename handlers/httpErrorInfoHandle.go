package handlers

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
)

// @Summary error info
// @Tags Error
// @version 0.0.1
// @description error info
// @Produce json
// @Param request body models.ErrorInfoParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/error/info [post]
func HttpErrorInfo(c *gin.Context) {
	var params models.ErrorInfoParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["title"] = params.Title[:int(math.Min(float64(len(params.Title)), 512))]
	v["content"] = params.Content
	v["func"] = params.Func[:int(math.Min(float64(len(params.Func)), 128))]
	v["params"] = params.Params[:int(math.Min(float64(len(params.Params)), 512))]
	m = append(m, v)

	sqlIns := oo.NewSqler().Table(consts.TbNameErrorInfo).Insert(m)
	err = oo.SqlExec(sqlIns)
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
		Data: models.ResResult{
			Success: true,
		},
	})
}
