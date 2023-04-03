package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
	"time"
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

	var counts int
	sqlCount := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Account).Count()
	err = oo.SqlGet(sqlCount, &counts)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if counts == 0 {
		var m = make([]map[string]interface{}, 0)
		var v = make(map[string]interface{})
		v["account"] = params.Account
		m = append(m, v)
		sqlIns := oo.NewSqler().Table(consts.TbNameAccount).Insert(m)
		err = oo.SqlExec(sqlIns)
		if err != nil {
			oo.LogW("SQL err: %v", err)
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
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	//var myTokenEntities []models.HolderDataModel
	//sqlSel = oo.NewSqler().Table(consts.TbNameHolderData).Where("holder_address", params.Account).Select()
	//err = oo.SqlSelect(sqlSel, &myTokenEntities)
	//if err != nil {
	//	oo.LogW("SQL err: %v", err)
	//	c.JSON(http.StatusInternalServerError, models.Response{
	//		Code:    500,
	//		Message: "Something went wrong, Please try again later.",
	//	})
	//	return
	//}
	//var dataMyTokens = make([]models.ResMyTokens, 0)
	//for index := range myTokenEntities {
	//	dataMyTokens = append(dataMyTokens, models.ResMyTokens{
	//		TokenAddress: myTokenEntities[index].TokenAddress,
	//		ChainId:      myTokenEntities[index].ChainId,
	//		Balance:      myTokenEntities[index].Balance,
	//	})
	//}

	var adminDaoEntities []models.AdminModel
	sqlSel = oo.NewSqler().Table(consts.TbNameAdmin).
		Where("account", params.Account).
		Where("account_level='superAdmin' OR account_level='admin'").Select()
	err = oo.SqlSelect(sqlSel, &adminDaoEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataAdmin = make([]models.ResDao, 0)
	for index := range adminDaoEntities {
		var daoEntity []models.DaoModel
		sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", adminDaoEntities[index].ChainId).Where("dao_address", adminDaoEntities[index].DaoAddress).Select()
		err = oo.SqlSelect(sqlSel, &daoEntity)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusOK, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if len(daoEntity) > 0 {
			if !daoEntity[0].Deprecated {
				dataAdmin = append(dataAdmin, models.ResDao{
					DaoAddress:   adminDaoEntities[index].DaoAddress,
					ChainId:      adminDaoEntities[index].ChainId,
					AccountLevel: adminDaoEntities[index].AccountLevel,
					DaoName:      daoEntity[0].DaoName,
					DaoLogo:      daoEntity[0].DaoLogo,
				})
			}
		}
	}

	var memberEntities []models.MemberModel
	sqlSel = oo.NewSqler().Table(consts.TbNameMember).Where("account", params.Account).Where("join_switch", 1).Select()
	err = oo.SqlSelect(sqlSel, &memberEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var dataMember = make([]models.ResDao, 0)
	for index := range memberEntities {
		var success = false
		for indexAdmin := range adminDaoEntities {
			if adminDaoEntities[indexAdmin].ChainId == memberEntities[index].ChainId && adminDaoEntities[indexAdmin].DaoAddress == memberEntities[index].DaoAddress {
				success = true
				break
			}
		}
		if !success {
			var daoEntity []models.DaoModel
			sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", memberEntities[index].ChainId).Where("dao_address", memberEntities[index].DaoAddress).Select()
			err = oo.SqlSelect(sqlSel, &daoEntity)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				c.JSON(http.StatusOK, models.Response{
					Code:    500,
					Message: "Something went wrong, Please try again later.",
				})
				return
			}
			if len(daoEntity) > 0 {
				if !daoEntity[0].Deprecated {
					dataMember = append(dataMember, models.ResDao{
						DaoAddress:   memberEntities[index].DaoAddress,
						ChainId:      memberEntities[index].ChainId,
						AccountLevel: consts.LevelMember,
						DaoName:      daoEntity[0].DaoName,
						DaoLogo:      daoEntity[0].DaoLogo,
					})
				}
			}
		}
	}

	var followersCount int
	sqlSel = oo.NewSqler().Table(consts.TbNameAccountFollow).Where("followed", entity.Account).Where("status", 1).Count()
	err = oo.SqlGet(sqlSel, &followersCount)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var followingCount int
	sqlSel = oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", entity.Account).Where("status", 1).Count()
	err = oo.SqlGet(sqlSel, &followingCount)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusOK, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var email string
	if checkLogin(&params) {
		email = entity.Email.String
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResQueryAccount{
			Account:      entity.Account,
			AccountLogo:  entity.AccountLogo.String,
			Followers:    followersCount,
			Following:    followingCount,
			Nickname:     entity.Nickname.String,
			Introduction: entity.Introduction.String,
			Twitter:      entity.Twitter.String,
			Github:       entity.Github.String,
			Discord:      entity.Discord.String,
			Email:        email,
			Country:      entity.Country.String,
			Youtube:      entity.Youtube.String,
			Opensea:      entity.Opensea.String,
			//MyTokens:     dataMyTokens,
			AdminDao:             dataAdmin,
			MemberDao:            dataMember,
			AllDaosICreateOrJoin: entity.PushSwitch&(1<<0) > 0,
			NewDao:               entity.PushSwitch&(1<<1) > 0,
			AllDaoAirdrop:        entity.PushSwitch&(1<<2) > 0,
			AllDaoProposal:       entity.PushSwitch&(1<<3) > 0,
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

	var v = make(map[string]interface{})
	v["account_logo"] = params.Param.AccountLogo[:int(math.Min(float64(len(params.Param.AccountLogo)), 128))]
	v["nickname"] = params.Param.Nickname[:int(math.Min(float64(len(params.Param.Nickname)), 128))]
	v["introduction"] = params.Param.Introduction[:int(math.Min(float64(len(params.Param.Introduction)), 200))]
	v["twitter"] = params.Param.Twitter[:int(math.Min(float64(len(params.Param.Twitter)), 128))]
	v["github"] = params.Param.Github[:int(math.Min(float64(len(params.Param.Github)), 128))]
	v["discord"] = params.Param.Discord[:int(math.Min(float64(len(params.Param.Discord)), 128))]
	v["email"] = params.Param.Email[:int(math.Min(float64(len(params.Param.Email)), 128))]
	v["country"] = params.Param.Country[:int(math.Min(float64(len(params.Param.Country)), 128))]
	v["youtube"] = params.Param.Youtube[:int(math.Min(float64(len(params.Param.Youtube)), 128))]
	v["opensea"] = params.Param.Opensea[:int(math.Min(float64(len(params.Param.Opensea)), 128))]
	sqler := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Sign.Account).Update(v)
	err = oo.SqlExec(sqler)
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
	})
}

// @Summary account record list
// @Tags account
// @version 0.0.1
// @description account record list
// @Produce json
// @Param account query string true "account address"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResAccountRecordPage
// @Router /stpdao/v2/account/record [get]
func httpQueryRecordList(c *gin.Context) {
	accountParam := c.Query("account")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var entities []models.AccountRecordModel
	sqler := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("creator", accountParam)

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("time DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &entities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResAccountRecord, 0)
	for index := range entities {
		data = append(data, models.ResAccountRecord{
			Creator:    entities[index].Creator,
			Types:      entities[index].Types,
			ChainId:    entities[index].ChainId,
			Address:    entities[index].Address,
			ActivityId: entities[index].ActivityId,
			Avatar:     entities[index].Avatar,
			DaoName:    entities[index].DaoName,
			Titles:     entities[index].Titles,
			Time:       entities[index].Time,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAccountRecordPage{
			List:  data,
			Total: total,
		},
	})

}

// @Summary account sign record list
// @Tags account
// @version 0.0.1
// @description account sign record list
// @Produce json
// @Param chainId query  int false "chainId"
// @Param daoAddress query  string false "daoAddress"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResAccountSignPage
// @Router /stpdao/v2/account/sign/list [get]
func httpQueryAccountSignList(c *gin.Context) {
	count := c.Query("count")
	offset := c.Query("offset")
	chainId := c.Query("chainId")
	daoAddressParam := c.Query("daoAddress")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)
	chainIdParam, _ := strconv.Atoi(chainId)

	var entities []models.AccountSignModel
	sqler := oo.NewSqler().Table(consts.TbNameAccountSign).Where("timestamp", "<", time.Now().Unix()-60)

	if chainIdParam != 0 && daoAddressParam != "" {
		sqler = sqler.Where("chain_id", chainIdParam).Where("dao_address", daoAddressParam)
	}

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("timestamp DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &entities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResAccountSign, 0)
	for index := range entities {
		var entity models.AccountModel
		sqlSel := oo.NewSqler().Table(consts.TbNameAccount).Where("account", entities[index].Account).Select()
		err = oo.SqlGet(sqlSel, &entity)
		if err != nil && err != oo.ErrNoRows {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		data = append(data, models.ResAccountSign{
			ChainId:     entities[index].ChainId,
			DaoAddress:  entities[index].DaoAddress,
			Account:     entities[index].Account,
			Operate:     entities[index].Operate,
			Signature:   entities[index].Signature,
			Message:     entities[index].Message,
			Timestamp:   entities[index].Timestamp,
			AccountLogo: entity.AccountLogo.String,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAccountSignPage{
			List:  data,
			Total: total,
		},
	})

}

// @Summary account NFTs
// @Tags account
// @version 0.0.1
// @description account NFTs
// @Produce json
// @Param chainId query  int true "chainId 1 or 56"
// @Param account query  string true "account"
// @Param size query  int true "size,page default 20"
// @Param index query  int true "index,page default 1"
// @Success 200 {object} models.JsonRPCAccountNFT
// @Router /stpdao/v2/account/nfts [get]
func (svc *Service) httpQueryAccountNFTsList(c *gin.Context) {
	pageIndex := c.Query("index")
	PageSize := c.Query("size")
	chainId := c.Query("chainId")
	accountParam := c.Query("account")
	pageIndexParam, _ := strconv.Atoi(pageIndex)
	PageSizeParam, _ := strconv.Atoi(PageSize)
	chainIdParam, _ := strconv.Atoi(chainId)

	var res *models.JsonRPCAccountNFT
	key := fmt.Sprintf(`NFTs%d-%s-%d-%d`, chainIdParam, accountParam, pageIndexParam, PageSizeParam)
	cacheNFTs, ok := svc.mCache.Get(key)
	if !ok {
		var url string
		for indexScan := range svc.scanInfo {
			for indexUrl := range svc.scanInfo[indexScan].ChainId {
				if svc.scanInfo[indexScan].ChainId[indexUrl] == chainIdParam {
					url = svc.scanInfo[indexScan].ScanUrl[indexUrl]
					break
				}
			}
		}

		if url == "" {
			c.JSON(http.StatusOK, models.Response{
				Code:    200,
				Message: "unsupported chain.",
			})
			return
		}

		var err error

		res, err = utils.AccountNFTPortfolio(accountParam, url, pageIndexParam, PageSizeParam)
		if err != nil {
			oo.LogW("AccountNFTPortfolio err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		svc.mCache.Set(key, res, time.Duration(5)*time.Minute)
	} else {
		res = cacheNFTs.(*models.JsonRPCAccountNFT)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "ok",
		Data:    res.Result,
	})
}

// @Summary account follow
// @Tags account
// @version 0.0.1
// @description account follow
// @Produce json
// @Param request body models.FollowWithSignParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/account/update/follow [post]
func httpUpdateAccountFollow(c *gin.Context) {
	var params models.FollowWithSignParam
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

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", params.Sign.Account).
		Where("followed", params.Params.FollowAccount).Count()
	err = oo.SqlGet(sqlSel, &count)
	if err != nil && err != oo.ErrNoRows {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if params.Params.Status {
		if count == 0 {
			var m = make([]map[string]interface{}, 0)
			var v = make(map[string]interface{})
			v["account"] = params.Sign.Account
			v["followed"] = params.Params.FollowAccount
			v["status"] = 1
			m = append(m, v)
			sqlIns := oo.NewSqler().Table(consts.TbNameAccountFollow).Insert(m)
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
		if count == 1 {
			var v = make(map[string]interface{})
			v["status"] = 1
			sqlUp := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", params.Sign.Account).
				Where("followed", params.Params.FollowAccount).Update(v)
			err = oo.SqlExec(sqlUp)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				c.JSON(http.StatusInternalServerError, models.Response{
					Code:    500,
					Message: "Something went wrong, Please try again later.",
				})
				return
			}
		}
	}

	if !params.Params.Status {
		if count == 0 {
			c.JSON(http.StatusOK, models.Response{
				Code:    http.StatusOK,
				Message: "you are not following.",
			})
			return
		}
		if count == 1 {
			var v = make(map[string]interface{})
			v["status"] = 0
			sqlUp := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", params.Sign.Account).
				Where("followed", params.Params.FollowAccount).Update(v)
			err = oo.SqlExec(sqlUp)
			if err != nil {
				oo.LogW("SQL err: %v", err)
				c.JSON(http.StatusInternalServerError, models.Response{
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
	})
}

// @Summary account following list
// @Tags account
// @version 0.0.1
// @description account following list
// @Produce json
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Param account query  string true "account"
// @Success 200 {object} models.ResAccountFollowPage
// @Router /stpdao/v2/account/following/list [get]
func httpAccountFollowingList(c *gin.Context) {
	accountParam := c.Query("account")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var entities []models.AccountFollowModel
	sqler := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", accountParam).Where("status", 1)

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("update_time DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &entities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResAccountFollow, 0)
	for index := range entities {
		var mutualCount int
		sqlSel := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", entities[index].Followed).
			Where("followed", entities[index].Account).Where("status", 1).Count()
		err = oo.SqlGet(sqlSel, &mutualCount)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		var relation = "following"
		if mutualCount == 1 {
			relation = "mutualFollowing"
		}

		var entity models.AccountModel
		sqlSel = oo.NewSqler().Table(consts.TbNameAccount).Where("account", entities[index].Followed).Select()
		err = oo.SqlGet(sqlSel, &entity)
		if err != nil && err != oo.ErrNoRows {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		data = append(data, models.ResAccountFollow{
			Account:     entities[index].Account,
			FollowTime:  entities[index].UpdateTime,
			Following:   entities[index].Followed,
			AccountLogo: entity.AccountLogo.String,
			Nickname:    entity.Nickname.String,
			Relation:    relation,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAccountFollowPage{
			List:  data,
			Total: total,
		},
	})
}

// @Summary account followers list
// @Tags account
// @version 0.0.1
// @description account followers list
// @Produce json
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Param account query  string true "account"
// @Success 200 {object} models.ResAccountFollowersPage
// @Router /stpdao/v2/account/followers/list [get]
func httpAccountFollowersList(c *gin.Context) {
	accountParam := c.Query("account")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var entities []models.AccountFollowModel
	sqler := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("followed", accountParam).Where("status", 1)

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("update_time DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &entities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResAccountFollowers, 0)
	for index := range entities {
		var mutualCount int
		sqlSel := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", entities[index].Followed).
			Where("followed", entities[index].Account).Where("status", 1).Count()
		err = oo.SqlGet(sqlSel, &mutualCount)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		var relation = "following"
		if mutualCount == 1 {
			relation = "mutualFollowing"
		}

		var entity models.AccountModel
		sqlSel = oo.NewSqler().Table(consts.TbNameAccount).Where("account", entities[index].Account).Select()
		err = oo.SqlGet(sqlSel, &entity)
		if err != nil && err != oo.ErrNoRows {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		data = append(data, models.ResAccountFollowers{
			Account:     entities[index].Followed,
			FollowTime:  entities[index].UpdateTime,
			Followers:   entities[index].Account,
			AccountLogo: entity.AccountLogo.String,
			Nickname:    entity.Nickname.String,
			Relation:    relation,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAccountFollowersPage{
			List:  data,
			Total: total,
		},
	})
}

// @Summary account relation
// @Tags account
// @version 0.0.1
// @description account relation
// @Produce json
// @Param myself query string true "myself"
// @Param others query string true "others"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/account/relation [get]
func httpAccountRelation(c *gin.Context) {
	myselfParam := c.Query("myself")
	othersParam := c.Query("others")

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameAccountFollow).Where("account", myselfParam).
		Where("followed", othersParam).Where("status", 1).Count()
	err := oo.SqlGet(sqlSel, &count)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var isFollowing bool
	if count == 1 {
		isFollowing = true
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    isFollowing,
	})
}

// @Summary update push switch
// @Tags account
// @version 0.0.1
// @description update push switch
// @Produce json
// @Param request body models.UpdateAccountPushSwitchParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/account/push/setting [post]
func httpPushSetting(c *gin.Context) {
	var params models.UpdateAccountPushSwitchParam
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

	pushSwitch := 0b0
	if params.AllDaosICreateOrJoin {
		pushSwitch = pushSwitch | (1 << 0)
	} else {
		pushSwitch = pushSwitch & ^(1 << 0)
	}
	if params.NewDao {
		pushSwitch = pushSwitch | (1 << 1)
	} else {
		pushSwitch = pushSwitch & ^(1 << 1)
	}
	if params.AllDaoAirdrop {
		pushSwitch = pushSwitch | (1 << 2)
	} else {
		pushSwitch = pushSwitch & ^(1 << 2)
	}
	if params.AllDaoProposal {
		pushSwitch = pushSwitch | (1 << 3)
	} else {
		pushSwitch = pushSwitch & ^(1 << 3)
	}

	var v = make(map[string]interface{})
	v["push_switch"] = pushSwitch
	sqler := oo.NewSqler().Table(consts.TbNameAccount).Where("account", params.Sign.Account).Update(v)
	err = oo.SqlExec(sqler)
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
	})
}
