package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

const (
	idParam = "id"

	Timeout5s = 5 * time.Second
)

func getIDFromPath(ctx *gin.Context) string {
	return ctx.Param(idParam)
}

func createTimeoutContext(ctx *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Request.Context(), Timeout5s)
}

func FormatList[V, T any](arr []V, fn func(V) T) (results []T) {
	for _, v := range arr {
		results = append(results, fn(v))
	}
	return
}
