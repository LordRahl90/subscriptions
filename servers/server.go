package servers

import (
	"net/http"

	"subscriptions/domains/plans"
	"subscriptions/domains/products"
	"subscriptions/domains/subscription"
	"subscriptions/domains/users"
	"subscriptions/domains/vouchers"
	"subscriptions/services/purchase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	purchaseService *purchase.Service
)

// Server container for the server object
type Server struct {
	Router        *gin.Engine
	DB            *gorm.DB
	SigningSecret string

	userService         users.IUserService
	productService      products.Manager
	voucherService      vouchers.Manager
	planService         plans.Manager
	subscriptionService subscription.Manager
}

// New returns a new server instance based on the injected services
func New(
	us users.IUserService,
	ps products.Manager,
	vs vouchers.Manager,
	pls plans.Manager,
	sbs subscription.Manager,
) *Server {
	purchaseService = purchase.New(ps, pls, vs, sbs)
	s := &Server{
		userService:         us,
		productService:      ps,
		planService:         pls,
		subscriptionService: sbs,
	}
	s.setupRoutes()
	return s
}

// NewWithDefaults returns a new instance of server with the defaults initialized with provided db
// this also initializes all the necessary services.
func NewWithDefaults(db *gorm.DB) (*Server, error) {
	userService, err := users.New(db)
	if err != nil {
		return nil, err
	}

	productService, err := products.New(db)
	if err != nil {
		return nil, err
	}

	voucherService, err := vouchers.New(db)
	if err != nil {
		return nil, err
	}

	planService, err := plans.New(db)
	if err != nil {
		return nil, err
	}

	subscriptionService, err := subscription.New(db)
	if err != nil {
		return nil, err
	}

	purchaseService = purchase.New(productService, planService, voucherService, subscriptionService)
	s := &Server{
		DB:                  db,
		userService:         userService,
		productService:      productService,
		voucherService:      voucherService,
		planService:         planService,
		subscriptionService: subscriptionService,
	}
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
	s.subscriptionRoute()
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
