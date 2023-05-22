package servers

import (
	"fmt"

	"subscriptions/domains/core"
	"subscriptions/domains/users"
	"subscriptions/requests"
	"subscriptions/responses"

	"github.com/gin-gonic/gin"
)

func (s *Server) authenticate(ctx *gin.Context) {
	var req requests.Authenticate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	res, err := s.userService.FindByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		badRequestFromError(ctx, fmt.Errorf("email not registered"))
		return
	}
	if res == nil {
		badRequestFromError(ctx, fmt.Errorf("email not registered"))
		return
	}

	if !core.CheckPasswordHash(req.Password, res.Password) {
		badRequestFromError(ctx, fmt.Errorf("password does not match"))
		return
	}
	td := core.TokenData{
		UserID: res.ID,
		Email:  res.Email,
	}
	token, err := td.Generate()
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}

	response := responses.User{
		ID:        res.ID,
		Name:      res.Name,
		Email:     res.Email,
		Token:     token,
		CreatedAt: &res.CreatedAt,
	}
	success(ctx, response)
}

func (s *Server) createUser(ctx *gin.Context) {
	var user requests.CreateUser
	if err := ctx.ShouldBindJSON(&user); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	dbUser := &users.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := s.userService.Create(ctx.Request.Context(), dbUser); err != nil {
		badRequestFromError(ctx, err)
		return
	}

	td := core.TokenData{
		UserID: dbUser.ID,
		Email:  dbUser.Email,
	}
	token, err := td.Generate()
	if err != nil {
		badRequestFromError(ctx, err)
		return
	}

	res := responses.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Token:     token,
		CreatedAt: &dbUser.CreatedAt,
	}
	created(ctx, res)
}
