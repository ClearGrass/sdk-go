package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ClearGrass/sdk-go/encoding"
)

func main() {
	encodeData := "01412401662758412FF2F5000003F6000D000E005C000C00000146FF000C3ABDCD00312E302E36EA05"

	dataBytes, err := hex.DecodeString(encodeData)
	if err != nil {
		return
	}

	data, err := encoding.DecodeLoraRobbData(dataBytes)
	c, _ := json.Marshal(data)
	fmt.Println(string(c))
}
