package servers

import (
	"subscriptions/domains/plans"
	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/gin-gonic/gin"
)

func (s *Server) plansRoute() {
	plans := s.Router.Group("plans")
	{
		plans.GET("/:id", planDetails)
		plans.POST("", createPlan)
	}
}

func createPlan(ctx *gin.Context) {
	var req requests.SubscriptionPlan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	plan := &plans.SubscriptionPlan{
		ProductID:     req.ProductID,
		Amount:        req.Amount,
		Duration:      req.Duration,
		TrialDuration: req.TrialDuration,
	}

	if err := planService.Create(ctx.Request.Context(), plan); err != nil {
		internalError(ctx, err)
		return
	}

	created(ctx, responses.SubscriptionPlan{
		ID:            plan.ID,
		ProductID:     plan.ProductID,
		Amount:        plan.Amount,
		Duration:      plan.Duration,
		TrialDuration: plan.TrialDuration,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	})

}

func planDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	plan, err := planService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		badRequestFromError(ctx, err)
	}

	success(ctx, responses.SubscriptionPlan{
		ID:            plan.ID,
		ProductID:     plan.ProductID,
		Amount:        plan.Amount,
		Duration:      plan.Duration,
		TrialDuration: plan.TrialDuration,
		CreatedAt:     plan.CreatedAt,
		UpdatedAt:     plan.UpdatedAt,
	})
}
