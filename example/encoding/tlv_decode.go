package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ClearGrass/sdk-go/encoding"
)

func exampleDecodeData() {
	encodeData := "Q0c0QgA4AgA8ABEFADIuMC4zIgQAMDAwMCwBAAE0BQAxLjMuNzUFADIuMC4zHQEAARQTAMXA82VJ4i8EHAAAADICCQEAzwCWCw=="

	dataBytes, err := base64.StdEncoding.DecodeString(encodeData)
	if err != nil {
		return
	}

	data, err := encoding.TlvDecode(dataBytes)
	c, _ := json.Marshal(data)
	fmt.Println(string(c))
}

func exampleSetIntervalCmd() {
	msg := encoding.NewMessagePod(true)

	msg.CmdType = 0x32
	msg.IntervalSetting = &encoding.IntervalSetting{CollectInterval: 60, ReportInterval: 60}
	msg.SetDebug(1)

	ret, err := encoding.TlvEncode(msg)
	if err != nil {

	}

	fmt.Println(hex.EncodeToString(ret))
	fmt.Println(base64.StdEncoding.EncodeToString(ret))
}

func main() {
	//exampleDecodeData()
	exampleSetIntervalCmd()
}
