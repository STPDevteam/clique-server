package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
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
		c.JSON(http.StatusBadRequest, models.Response{
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
	sqlCount := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Account).Count()
	err = oo.SqlGet(sqlCount, &counts)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
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
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
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
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var myTokenEntities []models.HolderDataModel
	sqlSel = oo.NewSqler().Table(consts.TbNameHolderData).Where("holder_address", params.Account).Select()
	err = oo.SqlSelect(sqlSel, &myTokenEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
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
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataSuper = make([]models.ResDao, 0)
	for index := range superDaoEntities {
		dataSuper = append(dataSuper, models.ResDao{
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
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataAdmin = make([]models.ResDao, 0)
	for index := range adminDaoEntities {
		dataAdmin = append(dataAdmin, models.ResDao{
			DaoAddress:   adminDaoEntities[index].DaoAddress,
			ChainId:      adminDaoEntities[index].ChainId,
			AccountLevel: adminDaoEntities[index].AccountLevel,
		})
	}

	var activityEntities []models.EventHistoricalModel
	account0x64 := utils.FixTo0x64String(params.Account)
	sqlActivity := oo.NewSqler().Table(consts.TbNameEventHistorical).
		Where("event_type='CreateProposal' OR event_type='Vote'").
		Where("topic2", account0x64).
		Order("time_stamp DESC").Limit(5).Offset(0).Select()
	err = oo.SqlSelect(sqlActivity, &activityEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataActivity = make([]models.ResActivity, 0)
	for index := range activityEntities {
		dataIndex := activityEntities[index]
		proposalId := utils.Hex2Dec(dataIndex.Topic1)
		if dataIndex.EventType == consts.EvCreateProposal {
			dataActivity = append(dataActivity, models.ResActivity{
				EventType:  dataIndex.EventType,
				ChainId:    dataIndex.ChainId,
				DaoAddress: dataIndex.Address,
				ProposalId: proposalId,
			})
		}
		if dataIndex.EventType == consts.EvVote {
			optionIndex := utils.Hex2Dec(dataIndex.Topic3)
			amount, _ := utils.Hex2BigInt(dataIndex.Data[:66])
			dataActivity = append(dataActivity, models.ResActivity{
				EventType:   dataIndex.EventType,
				ChainId:     dataIndex.ChainId,
				DaoAddress:  dataIndex.Address,
				ProposalId:  proposalId,
				OptionIndex: optionIndex,
				Amount:      amount.String(),
			})
		}
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
			SuperDao:     dataSuper,
			AdminDao:     dataAdmin,
			Activity:     dataActivity,
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
		oo.LogW("SQL err: %v", err)
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
