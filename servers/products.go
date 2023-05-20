package servers

import "github.com/gin-gonic/gin"

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

func createProduct(ctx *gin.Context) {}
func purchase(ctx *gin.Context)      {}

func allProducts(ctx *gin.Context) {
	res, err := productService.Find(ctx.Request.Context())
	if err != nil {
		internalError(ctx, err)
		return
	}

	success(ctx, res)
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
