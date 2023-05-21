package servers

import "github.com/gin-gonic/gin"

func (s *Server) subscriptionRoute() {
	subs := s.Router.Group("subscriptions")
	{
		subs.POST("", purchase)
		subs.GET("me", userSubscriptions)
		subs.GET(":id", subscriptionDetail)
		subs.DELETE(":id", cancelSubscription)
		subs.PATCH(":id/pause", pauseSubscription)
		subs.PATCH(":id/unpause", unpauseSubscription)
	}
}

func purchase(ctx *gin.Context)            {}
func userSubscriptions(ctx *gin.Context)   {}
func subscriptionDetail(ctx *gin.Context)  {}
func pauseSubscription(ctx *gin.Context)   {}
func unpauseSubscription(ctx *gin.Context) {}
func cancelSubscription(ctx *gin.Context)  {}
