package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"stp_dao_v2/models"
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
	var addr []string
	var amt []string
	addr = append(addr, "0x18041866663b077bB6BF2bAFFAeA2451a2472eD7", "0x5718D9C95D15a766E9DdE6579D7B93Eaa88a26b8")
	amt = append(amt, "1000000000000000000", "1000000000000000000")
	data := models.AirdropAddressArray{
		Address: addr,
		Amount:  amt,
	}
	fmt.Println(data)
	fmt.Println(Keccak256("1"))
	fmt.Println(fmt.Sprintf("%x", common.LeftPadBytes([]byte("params.Handle"), 32)))
	var f float64
	f = float64(1) / float64(3)
	fmt.Println(f)
	uuid := GenerateUuid()
	fmt.Println(uuid)

	a, _ := Hex2Dec("4000000000000000000000000000000000000000000000000000000000000000")
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
