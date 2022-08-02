package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"stp_dao_v2/models"
	"strings"
)

func ScanBlock(currentBlock, url string) (*models.JsonRPCModel, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method":"eth_getLogs",
		"params":[{
			"fromBlock": "%s",
			"toBlock": "%s"
		}]
	}`, currentBlock, currentBlock)

	return jsonRPC(body, url)
}

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

func QueryLatestBlock(url string) (*models.BlockNumber, error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method":"eth_blockNumber",
		"params":[]
	}`)

	return jsonRPCBlock(body, url)
}

func jsonRPCBlock(body, url string) (data *models.BlockNumber, err error) {
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

func QueryDaoInfo(tokenAddress, data, url string) (model *models.JsonRPCInfoModel, err error) {
	body := fmt.Sprintf(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "eth_call",
		"params": [
			{
				"to": "%s",
				"data": "%s"
			},
			"latest"
   	 ]
	}`, tokenAddress, data)
	return jsonInfoRPC(body, url)
}

func jsonInfoRPC(body, url string) (data *models.JsonRPCInfoModel, err error) {
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
