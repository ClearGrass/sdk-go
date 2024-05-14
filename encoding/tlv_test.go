package encoding

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestTlvEncode_ReadingOffset(t *testing.T) {
	msg := NewMessagePod(true)

	msg.CmdType = 50
	msg.ReadingOffsetSetting = &ReadingOffsetSetting{
		Temperature: &ReadingOffset{OffsetValue: 2},
	}

	ret, err := TlvEncode(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(ret))
}
