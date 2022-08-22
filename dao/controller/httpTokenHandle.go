package controller

import (
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strconv"
)

// @Summary query Token list
// @Tags Token
// @version 0.0.1
// @description query Token list
// @Produce json
// @Param chainId query  string false "chainId"
// @Param creator query  string false "creator"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResTokenListPage
// @Router /stpdao/v2/token/list [get]
func httpTokenList(c *gin.Context) {
	offset := c.Query("offset")
	count := c.Query("count")
	offsetParam, _ := strconv.Atoi(offset)
	countParam, _ := strconv.Atoi(count)
	creatorParam := c.Query("creator")
	chainIdParam := c.Query("chainId")

	var total uint64
	var listEntities []models.EventHistoricalModel
	sqlSel := oo.NewSqler().Table(consts.TbNameEventHistorical).Where("event_type", consts.EvCreateERC20)
	if creatorParam != "" {
		creatorParam = utils.FixTo0x64String(creatorParam)
		sqlSel = sqlSel.Where("topic1", creatorParam)
	}
	if chainIdParam != "" {
		sqlSel = sqlSel.Where("chain_id", chainIdParam)
	}

	sqlCopy := *sqlSel
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqlSel
		sqlStr = sqlCopy.Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &listEntities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResTokenList, 0)
	for index := range listEntities {
		tokenAddress := utils.FixTo0x40String(listEntities[index].Data)
		contractAddress := listEntities[index].Address
		chainId := listEntities[index].ChainId

		var names []string
		sqlName := oo.NewSqler().Table(consts.TbNameDao).
			Where("token_address", tokenAddress).
			Where("token_chain_id", chainId).Select("dao_name")
		err = oo.SqlSelect(sqlName, &names)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		var totalSupply string
		sqlTotal := oo.NewSqler().Table(consts.TbNameHolderData).
			Where("token_address", tokenAddress).
			Where("holder_address", consts.ZeroAddress0x40).
			Where("chain_id", chainId).Select("balance")
		err = oo.SqlGet(sqlTotal, &totalSupply)
		if err != nil && err != oo.ErrNoRows {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		data = append(data, models.ResTokenList{
			TokenAddress:    tokenAddress,
			ContractAddress: contractAddress,
			ChainId:         chainId,
			DaoName:         names,
			TotalSupply:     totalSupply,
		})

	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResTokenListPage{
			List:  data,
			Total: total,
		},
	})
}

// @Summary query Token img
// @Tags Token
// @version 0.0.1
// @description query Token img
// @Produce json
// @Param chainId query int true "chainId"
// @Param tokenAddress query string true "tokenAddress"
// @Success 200 {object} models.ResTokenImg
// @Router /stpdao/v2/token/img [get]
func httpTokenImg(c *gin.Context) {
	chainId := c.Query("chainId")
	chainIdParam, _ := strconv.Atoi(chainId)
	tokenAddressParam := c.Query("tokenAddress")

	var entity []models.TokensImgModel
	sqlSel := oo.NewSqler().Table(consts.TbNameTokensImg).
		Where("chain_id", chainIdParam).
		Where("token_address", tokenAddressParam).Select()
	err := oo.SqlSelect(sqlSel, &entity)
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
		Data: models.ResTokenImg{
			ChainId:      chainIdParam,
			TokenAddress: tokenAddressParam,
			Thumb:        entity[0].Thumb,
			Small:        entity[0].Small,
			Large:        entity[0].Large,
		},
	})

}
