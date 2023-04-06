package middlewares

import (
	"net/http"
	"stp_dao_v2/errs"

	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		err := ctx.Errors.Last()
		if err != nil {
			if cErr, ok := err.Err.(*errs.CustomError); ok {
				ctx.JSON(http.StatusOK, cErr)
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code": cErr.Code,
					"msg":  cErr.Msg,
					"data": err.Err,
				})
			}
		}
	}
}
