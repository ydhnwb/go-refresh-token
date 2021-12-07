package handler

import (
	"github.com/gin-gonic/gin"
	common_dto "github.com/ydhnwb/go-refresh-token-example/dto/common"
	login_dto "github.com/ydhnwb/go-refresh-token-example/dto/login"
	register_dto "github.com/ydhnwb/go-refresh-token-example/dto/register"
	"github.com/ydhnwb/go-refresh-token-example/service"
)

type AuthHandlerInterface interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type authHandler struct {
	userService service.UserServiceInterface
	authService service.AuthServiceInterface
}

func NewAuthHandler(userService service.UserServiceInterface, authService service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{
		userService: userService,
		authService: authService,
	}
}

func (c *authHandler) Login(ctx *gin.Context) {
	var loginRequest login_dto.LoginRequest
	err := ctx.ShouldBind(&loginRequest)

	if err != nil {
		res := common_dto.BuildErrorResponse(400, err.Error())
		ctx.JSON(400, res)
		return
	}

	result, e := c.authService.Login(loginRequest)
	if e != nil {
		res := common_dto.BuildErrorResponse(401, e.Error())
		ctx.JSON(401, res)
		return
	}

	ctx.JSON(200, result)
}

func (c *authHandler) Register(ctx *gin.Context) {
	var createUserRequest register_dto.RegisterRequest
	e := ctx.ShouldBind(&createUserRequest)

	if e != nil {
		res := common_dto.BuildErrorResponse(400, e.Error())
		ctx.JSON(400, res)
		return
	}

	result, err := c.userService.CreateAccount(createUserRequest)
	if err != nil {
		res := common_dto.BuildErrorResponse(422, err.Error())
		ctx.JSON(422, res)
		return
	}

	ctx.JSON(201, result)
}

func (c *authHandler) RefreshToken(ctx *gin.Context) {
	var refreshTokenRequest login_dto.RefreshTokenRequest
	e := ctx.ShouldBind(&refreshTokenRequest)

	if e != nil {
		res := common_dto.BuildErrorResponse(400, e.Error())
		ctx.JSON(400, res)
		return
	}

	refreshTokenRes, err := c.authService.RefreshToken(refreshTokenRequest)
	if err != nil {
		res := common_dto.BuildErrorResponse(401, err.Error())
		ctx.JSON(401, res)
		return
	}

	ctx.JSON(200, refreshTokenRes)
}
