package client

import (
	"errors"
	"github.com/lidaqi001/micro/plugins/logger"
)

func err(msg string) error {
	e := errors.New(msg)
	logger.Error(e)
	return e
}

const NAME_IS_NULL = "create service: The ServiceName cannot be empty!"

const CALL_FUNC_IS_NULL = "create service: The CallFunc cannot be empty~"
