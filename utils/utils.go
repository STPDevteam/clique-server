package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func Keccak256(str string) string {
	b := []byte(str)
	hash := crypto.Keccak256(b)
	return hex.EncodeToString(hash)
}

func Hex2Dec(hex string) int {
	n, err := strconv.ParseUint(strings.TrimPrefix(hex, "0x"), 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

func Hex2Int64(hex string) (dec int64, err error) {
	r, err := regexp.Compile("0x[0-9a-fA-F]+")
	if err != nil {
		return 0, err
	}
	if !r.MatchString(hex) {
		return 0, fmt.Errorf("not a hex string: %s", hex)
	}

	dec, err = strconv.ParseInt(hex[2:], 16, 64)
	return dec, err
}

func Hex2BigInt(hex string) (dec *big.Int, err error) {
	r, err := regexp.Compile("0x[0-9a-fA-F]+")
	if err != nil {
		return nil, err
	}
	if !r.MatchString(hex) {
		return nil, fmt.Errorf("not a hex string: %s", hex)
	}

	dec, ok := new(big.Int).SetString(hex[2:], 16)
	if !ok {
		return nil, fmt.Errorf("failed to parse hex: %s", hex)
	}
	return dec, nil
}

func Dec2BigInt(val string) (dec *big.Int, err error) {
	r, err := regexp.Compile("[0-9]+")
	if err != nil {
		return nil, err
	}
	if !r.MatchString(val) {
		return nil, fmt.Errorf("not a dec string: %s", val)
	}

	dec, ok := new(big.Int).SetString(val, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse ten: %s", val)
	}
	return dec, nil
}

func FixTo0x64String(str string) string {
	return fmt.Sprintf("0x%064s", strings.TrimPrefix(str, "0x"))
}

func FixTo0x40String(str string) string {
	return fmt.Sprintf("0x%040s", strings.Trim(strings.TrimPrefix(str, "0x"), "0"))
}

func GenerateUuid() string {
	uuidWithHyphen := uuid.NewRandom()
	//uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuidWithHyphen.String()
}
