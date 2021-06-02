package service

import (
	"context"
)

type S struct {
	ctx context.Context
}

type SI interface {
	Service()
	Client()
}

func NewService(ctx context.Context) SI {
	if ctx == nil {
		ctx = context.Background()
	}
	return &S{ctx: ctx}
}

func (s *S) Service() {
	panic("implement me")
}

func (s *S) Client() {
	panic("implement me")
}
