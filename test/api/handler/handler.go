package handler

import "github.com/gin-gonic/gin"

type handler struct {
}

type Handler interface {
	WelCome() gin.HandlerFunc
	Client1() gin.HandlerFunc
}

func NewHandler() Handler {
	return &handler{}
}
