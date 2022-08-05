package utils

import (
	"fmt"
	"testing"
)

func TestPersonalSign(t *testing.T) {
	var (
		message = "Welcome come Clique"
		pri_key = "285e587fd9421292dd2f2e770adfedf05e651515ba7b3cdb5ba743098793be82"
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
		message = "Welcome come Clique"
		// signature = "68687e64717511a9f7d4f79469f781d8b42f056ce2aed62393532eb85f7b51267dafe19038c4d44dbaf63469a81f6835ad7f0dbcfc75f6354efab65ea6e9b0911c"
		signature = "d3ce2bad166505621c5a24b01816357f83385dd27b96ba699dae70eb5ae9743d6e9f410a645eb6564488500f4162ca219870698287e50ed2449103b56591632301"
		address   = "0x0c910d8B402CCe8d49F596a568728139564ca7F5"
	)
	verify, err := CheckPersonalSign(message, address, signature)
	if nil != err {
		t.Fatal(err)
	}
	t.Logf("verify = %v", verify)
}

/**
keccak256(
	abi.encodePacked(
		msg.sender, 40
		_nonce, 64
		signInfo_.chainId,
		signInfo_.tokenAddress,
		signInfo_.balance,
		uint256(signInfo_.signType)
	)
)

0000000000000000000000000000000000000000000000000000000000000000

ddcb698c459bc99ed0476e058c1aab02680aa5c5
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000001
dac17f958d2ee523a2206206994597c13d831ec7
0000000000000000000000000000000000000000000000000000000000000000
0000000000000000000000000000000000000000000000000000000000000000

56fdf65d7349647082676f281f0493f3a93d249ba2de7ca89c28b28a847e95155f64de3c9194c1c281c28ed829fc52c2137dbd960a83368203c76307800f3d311c
56fdf65d7349647082676f281f0493f3a93d249ba2de7ca89c28b28a847e95155f64de3c9194c1c281c28ed829fc52c2137dbd960a83368203c76307800f3d311c


0x
0000000000000000000000000000000000000000000000000000000000000100
0000000000000000000000000000000000000000000000000000000000000140
0000000000000000000000000000000000000000000000000000000000000180
00000000000000000000000000000000000000000000000000000000000001e0
00000000000000000000000000000000000000000000000000000000000003e0
0000000000000000000000000000000000000000000000000000000000000400
0000000000000000000000000000000000000000000000000000000000000420
0000000000000000000000000000000000000000000000000000000000000440

000000000000000000000000000000000000000000000000000000000000001a
4142434445464748494a4b4c4d4e4f505152535455565758595a000000000000
000000000000000000000000000000000000000000000000000000000000001a
4142434445464748494a4b4c4d4e4f505152535455565758595a000000000000
000000000000000000000000000000000000000000000000000000000000002f
536f6369616c2c50726f746f636f6c2c4e46542c4d65746176657273652c4761
6d696e672c446170702c4f746865720000000000000000000000000000000000
00000000000000000000000000000000000000000000000000000000000001d4
4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a4142434445464748494a4b4c4d4e4f505152535455565758595a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000012268747470733a2f2f67696d67332e62616964752e636f6d2f7365617263682f7372633d6874747025334125324625324670696373332e62616964752e636f6d25324666656564253246636165663736303934623336616361666537333433653434316161663237316130306539396333372e6a706567253346746f6b656e25334433363539336332346534613964326431333934353061363264613839386136392672656665723d687474702533412532462532467777772e62616964752e636f6d266170703d323032312673697a653d663336302c323430266e3d3026673d306e26713d373526666d743d6175746f3f7365633d3136353834323238303026743d6136653736653933653232636266396161613135343935333931326633633766000000000000000000000000000000000000000000000000000000000000

*/
func TestSignMessage(t *testing.T) {
	a := Hex2Dec("")
	fmt.Println(a)
	to, err := Hex2BigInt("0x0")
	str := to.String()
	fmt.Println(str, to, err)

	var (
		message = "dDCb698c459BC99eD0476e058c1aaB02680aA5c5" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000001" +
			"dAC17F958D2ee523a2206206994597C13D831ec7" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000"
		pri_key = "493fd480a29d55282e121993340a244e6ea1830778a7261a5aa6de9cf6248296"
	)
	signature, _ := SignMessage(message, pri_key)

	t.Logf("signature = %s", signature)
}
