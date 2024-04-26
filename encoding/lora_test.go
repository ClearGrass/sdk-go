package encoding

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeLoraRobbData(t *testing.T) {
	hexDataList := []string{
		"01412401662758412FF2F5000003F6000D000E005C000C00000146FF000C3ABDCD00312E302E36EA05",
		"01411B0066276E0403842FF2F5000003F6000D000E005C000C00000146FF9076",
	}

	base64DataList := []string{
		"AUEkAWYq9Tgu8isAAAILAAYABgBsACYAAAG8/WA0MZYWADEuMS4weNA=",
	}

	for _, src := range hexDataList {

		bs, _ := hex.DecodeString(src)
		fmt.Println(hex.EncodeToString(bs))

		out, _ := DecodeLoraRobbData(bs)
		c, _ := json.Marshal(out)
		fmt.Println(string(c))
	}

	for _, src := range base64DataList {
		bs, _ := base64.StdEncoding.DecodeString(src)
		fmt.Println(hex.EncodeToString(bs))

		out, _ := DecodeLoraRobbData(bs)
		c, _ := json.Marshal(out)
		fmt.Println(string(c))
	}
}
