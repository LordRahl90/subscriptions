package servers

import (
	"net/http"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/users"
	"subscriptions/domains/vouchers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	initError      error
	userService    users.IUserService
	productService products.Manager
	voucherService vouchers.Manager
	planService    plans.Manager
)

// Server container for the server object
type Server struct {
	Router        *gin.Engine
	DB            *gorm.DB
	SigningSecret string
}

// New returns a new instance of server with the provided db
// this also initializes all the necessary services.
func New(db *gorm.DB) (*Server, error) {
	userService, initError = users.New(db)
	if initError != nil {
		return nil, initError
	}

	productService, initError = products.New(db)
	if initError != nil {
		return nil, initError
	}

	voucherService, initError = vouchers.New(db)
	if initError != nil {
		return nil, initError
	}

	planService, initError = plans.New(db)
	if initError != nil {
		return nil, initError
	}

	s := &Server{DB: db}
	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	r := gin.New()
	r.POST("/login", s.authenticate)
	r.POST("/user/create", s.createUser)
	r.Use(s.authenticated())

	s.Router = r
	s.productsRoute()
	s.vouchersRoute()
	s.plansRoute()
}

func success(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, response)
}

func created(c *gin.Context, response interface{}) {
	c.JSON(http.StatusCreated, response)
}

func unAuthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": err.Error(),
		"data":    nil,
	})
}

func badRequestFromError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": err.Error(),
		"data":    nil,
	})
}

func internalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": err.Error(),
		"data":    nil,
	})
}
