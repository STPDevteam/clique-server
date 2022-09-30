package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	solTree "github.com/0xKiwi/sol-merkle-tree-go"
	oo "github.com/Anna2024/liboo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"os"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
)

// @Summary create airdrop
// @Tags Airdrop
// @version 0.0.1
// @description create airdrop
// @Produce json
// @Param request body models.CreateAirdropParam true "request"
// @Success 200 {object} models.ResAirdropId
// @Router /stpdao/v2/airdrop/create [post]
func (svc *Service) httpCreateAirdrop(c *gin.Context) {
	var params models.CreateAirdropParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}
	if params.StartTime >= params.EndTime || params.EndTime >= params.AirdropStartTime || params.AirdropStartTime >= params.AirdropEndTime {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	if !checkAirdropAdminAndTimestamp(&params.Sign) {
		oo.LogD("SignData err not auth")
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Data:    models.ResResult{Success: false},
			Message: "SignData err not auth",
		})
		return
	}

	var approve bool
	sqlSel := oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", params.Sign.ChainId).Where("dao_address", params.Sign.DaoAddress).Select("approve")
	err = oo.SqlGet(sqlSel, &approve)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	if !approve {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusOK,
			Message: "not approved",
		})
		return
	}

	encoded, err := json.Marshal(params.CollectInformation)
	if err != nil {
		oo.LogW("json.Marshal %v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Json Marshal Failed.",
		})
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["creator"] = params.Sign.Account
	v["chain_id"] = params.Sign.ChainId
	v["dao_address"] = params.Sign.DaoAddress
	v["title"] = params.Title
	v["airdrop_address"] = ""
	v["description"] = params.Description
	v["collect_information"] = string(encoded)
	v["token_chain_id"] = params.TokenChainId
	v["token_address"] = params.TokenAddress
	v["max_airdrop_amount"] = params.MaxAirdropAmount
	v["start_time"] = params.StartTime
	v["end_time"] = params.EndTime
	v["airdrop_start_time"] = params.AirdropStartTime
	v["airdrop_end_time"] = params.AirdropEndTime
	m = append(m, v)
	sqlIn := oo.NewSqler().Table(consts.TbNameAirdrop).Insert(m)
	err = oo.SqlExec(sqlIn)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var airdropId int
	sqlSel = fmt.Sprintf(`SELECT LAST_INSERT_ID()`)
	err = oo.SqlGet(sqlSel, &airdropId)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	// TokenChainId+Account+AirdropId
	resAccount := strings.TrimPrefix(params.Sign.Account, "0x")
	resTokenChainId := fmt.Sprintf("%064x", params.TokenChainId)
	resAirdropId := fmt.Sprintf("%064x", airdropId)
	message := fmt.Sprintf("%s%s%s", resTokenChainId, resAccount, resAirdropId)
	signature, err := utils.SignMessage(message, svc.appConfig.SignMessagePriKey)
	if err != nil {
		oo.LogW("SignMessage err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Signature err",
		})
		return
	}
	signature = fmt.Sprintf("0x%s", signature)

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAirdropId{
			AirdropId: airdropId,
			Signature: signature,
		},
	})

}

