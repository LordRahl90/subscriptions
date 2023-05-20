package servers

import "github.com/gin-gonic/gin"

func (s *Server) plansRoute() {
	plans := s.Router.Group("plans")
	{
		plans.GET("/:id", planDetails)
		plans.POST("/", createPlan)
	}
}

func createPlan(ctx *gin.Context)  {}
func planDetails(ctx *gin.Context) {}
