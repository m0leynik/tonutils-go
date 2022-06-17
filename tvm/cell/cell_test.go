package cell

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestCell_HashSign(t *testing.T) {
	cc1 := BeginCell().MustStoreUInt(111, 63).EndCell()
	cc2 := BeginCell().MustStoreUInt(772227, 63).MustStoreRef(cc1).EndCell()
	cc3 := BeginCell().MustStoreUInt(333, 63).MustStoreRef(cc2).EndCell()
	cc := BeginCell().MustStoreUInt(777, 63).MustStoreRef(cc3).EndCell()

	b, _ := hex.DecodeString("bb2509fe3cff8f1faae19213774d218c018f9616cd397850c8ad9038db84eaa9")

	if !bytes.Equal(cc.Hash(), b) {
		t.Log(hex.EncodeToString(cc.Hash()))
		t.Log(hex.EncodeToString(b))
		t.Fatal("hash diff")
	}

	pub, priv, _ := ed25519.GenerateKey(nil)
	if !ed25519.Verify(pub, cc.Hash(), cc.Sign(priv)) {
		t.Fatal("sign not match")
	}
}

func TestBOC(t *testing.T) {
	str := "b5ee9c7201021b010003b2000271c000ab558f4db84fd31f61a273535c670c091ffc619b1cdbbe5769a0bf28d3b8fea236865b4312ab35600000625f2d741f0d6773533c74d34001020114ff00f4a413f4bcf2c80b0301510000002629a9a317c878acda0aa0cfacdab9bff8bca840e7d10d8a41d1ee96caf7ac645016af94dfc0160201200405020148060704f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff1213141502e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d08090201200a0b007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e20201200c0d0059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f31840201580e0f0011b8c97ed44d0d70b1f8003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c6002012010110019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc0006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed5402057fc01817180042bf8e1b0bc5dfcda03e92f9b4b9ffc438595770c0686d91bde674ad610dba9bc66e020148191a0041bf0f895e56f2933fdc5f7c21bc29292fdf0415b7368b9a3eef5bd23ced3021278a0041bf16fc68f92304fb493ca52b5ddefabc42a2131f3e45442b1f2ae45156b2972bea"
	data, _ := hex.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}

	boc := c.ToBOCWithFlags(false)

	if str != hex.EncodeToString(boc) {
		t.Log(str)
		t.Log(hex.EncodeToString(boc))
		t.Fatal("boc not same")
	}
}

func TestSmallBOC(t *testing.T) {
	str := "b5ee9c72010101010002000000"

	c := BeginCell().EndCell()

	boc := c.ToBOCWithFlags(false)

	if str != hex.EncodeToString(boc) {
		t.Log(str)
		t.Log(hex.EncodeToString(boc))
		t.Fatal("boc not same")
	}
}

func TestBOCWithCRC(t *testing.T) {
	str := "b5ee9c7241021b010003b2000271c000ab558f4db84fd31f61a273535c670c091ffc619b1cdbbe5769a0bf28d3b8fea236865b4312ab35600000625f2d741f0d6773533c74d34001020114ff00f4a413f4bcf2c80b0301510000002629a9a317c878acda0aa0cfacdab9bff8bca840e7d10d8a41d1ee96caf7ac645016af94dfc0160201200405020148060704f8f28308d71820d31fd31fd31f02f823bbf264ed44d0d31fd31fd3fff404d15143baf2a15151baf2a205f901541064f910f2a3f80024a4c8cb1f5240cb1f5230cbff5210f400c9ed54f80f01d30721c0009f6c519320d74a96d307d402fb00e830e021c001e30021c002e30001c0039130e30d03a4c8cb1f12cb1fcbff1213141502e6d001d0d3032171b0925f04e022d749c120925f04e002d31f218210706c7567bd22821064737472bdb0925f05e003fa403020fa4401c8ca07cbffc9d0ed44d0810140d721f404305c810108f40a6fa131b3925f07e005d33fc8258210706c7567ba923830e30d03821064737472ba925f06e30d08090201200a0b007801fa00f40430f8276f2230500aa121bef2e0508210706c7567831eb17080185004cb0526cf1658fa0219f400cb6917cb1f5260cb3f20c98040fb0006008a5004810108f45930ed44d0810140d720c801cf16f400c9ed540172b08e23821064737472831eb17080185005cb055003cf1623fa0213cb6acb1fcb3fc98040fb00925f03e20201200c0d0059bd242b6f6a2684080a06b90fa0218470d4080847a4937d29910ce6903e9ff9837812801b7810148987159f31840201580e0f0011b8c97ed44d0d70b1f8003db29dfb513420405035c87d010c00b23281f2fff274006040423d029be84c6002012010110019adce76a26840206b90eb85ffc00019af1df6a26840106b90eb858fc0006ed207fa00d4d422f90005c8ca0715cbffc9d077748018c8cb05cb0222cf165005fa0214cb6b12ccccc973fb00c84014810108f451f2a7020070810108d718fa00d33fc8542047810108f451f2a782106e6f746570748018c8cb05cb025006cf165004fa0214cb6a12cb1fcb3fc973fb0002006c810108d718fa00d33f305224810108f459f2a782106473747270748018c8cb05cb025005cf165003fa0213cb6acb1f12cb3fc973fb00000af400c9ed5402057fc01817180042bf8e1b0bc5dfcda03e92f9b4b9ffc438595770c0686d91bde674ad610dba9bc66e020148191a0041bf0f895e56f2933fdc5f7c21bc29292fdf0415b7368b9a3eef5bd23ced3021278a0041bf16fc68f92304fb493ca52b5ddefabc42a2131f3e45442b1f2ae45156b2972bea32690605"
	data, _ := hex.DecodeString(str)

	c, err := FromBOC(data)
	if err != nil {
		t.Fatal(err)
	}

	boc := c.ToBOC()

	if str != hex.EncodeToString(boc) {
		t.Fatal("boc not same")
	}
}

func TestCell_Hash1(t *testing.T) {
	emptyHash, _ := new(big.Int).SetString("68134197439415885698044414435951397869210496020759160419881882418413283430343", 10)

	if !bytes.Equal(BeginCell().EndCell().Hash(), emptyHash.Bytes()) {
		t.Fatal("empty cell hash incorrect")
		return
	}

	refRef57bitsHash, _ := new(big.Int).SetString("111217512120054409408353636830563617100690868120902564345366655075146083288188", 10)

	if !bytes.Equal(BeginCell().MustStoreUInt(7, 5).MustStoreRef(
		BeginCell().MustStoreRef(
			BeginCell().MustStoreUInt(777777888, 57).EndCell(),
		).EndCell(),
	).EndCell().Hash(), refRef57bitsHash.Bytes()) {
		t.Fatal("refRef57bits cell hash incorrect")
		return
	}
}
