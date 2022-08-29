package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
)

// @Summary account info
// @Tags account
// @version 0.0.1
// @description account info
// @Produce json
// @Param request body models.SignData true "request"
// @Success 200 {object} models.ResQueryAccount
// @Router /stpdao/v2/account/query [post]
func httpQueryAccount(c *gin.Context) {
	var params models.SignData
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	if !checkLogin(&params) {
		oo.LogD("SignData err not auth")
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Data:    models.ResResult{Success: false},
			Message: "SignData err not auth",
		})
		return
	}

	var counts int
	sqler := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Account).Count()
	err = oo.SqlGet(sqler, &counts)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if counts == 0 {
		sqlIns := fmt.Sprintf(`INSERT INTO %s (account) VALUES ('%s')`,
			consts.TbNameAccount,
			params.Account,
		)
		err = oo.SqlExec(sqlIns)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusOK, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
	}

	var entity models.AccountModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Account).Select()
	err = oo.SqlGet(sqlSel, &entity)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var myTokenEntities []models.HolderDataModel
	sqlSel = oo.NewSqler().Table(consts.TbNameHolderData).Where("holder_address", params.Account).Select()
	err = oo.SqlSelect(sqlSel, &myTokenEntities)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataMyTokens = make([]models.ResMyTokens, 0)
	for index := range myTokenEntities {
		dataMyTokens = append(dataMyTokens, models.ResMyTokens{
			TokenAddress: myTokenEntities[index].TokenAddress,
			ChainId:      myTokenEntities[index].ChainId,
			Balance:      myTokenEntities[index].Balance,
		})
	}

	var superDaoEntities []models.AdminModel
	sqlSel = oo.NewSqler().Table(consts.TbNameAdmin).
		Where("account", params.Account).
		Where("account_level", consts.LevelSuperAdmin).Select()
	err = oo.SqlSelect(sqlSel, &superDaoEntities)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var superData = make([]models.ResDao, 0)
	for index := range superDaoEntities {
		superData = append(superData, models.ResDao{
			DaoAddress:   superDaoEntities[index].DaoAddress,
			ChainId:      superDaoEntities[index].ChainId,
			AccountLevel: superDaoEntities[index].AccountLevel,
		})
	}

	var adminDaoEntities []models.AdminModel
	sqlSel = oo.NewSqler().Table(consts.TbNameAdmin).
		Where("account", params.Account).
		Where("account_level", consts.LevelAdmin).Select()
	err = oo.SqlSelect(sqlSel, &adminDaoEntities)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var adminData = make([]models.ResDao, 0)
	for index := range adminDaoEntities {
		adminData = append(adminData, models.ResDao{
			DaoAddress:   adminDaoEntities[index].DaoAddress,
			ChainId:      adminDaoEntities[index].ChainId,
			AccountLevel: adminDaoEntities[index].AccountLevel,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResQueryAccount{
			Account:      entity.Account,
			AccountLogo:  entity.AccountLogo.String,
			Nickname:     entity.Nickname.String,
			Introduction: entity.Introduction.String,
			Twitter:      entity.Twitter.String,
			Github:       entity.Github.String,
			MyTokens:     dataMyTokens,
			SuperDao:     superData,
			AdminDao:     adminData,
		},
	})
}

// @Summary update account info
// @Tags account
// @version 0.0.1
// @description update account info
// @Produce json
// @Param request body models.UpdateAccountWithSignParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/account/update [post]
func httpUpdateAccount(c *gin.Context) {
	var params models.UpdateAccountWithSignParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	if !checkLogin(&params.Sign) {
		oo.LogD("SignData err not auth")
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Data:    models.ResResult{Success: false},
			Message: "SignData err not auth",
		})
		return
	}

	//maybe error: account_logo=""
	sqler := fmt.Sprintf(`UPDATE %s SET account_logo='%s',nickname='%s',introduction='%s',twitter='%s',github='%s' WHERE account='%s'`,
		consts.TbNameAccount,
		params.Param.AccountLogo,
		params.Param.Nickname,
		params.Param.Introduction,
		params.Param.Twitter,
		params.Param.Github,
		params.Sign.Account,
	)
	err = oo.SqlExec(sqler)
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
	})
}
