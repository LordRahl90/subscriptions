package servers

import (
	"subscriptions/domains/vouchers"
	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/gin-gonic/gin"
)

func (s *Server) vouchersRoute() {
	vouchers := s.Router.Group("vouchers")
	{
		vouchers.GET("/:id/products", voucherProducts)
		vouchers.POST("/valid", checkValidity)
		vouchers.POST("", s.createVoucher)
	}
}

func (s *Server) createVoucher(ctx *gin.Context) {
	var req requests.Voucher

	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	voucher := &vouchers.Voucher{
		VoucherType: vouchers.VoucherTypeFromString(req.VoucherType),
		ProductID:   req.ProductID,
		Code:        req.Code,
		Percentage:  req.Percentage,
		Amount:      req.Amount,
	}

	if err := voucher.Validate(); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	if err := s.voucherService.Create(ctx.Request.Context(), voucher); err != nil {
		internalError(ctx, err)
		return
	}

	created(ctx, responses.Voucher{
		ID:          voucher.ID,
		VoucherType: voucher.VoucherType.String(),
		ProductID:   voucher.ProductID,
		Code:        voucher.Code,
		Active:      voucher.Active,
		Amount:      voucher.Amount,
		Percentage:  voucher.Percentage,
		CreatedAt:   voucher.CreatedAt,
		UpdatedAt:   voucher.UpdatedAt,
	})
}

func voucherProducts(ctx *gin.Context) {}
func checkValidity(ctx *gin.Context)   {}
