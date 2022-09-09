package utils

import (
	"fmt"
	"testing"
)

func TestPersonalSign(t *testing.T) {
	var (
		message = "Welcome come Clique"
		pri_key = "57437f659873be891d0adf7515ba7bd95e3cdb0d2285ef2e7ba242129285edf0"
	)
	address, signature, err := PersonalSign(pri_key, message)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("address = %s", address)
	t.Logf("signature = %s", signature)
}

func TestCheckPersonalSign(t *testing.T) {
	var (
		message   = "Welcome come Clique"
		signature = "d3ce2bad1668165dd2f410a645357f8338eb656448850049105f4162ca219870698287e50ed2405621c5a247b96ba699dae70eb5ae9743d6e9b013b56591632301"
		address   = "0x0c99F596a56872810d8B4139564ca7F02CCe8d45"
	)
	verify, err := CheckPersonalSign(message, address, signature)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("verify = %v", verify)
}

func TestSignMessage(t *testing.T) {
	var f float64
	f = float64(1) / float64(3)
	fmt.Println(f)
	uuid := GenerateUuid()
	fmt.Println(uuid)

	a := Hex2Dec("")
	fmt.Println(a)
	to, err := Hex2BigInt("0x0")
	str := to.String()
	fmt.Println(str, to, err)

	var (
		message = "dDCb698c476e058c59BC99eD041aaB02680aA5c5" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000001" +
			"dAC17F958Da2206202ee5236994597C13D831ec7" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000"
		pri_key = "493fd480a29d55282e130778a726121993340a244e6ea18a5aa6de9cf6248296"
	)
	signature, _ := SignMessage(message, pri_key)

	t.Logf("signature = %s", signature)
}
