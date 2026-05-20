package encoding

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeBleData(t *testing.T) {
	base64DataList := []string{
		"0201061416CDFD884FCB3E85342D58010410010B02020164",
		"0201061416CDFDC812D3E560342D580804001100000F01FE",
		"0201061416CDFDC812D3E560342D580804011700000F018F",
	}

	for _, src := range base64DataList {
		out, err := DecodeBleData(src)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		c, _ := json.Marshal(out)
		fmt.Println(string(c))
	}
}
