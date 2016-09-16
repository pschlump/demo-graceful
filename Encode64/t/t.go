package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	data := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	str, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("Error:%s\n", err)
	}
	fmt.Printf("%x\n", str)
}
