package helper

import (
	"log"
	"reflect"
)

func IsNil(value interface{}) bool {

	if val := reflect.ValueOf(value); val.IsNil() {
		log.Println("sing：返回值为空")
		return true
	}
	return false
}

func Empty(ip string) bool {
	return len(ip) == 0
}