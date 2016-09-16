package Encode64

import (
	"encoding/base64"
	"fmt"
)

func Base64Encode(s string) (rv string) {
	rv = base64.StdEncoding.EncodeToString([]byte(s))
	return
}

func Base64Decode(s string) (rv string) {
	t, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Printf("Error decoding string=[%s]: err=%s\n", s, err)
	}
	rv = string(t)
	return
}
