package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ClearGrass/sdk-go/encoding"
)

func main() {
	encodeData := "AUEkAWYq9Tgw8isAAAIlAAYABgBsACYAAAG8/wA0MzYwADEuMS4weNA="

	dataBytes, err := base64.StdEncoding.DecodeString(encodeData)
	if err != nil {
		return
	}

	data, err := encoding.DecodeLoraRobbData(dataBytes)
	c, _ := json.Marshal(data)
	fmt.Println(string(c))
}
