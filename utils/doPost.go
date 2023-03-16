package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"stp_dao_v2/models"
	"strings"
)

func ScanBlock(currentBlock, url string) (*models.JsonRPCScanBlockModel, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method":"eth_getLogs",
		"params":[{
			"fromBlock": "%s",
			"toBlock": "%s"
		}]
	}`, currentBlock, currentBlock)

	return jsonRPCScanBlock(body, url)
}

func jsonRPCScanBlock(body, url string) (data *models.JsonRPCScanBlockModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func QueryLatestBlock(url string) (*models.JsonRPCModel, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method":"eth_blockNumber",
		"params":[]
	}`)

	return jsonRPC(body, url)
}

func QueryTimesTamp(block, url string) (model *models.JsonRPCTimesTampModel, err error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": [
			"%s",
			false
		],
		"id": 1
	}`, block)
	return jsonTimesTampRPC(body, url)
}

func jsonTimesTampRPC(body, url string) (data *models.JsonRPCTimesTampModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTransactionByHashFrom(hash, url string) (model *models.JsonRPCTransactionByHashModel, err error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getTransactionByHash",
		"params": [
			"%s"
		],
		"id": 1
	}`, hash)
	return jsonTransactionByHashRPC(body, url)
}

func jsonTransactionByHashRPC(body, url string) (data *models.JsonRPCTransactionByHashModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReserveToken(hash, url string) (model *models.JsonRPCReserveTokenModel, err error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "eth_getTransactionByHash",
		"params": [
			"%s"
		],
		"id": 1
	}`, hash)
	return jsonReserveTokenRPC(body, url)
}

func jsonReserveTokenRPC(body, url string) (data *models.JsonRPCReserveTokenModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func QueryBalance(tokenAddress, accountAddress, url string) (model *models.JsonRPCBalanceModel, err error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "erc20_balance",
		"params": {
			"contractAddress":"%s",
			"accountAddress":"%s"
		}
	}`, tokenAddress, accountAddress)
	return jsonBalanceRPC(body, url)
}

func QuerySpecifyBalance(tokenAddress, accountAddress, url string, blockNumber int64) (model *models.JsonRPCBalanceModel, err error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "erc20_balance",
		"params": {
			"contractAddress":"%s",
			"accountAddress":"%s",
			"blockNumber":%d
		}
	}`, tokenAddress, accountAddress, blockNumber)
	return jsonBalanceRPC(body, url)
}

func jsonBalanceRPC(body, url string) (data *models.JsonRPCBalanceModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//
//func QueryMethodEthCall(to, data, url string) (*models.JsonRPCModel, error) {
//	body := fmt.Sprintf(`{
//		"id": 1,
//		"jsonrpc":"2.0",
//		"method": "eth_call",
//		"params": [{"to":"%s","data":"%s"}, "latest"]
//	}`, to, data)
//
//	return jsonRPC(body, url)
//}

func jsonRPC(body, url string) (data *models.JsonRPCModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DoPost(url string, contentType string, body string) (res []byte, err error) {
	resp, err := http.Post(url, contentType, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func DoGet(url string) (res []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func GetTokensId(url string) ([]models.TokensInfo, error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	data := make([]models.TokensInfo, 0)
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTokensPrice(url string) (data map[string]map[string]float64, err error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetTokenImg(url string) (data *models.TokenImg, err error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetV1ProposalHistory(url string) (data *models.V1ProposalHistory, err error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetKlaytnBlock(url string) (data *models.KlaytnBlock, err error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func QueryMethodEthCallByTag(to, data, url, tag string) (*models.JsonRPCModel, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method": "eth_call",
		"params": [{"to":"%s","data":"%s"}, "%s"]
	}`, to, data, tag)

	return jsonRPC(body, url)
}

func QueryMethodEthCall(to, data, url string) (*models.JsonRPCModel, error) {
	return QueryMethodEthCallByTag(to, data, url, "latest")
}

func QueryErc20TokenHolders(token, url string) (*models.JsonRPCTokenHoldersModel, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "erc20_tokenHolders",
		"params": {
			"contractAddress":"%s",
			"pageSize":1,
			"pageIndex":1
		}
	}`, token)

	return jsonRPCTokenHolders(body, url)
}

func jsonRPCTokenHolders(body, url string) (data *models.JsonRPCTokenHoldersModel, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetBlockNumberFromTimestamp(url string) (data *models.JsonRPCGetBlockNumber, err error) {
	res, err := DoGet(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func AccountNFTPortfolio(account, url string, pageIndex, pageSize int) (*models.JsonRPCAccountNFT, error) {
	body := fmt.Sprintf(`{
		"jsonrpc": "2.0",
			"id": 0,
			"method": "account_nftPortfolio",
			"params": {
			"accountAddress": "%s",
				"pageIndex": %d,
				"pageSize": %d
		}
	}`, account, pageIndex, pageSize)

	return jsonRPCAccountNFTPortfolio(body, url)
}

func jsonRPCAccountNFTPortfolio(body, url string) (data *models.JsonRPCAccountNFT, err error) {
	res, err := DoPost(
		url,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
