package utils

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	//0000000000000000000000000000000000000000000000000000000000000020
	//000000000000000000000000000000000000000000000000000000000000003a
	//687474703a2f2f64657661706976322e6d79636c697175652e696f2f73746174
	//69632f313636323031353335393931333334383533302e737667000000000000

	var input = "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003a687474703a2f2f64657661706976322e6d79636c697175652e696f2f7374617469632f313636323031353335393931333334383533302e737667000000000000"
	var outputParameters []string
	outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")
	//outputParameters = append(outputParameters, "string")

	data, err := Decode(outputParameters, input)
	if err != nil {
		println(err)
	} else {
		fmt.Println(data[0])
	}
}
