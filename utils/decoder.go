package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// DecodeDistribution
//function createERC20(
//    string memory name_,
//    string memory symbol_,
//    string memory logoUrl_,
//    uint8 decimal_,
//    uint256 totalSupply_,
//    DistributionParam[] calldata distributions_
//)
//struct DistributionParam {
//    address recipient;
//    uint256 amount;
//    uint256 lockDate;
//}
//   /**

type Distribution struct {
	Recipient string
	Amount    string
	LockDate  uint64
}

func DecodeDistribution(input string) (distributions []Distribution, err error) {
	// method id: 0xacdd39b6
	if !strings.HasPrefix(input, "0xacdd39b6") {
		return nil, errors.New("unknown data")
	}
	input = strings.Replace(input, "0xacdd39b6", "", -1)

	// find distributions index
	var tmpDataInt *big.Int
	tmpDataInt, err = parseUint256(input[64*5 : 64*5+64])
	if err != nil {
		return nil, err
	}

	offset := int(tmpDataInt.Uint64() * 2)
	if len(input) < offset+64 {
		err = errors.New("invalid data")
		return nil, err
	}

	tmpDataInt, err = parseUint256(input[offset : offset+64])
	if err != nil {
		return nil, err
	}

	inputIndex := offset + 64
	length := int(tmpDataInt.Uint64())
	for index := 0; index < length; index++ {
		var amount, lockDate *big.Int
		recipient := fmt.Sprintf("0x%s", input[inputIndex+24:inputIndex+64])
		amount, err = parseUint256(input[inputIndex+64 : inputIndex+64+64])
		if err != nil {
			return nil, err
		}
		lockDate, err = parseUint256(input[inputIndex+64+64 : inputIndex+64+64+64])
		if err != nil {
			return nil, err
		}

		inputIndex += 3 * 64
		distributions = append(distributions, Distribution{
			Recipient: recipient,
			Amount:    amount.String(),
			LockDate:  lockDate.Uint64(),
		})
	}

	return distributions, nil
}

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
				result = append(result, fmt.Sprintf("0x%s", input[inputIndex+24:inputIndex+64]))
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
				var b uint64
				b, err = strconv.ParseUint(input[inputIndex:inputIndex+64], 16, 64)
				if err != nil {
					break
				}
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
				var data uint64
				data, err = strconv.ParseUint(input[inputIndex:inputIndex+64], 16, 64)
				if err != nil {
					break
				}

				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "address[]" {
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
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					data = append(data, fmt.Sprintf("0x%s", input[dataOffset+64+(index*64)+24:dataOffset+64+(index*64)+64]))
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "uint256[]" {
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
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []*big.Int
				for index := 0; index < dataLength; index++ {
					var datum *big.Int
					datum, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					if err != nil {
						break
					}
					data = append(data, datum)
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "string[]" {
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
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					var offset, length *big.Int
					offset, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					length, err = parseUint256(input[dataOffset+64+int(offset.Int64()*2) : dataOffset+64+int(offset.Int64()*2)+64])
					var datum string
					datum, err = parseString(input[dataOffset+64+int(offset.Int64()*2)+64 : dataOffset+64+int(offset.Int64()*2)+64+int(length.Int64()*2)])
					data = append(data, datum)
				}
				result = append(result, data)
				inputIndex += 64
			} else if outputParameter == "bytes[]" {
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
				dataLength := int(tmpDataInt.Uint64())
				if len(input) < dataOffset+64+(dataLength*64) {
					err = errors.New("invalid data")
					break
				}
				var data []string
				for index := 0; index < dataLength; index++ {
					var offset, length *big.Int
					offset, err = parseUint256(input[dataOffset+64+(index*64) : dataOffset+64+(index*64)+64])
					length, err = parseUint256(input[dataOffset+64+int(offset.Int64()*2) : dataOffset+64+int(offset.Int64()*2)+64])
					var datum string
					datum = input[dataOffset+64+int(offset.Int64()*2)+64 : dataOffset+64+int(offset.Int64()*2)+64+int(length.Int64()*2)]
					data = append(data, fmt.Sprintf("0x%s", datum))
				}
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
