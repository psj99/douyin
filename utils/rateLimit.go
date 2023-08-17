package utils

import (
	"douyin/service/types/response"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// gin中间件
// 限流(起始/最大请求处理量, 每秒恢复量)
func MiddlewareRateLimit(capacity int64, recover int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(time.Second, capacity, recover)
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			ZapLogger.Warnf("MiddlewareRateLimit warn: 达到处理量上限")
			ctx.JSON(http.StatusTooManyRequests, &response.Status{
				Status_Code: -1,
				Status_Msg:  "请求过频",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
