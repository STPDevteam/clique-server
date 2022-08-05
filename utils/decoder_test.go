package utils

import (
	"testing"
)

func TestDecode(t *testing.T) {

	var input = ""
	var outputParameters []string
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")
	outputParameters = append(outputParameters, "string")

	data, err := Decode(outputParameters, input)
	if err != nil {
		println(err)
	} else {
		println(data)
	}
}
