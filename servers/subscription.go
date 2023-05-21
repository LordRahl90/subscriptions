package servers

import (
	"errors"
	"fmt"
	"subscriptions/domains/core"
	"subscriptions/domains/subscription"
	"subscriptions/requests"
	"subscriptions/responses"
	"subscriptions/services/purchase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		TotalDuration: sub.Duration + sub.TrialPeriod,
		Price:         sub.Amount,
		Tax:           sub.Tax,
		Discount:      sub.Discount,
		Total:         sub.Total,
		Status:        sub.Status.String(),
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
		internalError(ctx, err)
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
			TotalDuration: res[i].Duration + res[i].TrialPeriod,
			Price:         res[i].Amount,
			Tax:           res[i].Tax,
			Discount:      res[i].Discount,
			Total:         res[i].Total,
			Status:        res[i].Status.String(),
		}
	}
	success(ctx, response)
}

func (s *Server) subscriptionDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}
	sub, err := s.subscriptionService.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx)
			return
		}
		internalError(ctx, err)
		return
	}
	if sub.UserID != userData.(*core.TokenData).UserID {
		// can be made available for admins
		unAuthorized(ctx, fmt.Errorf("you are not allowed to access this subscription details"))
		return
	}

	success(ctx, responses.Subscription{
		ID:            sub.ID,
		StartDate:     sub.CreatedAt,
		EndDate:       sub.EndDate(),
		Duration:      sub.Duration,
		TrialDuration: sub.TrialPeriod,
		TotalDuration: sub.Duration + sub.TrialPeriod,
		Price:         sub.Amount,
		Tax:           sub.Tax,
		Discount:      sub.Discount,
		Total:         sub.Total,
		Status:        sub.Status.String(),
	})
}

func (s *Server) pauseSubscription(ctx *gin.Context) {
	id := ctx.Param("id")
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}

	sub, err := s.subscriptionService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx)
			return
		}
		badRequestFromError(ctx, err)
		return
	}
	if sub.UserID != userData.(*core.TokenData).UserID {
		// can be made available for admins
		unAuthorized(ctx, fmt.Errorf("you are not allowed to access this subscription details"))
		return
	}

	if sub.IsTrial() {
		badRequestFromError(ctx, fmt.Errorf("you cannot pause during trial period"))
		return
	}

	if err := s.subscriptionService.UpdateStatus(ctx, id, subscription.StatusPaused); err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, "subscription updated successfully")
}

func (s *Server) unpauseSubscription(ctx *gin.Context) {
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}
	id := ctx.Param("id")

	sub, err := s.subscriptionService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx)
			return
		}
		badRequestFromError(ctx, err)
		return
	}
	if sub.UserID != userData.(*core.TokenData).UserID {
		// can be made available for admins
		unAuthorized(ctx, fmt.Errorf("you are not allowed to access this subscription details"))
		return
	}

	if err := s.subscriptionService.UpdateStatus(ctx, id, subscription.StatusActive); err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, "subscription cancelled successfully")
}

func (s *Server) cancelSubscription(ctx *gin.Context) {
	userData, ok := ctx.Get("userInfo")
	if !ok {
		unAuthorized(ctx, fmt.Errorf("please provide token"))
		return
	}
	id := ctx.Param("id")

	sub, err := s.subscriptionService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			notFound(ctx)
			return
		}
		badRequestFromError(ctx, err)
		return
	}
	if sub.UserID != userData.(*core.TokenData).UserID {
		// can be made available for admins
		unAuthorized(ctx, fmt.Errorf("you are not allowed to access this subscription details"))
		return
	}

	if err := s.subscriptionService.UpdateStatus(ctx, id, subscription.StatusCancelled); err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, "subscription cancelled successfully")
}
