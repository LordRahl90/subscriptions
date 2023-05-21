package servers

import (
	"fmt"
	"subscriptions/domains/core"
	"subscriptions/requests"
	"subscriptions/responses"
	"subscriptions/services/purchase"

	"github.com/gin-gonic/gin"
)

func (s *Server) subscriptionRoute() {
	subs := s.Router.Group("subscriptions")
	{
		subs.POST("", newPurchase)
		subs.GET("", s.userSubscriptions)
		subs.GET(":id", s.subscriptionDetail)
		subs.DELETE(":id", s.cancelSubscription)
		subs.PATCH(":id/pause", s.pauseSubscription)
		subs.PATCH(":id/unpause", s.unpauseSubscription)
	}
}

func newPurchase(ctx *gin.Context) {
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}
	var req requests.Purchase
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	pch := &purchase.Purchase{
		UserID:    userData.(*core.TokenData).UserID,
		ProductID: req.ProductID,
		PlanID:    req.PlanID,
		Voucher:   req.Voucher,
	}
	sub, err := purchaseService.Process(ctx.Request.Context(), pch)
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}
	created(ctx, responses.Subscription{
		ID:            sub.ID,
		StartDate:     sub.CreatedAt,
		EndDate:       sub.EndDate(),
		Duration:      sub.Duration,
		TrialDuration: sub.TrialPeriod,
		Price:         sub.Amount,
		Tax:           sub.Tax,
		Discount:      sub.Discount,
		Total:         sub.Total,
	})
}
func (s *Server) userSubscriptions(ctx *gin.Context) {
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}
	res, err := s.subscriptionService.Find(ctx.Request.Context(), userData.(*core.TokenData).UserID)
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}

	response := make([]responses.Subscription, len(res))
	for i := range res {
		response[i] = responses.Subscription{
			ID:            res[i].ID,
			StartDate:     res[i].CreatedAt,
			EndDate:       res[i].EndDate(),
			Duration:      res[i].Duration,
			TrialDuration: res[i].TrialPeriod,
			Price:         res[i].Amount,
			Tax:           res[i].Tax,
			Discount:      res[i].Discount,
			Total:         res[i].Total,
		}
	}
	success(ctx, response)
}

func (s *Server) subscriptionDetail(ctx *gin.Context)  {}
func (s *Server) pauseSubscription(ctx *gin.Context)   {}
func (s *Server) unpauseSubscription(ctx *gin.Context) {}
func (s *Server) cancelSubscription(ctx *gin.Context)  {}
