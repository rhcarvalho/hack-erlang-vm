package erl

import (
	"bytes"
	"fmt"
)

const oneByte = 8

func bytesToBitString(b []byte) string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "<<")
	for i, v := range b {
		if i > 0 {
			fmt.Fprint(&buf, ",")
		}
		fmt.Fprint(&buf, v)
	}
	fmt.Fprint(&buf, ">>")
	return buf.String()
}
