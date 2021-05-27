package handler

import "github.com/gin-gonic/gin"

type handler struct {
}

type Handler interface {
	WelCome() gin.HandlerFunc
	Client() gin.HandlerFunc
	ClientAsyncA() gin.HandlerFunc
	ClientAsyncB() gin.HandlerFunc
}

func NewHandler() Handler {
	return &handler{}
}
