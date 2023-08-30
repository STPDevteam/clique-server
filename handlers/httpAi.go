package handlers

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"stp_dao_v2/errs"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
)

// @Summary ai chat
// @Tags ai
// @version 0.0.1
// @description request header: Authorization=Bearer ${JWT Token}
// @Produce json
// @Param request body models.ReqAiChai true "request"
// @Success 200 {object} models.ResAiChat
// @Router /stpdao/v2/ai [post]
func Ai(c *gin.Context) {
	//var ok bool
	//var user *db.TbAccountModel
	//user, ok = parseJWTCache(c)
	//if !ok {
	//	return
	//}

	var params models.ReqAiChai
	if handleErrorIfExists(c, c.ShouldBindJSON(&params), errs.ErrParam) {
		return
	}

	var message string
	for _, v := range params.Content {
		if message == "" {
			message = fmt.Sprintf(`{"role": "user", "content": "%s"}`, v)
		} else {
			message = fmt.Sprintf(`%s,{"role": "user", "content": "%s"}`, message, v)
		}
	}

	url := "https://api.openai.com/v1/chat/completions"
	chat, err := utils.AiChat("1", message, url, viper.GetString("app.openai_bearer_key"))
	if handleErrorIfExists(c, err, errs.ErrServer) {
		oo.LogW("AiChat err:%v", err)
		return
	}

	jsonData(c, chat)
}
