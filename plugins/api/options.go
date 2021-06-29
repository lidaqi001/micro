package api

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Option func(opts *Options)

type Options struct {
	Context context.Context

	Router func(engine *gin.Engine)
}

type routeKey struct{}

func Route(val func(engine *gin.Engine)) Option {
	return SetOption(routeKey{}, val)
}
