package service

import (
	"github.com/mashingan/smapping"
	common_dto "github.com/ydhnwb/go-refresh-token-example/dto/common"
	register_dto "github.com/ydhnwb/go-refresh-token-example/dto/register"
	user_dto "github.com/ydhnwb/go-refresh-token-example/dto/user"
	"github.com/ydhnwb/go-refresh-token-example/entity"
	"github.com/ydhnwb/go-refresh-token-example/repo"
)

type UserServiceInterface interface {
	CreateAccount(createUserRequest register_dto.RegisterRequest) (*common_dto.StdResponse, error)
	FindUserByEmail(email string) (*user_dto.UserResponse, error)
	FindUserByID(id uint) (*user_dto.UserResponse, error)
}

type userService struct {
	userRepo repo.UserRepoInterface
}

func NewUserService(userRepo repo.UserRepoInterface) UserServiceInterface {
	return &userService{
		userRepo: userRepo,
	}
}

func (c *userService) CreateAccount(createUserRequest register_dto.RegisterRequest) (*common_dto.StdResponse, error) {
	var user entity.User
	err := smapping.FillStruct(&user, smapping.MapFields(&createUserRequest))
	if err != nil {
		return nil, err
	}

	_, e := c.userRepo.CreateAccount(user)
	if e != nil {
		return nil, e
	}

	stdRes := common_dto.StdResponse{
		Code: 201,
		Msg:  "User created. Please login",
	}

	return &stdRes, nil
}

func (c *userService) FindUserByEmail(email string) (*user_dto.UserResponse, error) {
	user, e := c.userRepo.FindUserByEmail(email)
	if e != nil {
		return nil, e
	}

	userResponse := user_dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return &userResponse, nil
}

func (c *userService) FindUserByID(id uint) (*user_dto.UserResponse, error) {
	user, e := c.userRepo.FindUserByID(id)
	if e != nil {
		return nil, e
	}

	userResponse := user_dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return &userResponse, nil
}
