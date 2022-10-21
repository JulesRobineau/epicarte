package middleware

import (
	"fmt"
	"gin-template/utils/jwt"
	"github.com/gin-gonic/gin"
	"time"
)

const NoId = "Ø"

func CustomLogger(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	user := NoId
	if claimsStr, ok := param.Keys["claims"]; ok {
		claims := claimsStr.(*jwt.Claims)
		user = fmt.Sprintf("%d", claims.UserId)
	}

	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %20s | %15s |%s %-7s %s %#v\n%s",
		param.TimeStamp.Format("01/02/2006 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		user,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