// @Summary airdrop need collect information
// @Tags Airdrop
// @version 0.0.1
// @description airdrop need collect information
// @Produce json
// @Param id query int true "id"
// @Success 200 {object} models.CollectInformationInfo
// @Router /stpdao/v2/airdrop/collect [get]
func httpCollectInformation(c *gin.Context) {
	idParam := c.Query("id")

	var entity []models.AirdropModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", idParam).Select()
	err := oo.SqlSelect(sqlSel, &entity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data models.CollectInformationInfo
	err = json.Unmarshal([]byte(entity[0].CollectInformation), &data)
	if err != nil {
		oo.LogW("json.Unmarshal %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var addressNum = 0
	var addressArray models.AirdropAddressArray
	if len(entity[0].AirdropAddress) != 0 {
		err = json.Unmarshal([]byte(entity[0].AirdropAddress), &addressArray)
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Json Unmarshal Failed.",
			})
			return
		}
		addressNum = len(addressArray.Address)
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAirdropInfo{
			Creator:          entity[0].Creator,
			ChainId:          entity[0].ChainId,
			DaoAddress:       entity[0].DaoAddress,
			Title:            entity[0].Title,
			Description:      entity[0].Description,
			TokenChainId:     entity[0].TokenChainId,
			TokenAddress:     entity[0].TokenAddress,
			StartTime:        entity[0].StartTime,
			EndTime:          entity[0].EndTime,
			AirdropStartTime: entity[0].AirdropStartTime,
			AirdropEndTime:   entity[0].AirdropEndTime,
			AddressNum:       addressNum,
			Collect:          data,
		},
	})

}

// @Summary airdrop collect user information
// @Tags Airdrop
// @version 0.0.1
// @description airdrop collect user information
// @Produce json
// @Param request body models.UserInformationParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/airdrop/save/user [post]
func httpSaveUserInformation(c *gin.Context) {
	var params models.UserInformationParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var endTime int64
	sqlSel := oo.NewSqler().Table(consts.TbNameAirdrop).Select("end_time")
	err = oo.SqlGet(sqlSel, &endTime)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	if endTime < time.Now().Unix() {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusOK,
			Message: "event has ended",
			Data: models.ResResult{
				Success: false,
			},
		})
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["airdrop_id"] = params.AirdropId
	v["account"] = params.Account
	v["twitter"] = params.Twitter
	v["telegram"] = params.Telegram
	v["email"] = params.Email
	v["txid"] = params.TXID
	v["other"] = params.Other
	v["discord_username"] = params.DiscordUsername
	m = append(m, v)
	sqlIns := oo.NewSqler().Table(consts.TbNameAirdropUserSubmit).Insert(m)
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

