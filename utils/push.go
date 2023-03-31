package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"strings"
)

type PushFeed struct {
	PayloadId uint64 `json:"payload_id"`
	Sender    string `json:"sender"`
	Epoch     string `json:"epoch"`
	Payload   struct {
		Data struct {
			App   string `json:"app"`
			Sid   string `json:"sid"`
			Url   string `json:"url"`
			Acta  string `json:"acta"`
			Aimg  string `json:"aimg"`
			Asub  string `json:"asub"`
			Icon  string `json:"icon"`
			Type  uint64 `json:"type"`
			Epoch string `json:"epoch"`
			//ETime string `json:"etime"`
			Hidden string `json:"hidden"`
			//Sectype string `json:"sectype"`
			//AdditionalMeta string `json:"additionalMeta"`
		} `json:"data"`
		Recipients   string `json:"recipients"`
		Notification struct {
			Body  string `json:"body"`
			Title string `json:"title"`
		} `json:"notification"`
		VerificationProof string `json:"verificationProof"`
	} `json:"payload"`
	Source string `josn:"source"`
	//ETime string `json:"etime"`
}

type PageFeeds struct {
	Feeds     []PushFeed `json:"feeds"`
	Itemcount uint64     `json:"itemcount"`
}

type PushAPI struct {
	EndPoint       string
	Types          apitypes.Types
	ChainId        int64
	Source         string
	ChannelAddress string
	Domain         apitypes.TypedDataDomain
	Signer         *ecdsa.PrivateKey
}

func NewPushAPI(endPoint, source, channelAddress string, chainId int64, signer string) (api *PushAPI, err error) {
	api = &PushAPI{
		EndPoint: endPoint,
		Types: apitypes.Types{
			"EIP712Domain": {
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "chainId",
					Type: "uint256",
				},
				{
					Name: "verifyingContract",
					Type: "address",
				},
			},
			"Data": {
				{
					Name: "data",
					Type: "string",
				},
			},
		},
		ChainId:        chainId,
		Source:         source,
		ChannelAddress: channelAddress,
		Domain: apitypes.TypedDataDomain{
			Name:              "EPNS COMM V1",
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: "0xb3971BCef2D791bc4027BbfedFb47319A4AAaaAa",
		},
	}

	var bytes []byte
	if bytes, err = hex.DecodeString(signer); err != nil {
		return nil, err
	}
	if api.Signer, err = crypto.ToECDSA(bytes); err != nil {
		return nil, err
	}
	return api, nil
}

func (api *PushAPI) GetAddress() string {
	return crypto.PubkeyToAddress(api.Signer.PublicKey).String()
}

func (api *PushAPI) GetCAIPAddress(address string) string {
	return fmt.Sprintf("eip155:%d:%s", api.ChainId, address)
}

func (api *PushAPI) GetChannel() string {
	return strings.ToLower(fmt.Sprintf("eip155:%s", api.ChannelAddress))
}

func (api *PushAPI) GetChannelCAIPAddress() string {
	return fmt.Sprintf("eip155:%d:%s", api.ChainId, api.ChannelAddress)
}

func (api *PushAPI) GetFeeds(page, pageSize uint64) (ret *PageFeeds, err error) {
	var url = fmt.Sprintf(
		"%s/v1/users/%s/feeds?page=%d&limit=%d",
		api.EndPoint,
		api.GetCAIPAddress(api.GetAddress()),
		page,
		pageSize,
	)
	var respBytes []byte
	if respBytes, err = DoGet(url); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(respBytes, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (api *PushAPI) getRecipients(notificationType int, recipients *[]string) (interface{}, error) {
	if notificationType == 1 {
		return api.GetChannelCAIPAddress(), nil
	}
	if notificationType == 4 {
		if recipients != nil {
			var ret []string = make([]string, 0)
			for _, recipient := range *recipients {
				ret = append(ret, api.GetCAIPAddress(recipient))
			}
			return ret, nil
		}
	}
	return nil, fmt.Errorf("recipients was empty")
}

// 1: broadcast; 4: subset;
func (api *PushAPI) SendNotification(notificationType int, uid string, title, body string, recipients *[]string) (err error) {
	if !oo.InArray(notificationType, []int{1, 4}) {
		return fmt.Errorf("notification type was unexpected")
	}

	var recipientObject interface{}
	if recipientObject, err = api.getRecipients(notificationType, recipients); err != nil {
		return err
	}

	var msgBytes []byte
	if msgBytes, err = json.Marshal(map[string]interface{}{
		"notification": map[string]string{
			"title": title,
			"body":  body,
		},
		"data": map[string]string{
			"acta": "",
			"aimg": "",
			"amsg": "",
			"asub": "",
			"type": fmt.Sprintf("%d", notificationType),
		},
		"recipients": recipientObject,
	}); err != nil {
		return err
	}
	data := apitypes.TypedData{
		Types:  api.Types,
		Domain: api.Domain,
		Message: map[string]interface{}{
			"data": fmt.Sprintf("2+%s", string(msgBytes)),
		},
	}
	var domainSeparator, typedDataHash hexutil.Bytes
	if domainSeparator, err = data.HashStruct("EIP712Domain", data.Domain.Map()); err != nil {
		return err
	}
	if typedDataHash, err = data.HashStruct("Data", data.Message); err != nil {
		return err
	}
	dataHash := crypto.Keccak256([]byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))))

	var signature []byte
	if signature, err = crypto.Sign(dataHash, api.Signer); err != nil {
		return err
	}

	if signature[64] < 27 {
		signature[64] += 27
	}

	bodyStruct := map[string]interface{}{
		"verificationProof": fmt.Sprintf("eip712v2:0x%0130s::uid::%s", common.Bytes2Hex(signature), uid),
		"identity":          data.Message["data"],
		"sender":            api.GetChannelCAIPAddress(),
		"source":            api.Source,
		"recipient":         api.GetChannelCAIPAddress(),
	}
	var bodyBytes []byte
	if bodyBytes, err = json.Marshal(bodyStruct); err != nil {
		return err
	}

	_, err = DoPost(
		fmt.Sprintf("%s/v1/payloads/", api.EndPoint),
		"application/json",
		string(bodyBytes),
	)
	return err
}
