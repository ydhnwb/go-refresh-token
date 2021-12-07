package handler

import (
	"github.com/gin-gonic/gin"
	common_dto "github.com/ydhnwb/go-refresh-token-example/dto/common"
	"github.com/ydhnwb/go-refresh-token-example/service"
)

type UserHandlerInterface interface {
	MyProfile(ctx *gin.Context)
}

type userHandler struct {
	userService service.UserServiceInterface
	jwtService  service.JWTServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface, jwtService service.JWTServiceInterface) UserHandlerInterface {
	return &userHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userHandler) MyProfile(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	user, err := c.jwtService.ExtractToken(token)
	if err != nil {
		res := common_dto.BuildErrorResponse(400, err.Error())
		ctx.JSON(400, res)
		return
	}

	ctx.JSON(200, user)
}
