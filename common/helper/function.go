package helper

import (
	"os"
	"reflect"
)

func IsNil(value interface{}) bool {

	if val := reflect.ValueOf(value); val.IsNil() {
		//log.Println("sing：返回值为空")
		return true
	}
	return false
}

func Empty(ip string) bool {
	return len(ip) == 0
}

//CreateDir  文件夹创建
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Chmod(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
