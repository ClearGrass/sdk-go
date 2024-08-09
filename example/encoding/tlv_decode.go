package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ClearGrass/sdk-go/encoding"
)

func main() {
	encodeData := "Q0c0QgA4AgA8ABEFADIuMC4zIgQAMDAwMCwBAAE0BQAxLjMuNzUFADIuMC4zHQEAARQTAMXA82VJ4i8EHAAAADICCQEAzwCWCw=="

	dataBytes, err := base64.StdEncoding.DecodeString(encodeData)
	if err != nil {
		return
	}

	data, err := encoding.TlvDecode(dataBytes)
	c, _ := json.Marshal(data)
	fmt.Println(string(c))
}
