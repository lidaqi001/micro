package service

import (
	"errors"
	"github.com/lidaqi001/micro/plugins/logger"
)

func sErr(v interface{}) error {
	var e error

	switch v.(type) {
	case string:
		e = errors.New(v.(string))
	case error:
		e = v.(error)
	}

	logger.Error("service error: ", e)

	return e
}

const SERVICE_NAME_IS_NULL = "create service: The ServiceName cannot be empty!"

const CALL_FUNC_IS_NULL = "create service: The CallFunc cannot be empty~"
