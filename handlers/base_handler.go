package handlers

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/errs"
)

// req model
type ReqPagination struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

func jsonData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data, "msg": "success"})
}

func jsonPagination(c *gin.Context, data interface{}, total int64, pagination ReqPagination) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code":   http.StatusOK,
			"msg":    "success",
			"data":   data,
			"total":  total,
			"offset": pagination.Offset,
			"limit":  pagination.Limit,
		},
	)
}

func jsonSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "success"})
}

func handleError(c *gin.Context, cErr *errs.CustomError) {
	oo.LogW("%s: custom error: %v", c.FullPath(), cErr)
	c.Abort()
	c.Error(cErr)
}

func handleErrorIfExists(c *gin.Context, err error, cErr *errs.CustomError) bool {
	if err != nil {
		oo.LogW("%s: error : %v, custom error: %v", c.FullPath(), err, cErr)
		handleError(c, cErr)
		return true
	}
	return false
}

func handleErrorIfExistsExceptNoRows(c *gin.Context, err error, cErr *errs.CustomError) bool {
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("%s: error : %v, custom error: %v", c.FullPath(), err, cErr)
		handleError(c, cErr)
		return true
	}
	return false
}

func HandlerPagination(c *gin.Context) {
	var err error
	var pagination ReqPagination
	err = c.ShouldBindQuery(&pagination)
	if handleErrorIfExists(c, err, errs.ErrParam) {
		return
	}
	if pagination.Limit > 100 {
		handleError(c, errs.ErrParam)
		return
	}
	if pagination.Limit == 0 {
		handleError(c, errs.ErrParam)
		return
	}
}
