package servers

import (
	"fmt"
	"strings"

	"subscriptions/domains/core"

	"github.com/gin-gonic/gin"
)

var nonAuthPath = map[string]struct{}{
	"/products":           {},
	"/products/:id":       {},
	"/products/:id/plans": {},
	"/plans/:id":          {},
}

func (s *Server) authenticated() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		if _, ok := nonAuthPath[ctx.FullPath()]; ok {
			ctx.Next()
			return
		}
		headers := ctx.Request.Header
		authHeader, ok := headers["Authorization"]
		if !ok {
			unAuthorized(ctx, fmt.Errorf("authorization token not provided"))
			ctx.Abort()
			return
		}
		authData := strings.Split(authHeader[0], " ")
		if len(authData) != 2 {
			unAuthorized(ctx, fmt.Errorf("invalid token format provided"))
			ctx.Abort()
			return
		}
		authToken := authData[1]
		userInfo, err := core.Decode(authToken)
		if err != nil || userInfo == nil {
			unAuthorized(ctx, err)
			ctx.Abort()
			return
		}
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}
