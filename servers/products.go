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
		product.POST("", createProduct)
		product.GET("", allProducts)
		product.GET("/:id", singleProduct)
		product.GET("/:id/plans", productPlans)
		product.POST("/:id", purchase)
	}
}

func createProduct(ctx *gin.Context) {
	var req *requests.Product

	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	product := &products.Product{
		Name:        req.Name,
		Description: req.Description,
		Tax:         req.Tax,
		TrialExists: req.TrialExists,
	}

	if err := productService.Create(ctx, product); err != nil {
		internalError(ctx, err)
		return
	}

	created(ctx, responses.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Tax:         product.Tax,
		TrialExists: product.TrialExists,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}

func purchase(ctx *gin.Context) {}

func allProducts(ctx *gin.Context) {
	res, err := productService.Find(ctx.Request.Context())
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
			Tax:         res[i].Tax,
			TrialExists: res[i].TrialExists,
			CreatedAt:   res[i].CreatedAt,
			UpdatedAt:   res[i].UpdatedAt,
		}
	}

	success(ctx, result)
}

func singleProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := productService.FindOne(ctx.Request.Context(), id)
	if err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, res)
}

func productPlans(ctx *gin.Context) {
	productID := ctx.Param("id")
	res, err := planService.Find(ctx.Request.Context(), productID)
	if err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, res)
}
