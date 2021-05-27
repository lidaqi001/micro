package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"github.com/lidaqi001/micro/examples/api/jwt"
)

const (
	TOKEN_FMT_INVALID = "Authorization Format invalid !"
	TOKEN_NOT_EXIST   = "Token not exist !"
	TOKEN_INVALID     = "Token is invalid !"
	REFUSE_CODE       = 403
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("进入Auth中间件")

		// 请求前

		token := c.Request.Header.Get("Authorization")

		var stringToken string
		switch {
		case len(token) > 0:
			// 从header中取token

			splitToken := strings.Split(token, "Bearer ")
			log.Printf("%v,%v", len(splitToken), splitToken)
			// header参数格式错误终止
			if len(splitToken) != 2 {
				refuse_response(c, TOKEN_FMT_INVALID)
				return
			}
			stringToken = splitToken[1]

		case len(token) == 0:
			// header中不存在，从url中取

			if t, ok := c.GetQuery("token"); ok {
				stringToken = t
				break
			}
			// 都不存在，则终止
			refuse_response(c, TOKEN_NOT_EXIST)
			return
		}

		// 验证token
		if ok, _ := jwt.VerifyMiddleware(stringToken); !ok {
			refuse_response(c, TOKEN_INVALID)
			return
		}

		c.Set("token", stringToken)

		c.Next()

		// 请求后

		log.Println("离开Auth中间件")
	}
}

func refuse_response(c *gin.Context, msg string) {

	code := REFUSE_CODE
	c.Abort()
	c.JSON(code, gin.H{
		"message": msg,
	})
}
