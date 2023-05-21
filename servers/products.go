package servers

import (
	"subscriptions/domains/products"
	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/gin-gonic/gin"
)

func (s *Server) productsRoute() {
	product := s.Router.Group("products")
	{
		product.POST("", s.createProduct)
		product.GET("", s.allProducts)
		product.GET(":id", s.singleProduct)
		product.GET(":id/plans", s.productPlans)
	}
}

func (s *Server) createProduct(ctx *gin.Context) {
	var req *requests.Product

	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	product := &products.Product{
		Name:        req.Name,
		Description: req.Description,
		TaxRate:     req.Tax,
	}

	if err := s.productService.Create(ctx, product); err != nil {
		internalError(ctx, err)
		return
	}

	created(ctx, responses.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Tax:         product.TaxRate,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}

func (s *Server) allProducts(ctx *gin.Context) {
	res, err := s.productService.Find(ctx.Request.Context())
	if err != nil {
		internalError(ctx, err)
		return
	}
	result := make([]responses.Product, len(res))

	for i := range res {
		result[i] = responses.Product{
			ID:          res[i].ID,
			Name:        res[i].Name,
			Description: res[i].Description,
			Tax:         res[i].TaxRate,
			CreatedAt:   res[i].CreatedAt,
			UpdatedAt:   res[i].UpdatedAt,
		}
	}

	success(ctx, result)
}

func (s *Server) singleProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := s.productService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, responses.Product{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Tax:         res.TaxRate,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	})
}

func (s *Server) productPlans(ctx *gin.Context) {
	productID := ctx.Param("id")
	res, err := s.planService.Find(ctx.Request.Context(), productID)
	if err != nil {
		internalError(ctx, err)
		return
	}
	result := make([]responses.SubscriptionPlan, len(res))

	for i := range res {
		result[i] = responses.SubscriptionPlan{
			ID:            res[i].ID,
			ProductID:     res[i].ProductID,
			Amount:        res[i].Amount,
			Duration:      res[i].Duration,
			TrialDuration: res[i].TrialDuration,
			CreatedAt:     res[i].CreatedAt,
			UpdatedAt:     res[i].UpdatedAt,
		}
	}

	success(ctx, result)
}
