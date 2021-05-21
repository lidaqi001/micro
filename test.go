package main

import (
	"log"
	"strings"
)

func main() {
	res := strings.Split("127.0.0.1:1111", ":")
	log.Println(res[0])
}
