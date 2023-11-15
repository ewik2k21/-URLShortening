package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}

func RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		ctx.Next()
		duration := time.Since(t)
		Log.Sugar().Infoln(
			"uri", ctx.Request.RequestURI,
			"method", ctx.Request.Method,
			"duration", duration,
		)
	}
}

func ResponseLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Next()

		Log.Sugar().Infoln(
			"statusCode", ctx.Writer.Status(),
			"content-size", ctx.Writer.Size(),
			"location", ctx.GetHeader("Location"),
		)
	}
}
