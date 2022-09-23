package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"strconv"
	"time"
)

// @Summary query Dao list
// @Tags Dao
// @version 0.0.1
// @description query Dao list
// @Produce json
// @Param customizeOrderProposal query string false "customizeOrderProposal: ASC"
// @Param customizeOrderCreateTime query string false "customizeOrderCreateTime: ASC"
// @Param order query string false "account address"
// @Param account query string false "account address"
// @Param keyword query  string false "query keyword:Dao name,Dao address,Token address"
// @Param category query string false "category"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResDaoListPage
// @Router /stpdao/v2/dao/list [get]
func httpDaoList(c *gin.Context) {
	accountParam := c.Query("account")
	keywordParam := c.Query("keyword")
	categoryParam := c.Query("category")
	offset := c.Query("offset")
	count := c.Query("count")
	offsetParam, _ := strconv.Atoi(offset)
	countParam, _ := strconv.Atoi(count)
	customizeOrderProposalParam := c.Query("customizeOrderProposal")
	customizeOrderCreateTimeParam := c.Query("customizeOrderCreateTime")

	var orderStr string
	if customizeOrderProposalParam == "ASC" && customizeOrderCreateTimeParam == "" {
		orderStr = "ORDER BY weight ASC,create_time DESC"
	}
	if customizeOrderCreateTimeParam == "ASC" && customizeOrderProposalParam == "" {
		orderStr = "ORDER BY create_time ASC,weight DESC"
	}
	if customizeOrderProposalParam == "ASC" && customizeOrderCreateTimeParam == "ASC" {
		orderStr = "ORDER BY weight ASC,create_time ASC"
	}
	if customizeOrderProposalParam == "" && customizeOrderCreateTimeParam == "" {
		orderStr = "ORDER BY weight DESC,create_time DESC"
	}

	var sqlCount, sqlSel, sqlWhere, sqlOrderLimit, sqlSubquery string
	sqlCount = fmt.Sprintf(`SELECT COUNT(*) FROM %s `, consts.TbNameDao)
	sqlSel = fmt.Sprintf(`SELECT * FROM %s `, consts.TbNameDao)
	sqlOrderLimit = fmt.Sprintf(`%s Limit %d,%d `, orderStr, offsetParam, countParam)
	if keywordParam != "" {
		sqlWhere = fmt.Sprintf(`WHERE (dao_address='%s' OR token_address='%s' OR dao_name LIKE '%%%s%%') `, keywordParam, keywordParam, keywordParam)
	}
	if categoryParam != "" {
		var categoryId int
		sqlSelCategoryId := oo.NewSqler().Table(consts.TbNameCategory).Where("category_name", categoryParam).Select("id")
		err := oo.SqlGet(sqlSelCategoryId, &categoryId)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if keywordParam == "" {
			sqlSubquery = fmt.Sprintf(`WHERE id IN (SELECT dao_id FROM %s WHERE category_id = %d) `, consts.TbNameDaoCategory, categoryId)
		} else {
			sqlSubquery = fmt.Sprintf(`AND id IN (SELECT dao_id FROM %s WHERE category_id = %d) `, consts.TbNameDaoCategory, categoryId)
		}
	}

	sqlStrCount := fmt.Sprintf(`%s%s%s`, sqlCount, sqlWhere, sqlSubquery)
	sqlStrSel := fmt.Sprintf(`%s%s%s%s`, sqlSel, sqlWhere, sqlSubquery, sqlOrderLimit)

	var total uint64
	var daoListEntity []models.DaoModel
	err := oo.SqlGet(sqlStrCount, &total)
	if err == nil {
		err = oo.SqlSelect(sqlStrSel, &daoListEntity)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var listModel = make([]models.ResDaoList, 0)
	for index := range daoListEntity {

		var totalProposals uint64
		sqlTotal := oo.NewSqler().Table(consts.TbNameProposal).
			Where("chain_id", daoListEntity[index].ChainId).
			Where("dao_address", daoListEntity[index].DaoAddress).Count()
		err = oo.SqlGet(sqlTotal, &totalProposals)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var activeProposals uint64
		var now = time.Now().Unix()
		sqlActive := oo.NewSqler().Table(consts.TbNameProposal).
			Where("chain_id", daoListEntity[index].ChainId).
			Where("dao_address", daoListEntity[index].DaoAddress).
			Where("start_time", "<=", now).
			Where("end_time", ">=", now).Count()
		err = oo.SqlGet(sqlActive, &activeProposals)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var soonProposals uint64
		sqlSoon := oo.NewSqler().Table(consts.TbNameProposal).
			Where("chain_id", daoListEntity[index].ChainId).
			Where("dao_address", daoListEntity[index].DaoAddress).
			Where("start_time", ">=", now).Count()
		err = oo.SqlGet(sqlSoon, &soonProposals)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var closedProposals uint64
		sqlClosed := oo.NewSqler().Table(consts.TbNameProposal).
			Where("chain_id", daoListEntity[index].ChainId).
			Where("dao_address", daoListEntity[index].DaoAddress).
			Where("end_time", "<=", now).Count()
		err = oo.SqlGet(sqlClosed, &closedProposals)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var members uint64
		sqlMembers := oo.NewSqler().Table(consts.TbNameMember).
			Where("dao_address", daoListEntity[index].DaoAddress).
			Where("chain_id", daoListEntity[index].ChainId).
			Where("join_switch", 1).Count()
		err = oo.SqlGet(sqlMembers, &members)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var joinSwitch int
		if accountParam == "" {
			joinSwitch = 0
		} else {

			var entity []models.MemberModel
			sqlAcc := oo.NewSqler().Table(consts.TbNameMember).
				Where("account", accountParam).
				Where("dao_address", daoListEntity[index].DaoAddress).Select()
			err = oo.SqlSelect(sqlAcc, &entity)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				c.JSON(http.StatusInternalServerError, models.Response{
					Code:    500,
					Message: "Something went wrong, Please try again later.",
				})
				return
			}
			if len(entity) == 0 {
				joinSwitch = 0
			} else {
				joinSwitch = entity[0].JoinSwitch
			}

		}

		listModel = append(listModel, models.ResDaoList{
			DaoLogo:         daoListEntity[index].DaoLogo,
			DaoName:         daoListEntity[index].DaoName,
			DaoAddress:      daoListEntity[index].DaoAddress,
			ChainId:         daoListEntity[index].ChainId,
			TotalProposals:  totalProposals,
			ActiveProposals: activeProposals,
			SoonProposals:   soonProposals,
			ClosedProposals: closedProposals,
			Members:         members,
			JoinSwitch:      joinSwitch,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResDaoListPage{
			List:  listModel,
			Total: total,
		},
	})

}

// @Summary join or quit Dao
// @Tags Dao
// @version 0.0.1
// @description join or quit Dao
// @Produce json
// @Param request body models.JoinDaoWithSignParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/dao/member [post]
func httpDaoJoinOrQuit(c *gin.Context) {
	var params models.JoinDaoWithSignParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
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

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameMember).
		Where("dao_address", params.Params.DaoAddress).
		Where("account", params.Params.Account).
		Where("chain_id", params.Params.ChainId).Count()
	err = oo.SqlGet(sqlSel, &count)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var role string
	sqlSelRole := oo.NewSqler().Table(consts.TbNameAdmin).
		Where("dao_address", params.Params.DaoAddress).
		Where("chain_id", params.Params.ChainId).
		Where("account", params.Params.Account).Select("account_level")
	err = oo.SqlGet(sqlSelRole, &role)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if role != consts.LevelSuperAdmin {
		if (count == 0 && params.Params.JoinSwitch == 1) || (count == 1 && params.Params.JoinSwitch == 1) {

			sqlIns := fmt.Sprintf(`REPLACE INTO %s (dao_address,chain_id,account,join_switch) VALUES ('%s',%d,'%s',%d)`,
				consts.TbNameMember,
				params.Params.DaoAddress,
				params.Params.ChainId,
				params.Params.Account,
				params.Params.JoinSwitch,
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

		} else if count == 1 && params.Params.JoinSwitch == 0 {

			sqlUp := fmt.Sprintf(`UPDATE %s SET join_switch=%d WHERE dao_address='%s' AND account='%s' AND chain_id=%d`,
				consts.TbNameMember,
				params.Params.JoinSwitch,
				params.Params.DaoAddress,
				params.Params.Account,
				params.Params.ChainId,
			)
			err = oo.SqlExec(sqlUp)
			if err != nil {
				oo.LogW("%v", err)
				c.JSON(http.StatusOK, models.Response{
					Code:    500,
					Message: "Something went wrong, Please try again later.",
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResResult{
			Success: true,
		},
	})

}

// @Summary Dao Creator-Left
// @Tags Dao
// @version 0.0.1
// @description Dao Creator-Left
// @Produce json
// @Param account query string true "account address"
// @Success 200 {object} models.ResLeftDaoCreator
// @Router /stpdao/v2/dao/left [get]
func httpLeftDaoJoin(c *gin.Context) {
	accountParam := c.Query("account")

	var entities []models.MemberModel
	sqler := oo.NewSqler().Table(consts.TbNameMember).Where("account", accountParam).Where("join_switch", 1).Select()
	err := oo.SqlSelect(sqler, &entities)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResLeftDaoCreator, 0)
	for index := range entities {
		data = append(data, models.ResLeftDaoCreator{
			Account:    accountParam,
			DaoAddress: entities[index].DaoAddress,
			ChainId:    entities[index].ChainId,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})

}

// @Summary Dao info
// @Tags Dao
// @version 0.0.1
// @description Dao info
// @Produce json
// @Param account query string false "account address"
// @Param daoAddress query string true "dao Address"
// @Param chainId query string true "chainId"
// @Success 200 {object} models.ResDaoInfo
// @Router /stpdao/v2/dao/info [get]
func httpDaoInfo(c *gin.Context) {
	accountParam := c.Query("account")
	daoAddressParam := c.Query("daoAddress")
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)

	var members uint64
	sqlMembers := oo.NewSqler().Table(consts.TbNameMember).
		Where("dao_address", daoAddressParam).
		Where("chain_id", chainIdParam).
		Where("join_switch", 1).Count()
	err := oo.SqlGet(sqlMembers, &members)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var joinSwitch int
	if accountParam == "" {
		joinSwitch = 0
	} else {
		var entity []models.MemberModel
		sqlAcc := oo.NewSqler().Table(consts.TbNameMember).
			Where("account", accountParam).
			Where("dao_address", daoAddressParam).
			Where("chain_id", chainIdParam).Select()
		err = oo.SqlSelect(sqlAcc, &entity)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if len(entity) == 0 {
			joinSwitch = 0
		} else {
			joinSwitch = entity[0].JoinSwitch
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResDaoInfo{
			Members:    members,
			JoinSwitch: joinSwitch,
		},
	})
}

// @Summary Dao admins
// @Tags Dao
// @version 0.0.1
// @description Dao admins
// @Produce json
// @Param daoAddress query string true "dao Address"
// @Param chainId query string true "chainId"
// @Success 200 {object} models.ResAdminsList
// @Router /stpdao/v2/dao/admins [get]
func httpDaoAdmins(c *gin.Context) {
	daoAddressParam := c.Query("daoAddress")
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)

	var adminEntities []models.AdminModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAdmin).
		Where("dao_address", daoAddressParam).
		Where("chain_id", chainIdParam).
		Where("account_level='superAdmin' OR account_level='admin'").Distinct().Select("account")
	err := oo.SqlSelect(sqlSel, &adminEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResAdminsList, 0)
	for index := range adminEntities {
		data = append(data, models.ResAdminsList{
			Account: adminEntities[index].Account,
			//AccountLevel: adminEntities[index].AccountLevel,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})

}
