package helper

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

/*******************************************************
					助手函数
*******************************************************/

// The argument must be a chan, func, interface, map, pointer, or slice value;
// if it is not, IsNil panics.
func IsNil(value interface{}) bool {
	if val := reflect.ValueOf(value); val.IsNil() {
		return true
	}
	return false
}

// The len built-in function returns the length of v, according to its type:
//	Array: the number of elements in v.
//	Pointer to array: the number of elements in *v (even if v is nil).
//	Slice, or map: the number of elements in v; if v is nil, len(v) is zero.
//	String: the number of bytes in v.
//	Channel: the number of elements queued (unread) in the channel buffer;
//	         if v is nil, len(v) is zero.
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

// default : debug environment
func IsOpenDebug() (bool, error) {

	if debug := os.Getenv("DEBUG"); len(debug) > 0 {
		val, err := strconv.ParseInt(debug, 10, 64)
		if err != nil {
			return false, errors.New(fmt.Sprintf("isOpenDebug.error:%v", err))
		}

		// production environment
		if val == 0 {
			return false, nil
		}

	}

	// debug environment
	return true, nil
}
