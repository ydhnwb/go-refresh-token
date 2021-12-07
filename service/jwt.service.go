package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt"
	user_dto "github.com/ydhnwb/go-refresh-token-example/dto/user"
	"github.com/ydhnwb/go-refresh-token-example/repo"
	"gorm.io/gorm"
)

type JWTServiceInterface interface {
	ExtractToken(tokenString string) (*user_dto.UserResponse, error)
}

type jwtService struct {
	db       *gorm.DB
	userRepo repo.UserRepoInterface
}

func NewJWTService(db *gorm.DB, userRepo repo.UserRepoInterface) JWTServiceInterface {
	return &jwtService{
		db:       db,
		userRepo: userRepo,
	}
}

func (c *jwtService) ExtractToken(tokenString string) (*user_dto.UserResponse, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error middleware: %v", "Token is not valid")
		}
		return []byte("ydhnwb"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("Token check: OK! -> %v", claims["id"])
		t := fmt.Sprintf("%v", claims["id"])
		u64, err := strconv.ParseUint(t, 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		uid := uint(u64)
		user, e := c.userRepo.FindUserByID(uid)
		if e != nil {
			return nil, e
		}
		userResponse := user_dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		return &userResponse, nil
	} else {
		err = errors.New("token cannot be extracted")
		return nil, err
	}
}