// @Summary airdrop download user information
// @Tags Airdrop
// @version 0.0.1
// @description airdrop download user information
// @Produce json
// @Param id query int true "id"
// @Router /stpdao/v2/airdrop/user/download [get]
func httpDownloadUserInformation(c *gin.Context) {
	idParam := c.Query("id")

	var userEntities []models.AirdropUserSubmit
	sqlSel := oo.NewSqler().Table(consts.TbNameAirdropUserSubmit).Where("airdrop_id", idParam).Select()
	err := oo.SqlSelect(sqlSel, &userEntities)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var airdropEntity []models.AirdropModel
	sqlSel = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", idParam).Select()
	err = oo.SqlSelect(sqlSel, &airdropEntity)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data models.CollectInformationInfo
	err = json.Unmarshal([]byte(airdropEntity[0].CollectInformation), &data)
	if err != nil {
		oo.LogW("json.Unmarshal %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	if airdropEntity[0].EndTime <= time.Now().Unix() {
		err = os.MkdirAll("./file", os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		fileName := fmt.Sprintf("airdrop_users_%d.csv", time.Now().Unix())
		path := fmt.Sprintf("./file/%s", fileName)
		f, err := os.Create(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		defer f.Close()

		writer := csv.NewWriter(f)
		writer.Write([]string{"account", "discord", "twitter", "telegram", "email", "TXID", data.Other})

		var allArray = make([][]string, 0)
		for index := range userEntities {
			var array = make([]string, 0)
			u := userEntities[index]
			array = append(array, u.Account, u.DiscordUsername.String, u.Twitter.String, u.Telegram.String, u.Email.String, u.TXID.String, u.Other.String)
			allArray = append(allArray, array)
		}
		writer.WriteAll(allArray)
		writer.Flush()

		if err = writer.Error(); err != nil {
			oo.LogW("error writing csv:", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		file, err := os.Open(path)
		if err != nil {
			oo.LogW("file open err:", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		defer file.Close()

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(path)

	}
}

// @Summary save airdrop address
// @Tags Airdrop
// @version 0.0.1
// @description save airdrop address
// @Produce json
// @Param request body models.AirdropAdminSignData true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/airdrop/address [post]
func httpSaveAirdropAddress(c *gin.Context) {
	var params models.AirdropAdminSignData
	err := c.ShouldBindJSON(&params)
	if err != nil || len(params.Array.Address) != len(params.Array.Amount) {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	if !checkAirdropAdminAndTimestamp(&params) {
		oo.LogD("SignData err not auth")
		c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Data:    models.ResResult{Success: false},
			Message: "SignData err not auth",
		})
		return
	}

	encoded, err := json.Marshal(params.Array)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Json Marshal Failed.",
		})
		return
	}

	var v = make(map[string]interface{})
	v["airdrop_address"] = string(encoded)
	sqlUp := oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", params.AirdropId).Update(v)
	err = oo.SqlExec(sqlUp)
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

// @Summary claim airdrop address
// @Tags Airdrop
// @version 0.0.1
// @description claim airdrop address
// @Produce json
// @Param address query string true "address"
// @Param id query int true "id"
// @Success 200 {object} models.ResProof
// @Router /stpdao/v2/airdrop/proof [get]
func httpClaimAirdrop(c *gin.Context) {
	idParam := c.Query("id")
	addressParam := c.Query("address")

	var entity []models.AirdropModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", idParam).Select()
	err := oo.SqlSelect(sqlSel, &entity)
	if err != nil || len(entity[0].AirdropAddress) == 0 {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data models.AirdropAddressArray
	err = json.Unmarshal([]byte(entity[0].AirdropAddress), &data)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Json Unmarshal Failed.",
		})
		return
	}

	var totalAmount = new(big.Int)
	var addressLength = len(data.Address)
	var addressData = make([]models.AddressData, addressLength)
	for index := 0; index < addressLength; index++ {
		amount, err := utils.Dec2BigInt(data.Amount[index])
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		addressData[index] = models.AddressData{
			Id:      uint64(index),
			Amount:  amount,
			Address: common.HexToAddress(data.Address[index]),
		}

		totalAmount.Add(totalAmount, amount)
	}

	var nodes = make([][]byte, addressLength)
	for index, model := range addressData {
		packed := append(
			common.LeftPadBytes(big.NewInt(0).SetInt64(int64(index)).Bytes(), 32),
			append(
				model.Address.Bytes(),
				common.LeftPadBytes(model.Amount.Bytes(), 32)...,
			)...,
		)

		nodes[index] = crypto.Keccak256(packed)
	}

	var (
		merkleTree  *solTree.MerkleTree
		addrToProof map[string]models.ClaimInfo
	)
	merkleTree, err = solTree.GenerateTreeFromHashedItems(nodes)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	addrToProof = make(map[string]models.ClaimInfo, addressLength)
	for index, model := range addressData {
		proof, err := merkleTree.MerkleProof(nodes[index])
		if err != nil {
			oo.LogW("%v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		addrToProof[strings.ToLower(model.Address.String())] = models.ClaimInfo{
			Index:  uint64(index),
			Amount: model.Amount.String(),
			Proof:  utils.StringArrayFrom2DBytes(proof),
		}
	}

	claimInfo, ok := addrToProof[strings.ToLower(addressParam)]
	if !ok {
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusOK,
			Message: "ok",
			Data: models.ResProof{
				AirdropTotalAmount: totalAmount.String(),
				AirdropNumber:      addressLength,
				Title:              entity[0].Title,
				Amount:             "",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResProof{
			AirdropTotalAmount: totalAmount.String(),
			AirdropNumber:      addressLength,
			Title:              entity[0].Title,
			Index:              claimInfo.Index,
			Amount:             claimInfo.Amount,
			Proof:              claimInfo.Proof,
		},
	})
}
