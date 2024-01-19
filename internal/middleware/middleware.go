package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
	// TODO: implement method to authenticate & authorization
	fmt.Println("auth")
	ctx.Next()
}
