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
		plans.GET("/:id", s.planDetails)
		plans.POST("", s.createPlan)
	}
}

func (s *Server) createPlan(ctx *gin.Context) {
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

	if err := s.planService.Create(ctx.Request.Context(), plan); err != nil {
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

func (s *Server) planDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	plan, err := s.planService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		badRequestFromError(ctx, err)
		return
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
