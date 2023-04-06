package handlers

import (
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/errs"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
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

func checkLogin(sign *models.SignData) (ret bool) {
	ret, errSign := utils.CheckPersonalSign(consts.SignMessagePrefix, sign.Account, sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v, signature: %v", errSign, sign.Signature)
		return false
	}
	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}

func checkAirdropAdminAndTimestamp(sign *models.AirdropAdminSignData) (ret bool) {
	ret, errSign := utils.CheckPersonalSign(sign.Message, sign.Account, sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v", errSign)
		return false
	}

	var data models.AdminMessage
	err := json.Unmarshal([]byte(sign.Message), &data)
	if err != nil {
		oo.LogD("signMessage Unmarshal err %v", errSign)
		return false
	}

	if !utils.CheckAdminSignMessageTimestamp(data.Expired) {
		oo.LogD("signMessage CheckAdminSignMessageTimestamp err %v", errSign)
		return false
	}

	//if data.Type == "airdrop2" {
	//	root, err := merkelTreeRoot(sign.Array)
	//	log.Println(fmt.Sprintf("rootStr: %s", root))
	//	if err != nil || root != data.Root {
	//		oo.LogD("signMessage err rootMe:%v.root:%v", root, data.Root)
	//		return false
	//	}
	//}

	var count int
	var sqlSql string
	if data.Type == "airdrop1" {
		sqlSql = oo.NewSqler().Table(consts.TbNameAdmin).Where("chain_id", sign.ChainId).Where("dao_address", sign.DaoAddress).
			Where("account", sign.Account).Where("account_level='superAdmin' OR account_level='admin'").Count()
	} else if data.Type == "airdrop2" || data.Type == "airdropDownload" {
		sqlSql = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", sign.AirdropId).Where("creator", sign.Account).Count()
	}
	err = oo.SqlGet(sqlSql, &count)
	if err != nil || count == 0 {
		return false
	}

	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}

func checkAccountJoinOrQuit(data *models.JoinDaoWithSignParam) (ret bool) {
	message := fmt.Sprintf(`%d,%s,%s,%d`, data.Params.ChainId, data.Params.DaoAddress, data.Params.JoinSwitch, data.Params.Timestamp)
	ret, errSign := utils.CheckPersonalSign(message, data.Sign.Account, data.Sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v", errSign)
		return false
	}

	if !utils.CheckAdminSignMessageTimestamp(data.Params.Timestamp) {
		oo.LogD("signMessage CheckAdminSignMessageTimestamp err %v", errSign)
		return false
	}

	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}
