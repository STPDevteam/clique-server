package utils

import (
	"fmt"
	"testing"
)

func TestPersonalSign(t *testing.T) {
	//a := 0b0
	//a = a | (1 << 0)
	//a = a | (1 << 1)
	//a = a | (1 << 2)
	//a = a | (1 << 3)
	//fmt.Printf("%b\n", a)
	//a = a & ^(1 << 0)
	//a = a & ^(1 << 1)
	//a = a & ^(1 << 2)
	//a = a & ^(1 << 3)
	//fmt.Printf("%b\n", a)

	b := 0b1110
	ba := (b & (1 << 0)) > 0
	println(ba)
	fmt.Printf("%b\n", b)
}

func TestAi(t *testing.T) {
	openai_bearer_key := "Bearer "
	url := "https://api.openai.com/v1/chat/completions"
	chat, err := AiChat("0xdDCb698c459BC99eD0476e058c1aaB02680aA5c5", "Use golang to write a method to take the minimum value of two numbers", url, openai_bearer_key)
	if err != nil {
		println(err.Error())
		return
	}
	println(chat)
}
