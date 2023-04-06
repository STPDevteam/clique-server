package utils

import (
	"fmt"
	solTree "github.com/0xKiwi/sol-merkle-tree-go"
	oo "github.com/Anna2024/liboo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"stp_dao_v2/models"
)

func MerkelTreeRoot(data models.AirdropAddressArray) (string, error) {
	var addressLength = len(data.Address)
	var addressData = make([]models.AddressData, addressLength)
	for index := 0; index < addressLength; index++ {
		amount, err := Dec2BigInt(data.Amount[index])
		if err != nil {
			oo.LogW("%v", err)
			return "", err
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

	var merkleTree *solTree.MerkleTree
	merkleTree, err := solTree.GenerateTreeFromHashedItems(nodes)
	if err != nil {
		oo.LogW("%v", err)
		return "", err
	}
	rootStr := fmt.Sprintf("%#x", merkleTree.Root())
	return rootStr, nil
}
