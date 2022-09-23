package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

func Decode(outputParameters []string, input string) (result []interface{}, err error) {
	if outputParameters == nil || len(outputParameters) == 0 {
		result = append(result, input)
	} else {
		inputIndex := 0
		for _, outputParameter := range outputParameters {
			if outputParameter == "uint256" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var data *big.Int
				data, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "address" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				result = append(result, input[inputIndex+24:inputIndex+64])
				inputIndex += 64
			} else if outputParameter == "string" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var tmpDataInt *big.Int
				tmpDataInt, err = parseUint256(input[inputIndex : inputIndex+64])
				if err != nil {
					break
				}
				dataOffset := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64 {
					err = errors.New("invalid data")
					break
				}
				tmpDataInt, err = parseUint256(input[dataOffset : dataOffset+64])
				if err != nil {
					break
				}
				dataLength := int(tmpDataInt.Uint64() * 2)
				if len(input) < dataOffset+64+dataLength {
					err = errors.New("invalid data")
					break
				}
				var data string
				data, err = parseString(input[dataOffset+64 : dataOffset+64+dataLength])
				if err != nil {
					break
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "bool" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var b int
				b, _ = Hex2Dec(input[inputIndex : inputIndex+64])
				var data bool
				if b == 0 {
					data = false
				} else if b == 1 {
					data = true
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "uint8" {
				if len(input) < inputIndex+64 {
					err = errors.New("invalid data")
					break
				}
				var data int
				data, _ = Hex2Dec(input[inputIndex : inputIndex+64])

				result = append(result, data)
				inputIndex += 64
			} else {
				err = errors.New(fmt.Sprintf("unsupported %s data", outputParameter))
				break
			}
		}
	}

	return result, err
}

func parseUint256(input string) (*big.Int, error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(bytes), nil
}

func parseString(input string) (string, error) {
	bytes, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	result := string(bytes)

	return result, nil
}
