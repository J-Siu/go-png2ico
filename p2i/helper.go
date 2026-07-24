package p2i

import (
	"encoding/json"
	"fmt"
)

func PrintStruct(prefix string, s any) {
	if prefix != "" {
		fmt.Println(prefix)
	}
	if s != nil {
		if b, err := json.MarshalIndent(s, "", "  "); err == nil {
			fmt.Println(string(b))
		}
	}
}
