package controller

import (
	"encoding/json"
	"fmt"
	solTree "github.com/0xKiwi/sol-merkle-tree-go"
	oo "github.com/Anna2024/liboo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
)

// @Summary save airdrop address
// @Tags Airdrop
// @version 0.0.1
// @description save airdrop address
// @Produce json
// @Param request body models.AirdropAddressParam true "request"
// @Success 200 {object} models.ResAirdropId
// @Router /stpdao/v2/airdrop/address [post]
func httpSaveAirdropAddress(c *gin.Context) {
	var params models.AirdropAddressParam
	err := c.ShouldBindJSON(&params)
	if err != nil || len(params.Address) != len(params.Amount) {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	encoded, err := json.Marshal(params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Json Marshal Failed.",
		})
		return
	}

	var m = make([]map[string]interface{}, 0)
	var v = make(map[string]interface{})
	v["content"] = string(encoded)
	m = append(m, v)

	sqlIns := oo.NewSqler().Table(consts.TbNameAirdropAddress).Insert(m)
	err = oo.SqlExec(sqlIns)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var airdropId int
	sqlSel := fmt.Sprintf(`SELECT LAST_INSERT_ID()`)
	err = oo.SqlGet(sqlSel, &airdropId)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResAirdropId{
			AirdropIdId: airdropId,
		},
	})
}

// @Summary claim airdrop address
// @Tags Airdrop
// @version 0.0.1
// @description claim airdrop address
// @Produce json
// @Param address query string true "address"
// @Param id query int true "id"
// @Success 200 {object} models.ClaimInfo
// @Router /stpdao/v2/airdrop/proof [get]
func httpClaimAirdrop(c *gin.Context) {
	idParam := c.Query("id")
	addressParam := c.Query("address")

	var entity []models.AirdropAddressModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAirdropAddress).Where("id", idParam).Select()
	err := oo.SqlSelect(sqlSel, &entity)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data models.AirdropAddress
	err = json.Unmarshal([]byte(entity[0].Content), &data)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Json Unmarshal Failed.",
		})
		return
	}

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
			Data: models.ClaimInfo{
				Amount: "",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ClaimInfo{
			Index:  claimInfo.Index,
			Amount: claimInfo.Amount,
			Proof:  claimInfo.Proof,
		},
	})
}
