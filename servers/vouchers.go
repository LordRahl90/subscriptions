package servers

import "github.com/gin-gonic/gin"

func (s *Server) vouchersRoute() {
	vouchers := s.Router.Group("vouchers")
	{
		vouchers.GET("/:id/products", voucherProducts)
		vouchers.POST("/valid", checkValidity)
		vouchers.POST("", createVoucher)
	}
}

func createVoucher(ctx *gin.Context) {

}

func voucherProducts(ctx *gin.Context) {}
func checkValidity(ctx *gin.Context)   {}
