package encoding

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeBleData(t *testing.T) {
	base64DataList := []string{
		//"014709003C03840000000000000000000000000000000000000503E805780000010100C32C",
		"0201061416CDFD884FCB3E85342D58010410010B02020164",
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
