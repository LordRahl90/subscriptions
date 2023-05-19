package servers

import (
	"net/http"

	"subscriptions/domains/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	userService users.IUserService
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
	us, err := users.New(db)
	if err != nil {
		return nil, err
	}
	userService = us

	s := &Server{DB: db}
	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	r := gin.Default()
	r.POST("/login", s.authenticate)
	r.POST("/user/create", s.createUser)
	r.Use(s.authenticated())

	s.Router = r
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
