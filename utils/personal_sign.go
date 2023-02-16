package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func CheckSignMessageTimestamp(message string) bool {

	arr := strings.Split(message, " ")
	if len(arr) != 5 {
		return false
	}

	now := time.Now().Unix()
	gap := int64(1 * 60 * 60)

	num, _ := strconv.ParseInt(arr[4], 10, 64)

	if num < now-gap || num > now+gap {
		return false
	}

	return true
}

func CheckAdminSignMessageTimestamp(timestamp int64) bool {
	now := time.Now().Unix()
	if timestamp < now {
		return false
	}

	return true
}

func eip191MessageHash(message string) common.Hash {
	buf := []byte(message)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(buf), buf)
	return crypto.Keccak256Hash([]byte(msg))
}

func PersonalSign(pri_key, message string) (address, signature string, err error) {
	priv, err := crypto.HexToECDSA(pri_key)
	if nil != err {
		err = fmt.Errorf("HexToECDSA(%s) %v", pri_key, err)
		return
	}

	var (
		msg = eip191MessageHash(message).Bytes()
	)
	sig, err := crypto.Sign(msg, priv)
	if nil != err {
		err = fmt.Errorf("crypto.Sign %v", err)
		return
	}

	signature = common.Bytes2Hex(sig)
	address = crypto.PubkeyToAddress(priv.PublicKey).Hex()

	return
}

func CheckPersonalSign(message, address, signature string) (ret bool, err error) {
	var (
		msg = eip191MessageHash(message).Bytes()
		sig = common.FromHex(signature)
	)
	if len(sig) != 65 {
		err = fmt.Errorf("invalid signature len")
		return
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	if sig[64] != 0 && sig[64] != 1 {
		err = fmt.Errorf("invalid signature v")
		return
	}

	pub, err := crypto.SigToPub(msg, sig)
	if nil != err {
		return
	}

	ret = bytes.Equal(
		common.HexToAddress(address).Bytes(),
		crypto.PubkeyToAddress(*pub).Bytes(),
	)

	return
}

func SignMessage(message, secret string) (string, error) {

	by, _ := hex.DecodeString(message)
	sha3 := crypto.Keccak256(by)
	buf := []byte(fmt.Sprintf("\u0019Ethereum Signed Message:\n%d%s", len(sha3), sha3))
	sha3 = crypto.Keccak256(buf)

	s, _ := hex.DecodeString(secret)
	priv, err := crypto.ToECDSA(s)
	if err != nil {
		err = fmt.Errorf("ToECDSA(%s) err:%v", s, err)
		return "", err
	}

	sig, err := crypto.Sign(sha3, priv)
	if err != nil {
		err = fmt.Errorf("crypto.Sign err:%v", err)
		return "", err
	}
	sig1 := big.NewInt(0)
	sig1.Add(big.NewInt(0).SetBytes(sig), big.NewInt(27))
	sig = sig1.Bytes()

	signature := common.Bytes2Hex(sig)

	return fmt.Sprintf("%0130s", signature), nil
}
