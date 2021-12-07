package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	login_dto "github.com/ydhnwb/go-refresh-token-example/dto/login"
	user_dto "github.com/ydhnwb/go-refresh-token-example/dto/user"
	"github.com/ydhnwb/go-refresh-token-example/repo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(loginRequest login_dto.LoginRequest) (*login_dto.LoginResponse, error)
	RefreshToken(refreshTokenRequest login_dto.RefreshTokenRequest) (*login_dto.RefreshAccessTokenResponse, error)
}

type authService struct {
	// authRepo repo.AuthRepoInterface
	db       *gorm.DB
	userRepo repo.UserRepoInterface
}

func NewAuthService(db *gorm.DB, userRepo repo.UserRepoInterface) AuthServiceInterface {
	return &authService{
		db:       db,
		userRepo: userRepo,
	}
}

func (c *authService) RefreshToken(refreshTokenRequest login_dto.RefreshTokenRequest) (*login_dto.RefreshAccessTokenResponse, error) {
	user, err := c.validaterefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tokenExpiredAt := now.Add(time.Minute * time.Duration(1))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":               user.ID,
		"token_expired_at": tokenExpiredAt,
	})

	tokenString, _ := token.SignedString([]byte("ydhnwb"))

	loginRes := login_dto.RefreshAccessTokenResponse{
		ID:                    user.ID,
		Name:                  user.Name,
		Email:                 user.Email,
		AccessToken:           tokenString,
		AccesssTokenExpiredAt: tokenExpiredAt,
		Message:               "Only gives you new access token. Refresh token not generated. If your refresh token is expired. You need to login again from the app/beginning",
	}

	return &loginRes, nil
}

func (c *authService) Login(loginRequest login_dto.LoginRequest) (*login_dto.LoginResponse, error) {
	res, err := c.userRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, err
	}

	isPasswordValid := comparePassword(res.Password, []byte(loginRequest.Password))
	if !isPasswordValid {
		return nil, errors.New("Kombinasi password dan email salah")
	}

	now := time.Now()
	tokenExpiredAt := now.Add(time.Minute * time.Duration(1))
	// refreshTokenExpiredAt := now.AddDate(0, 1, 0)
	refreshTokenExpiredAt := now.Add(time.Minute * time.Duration(2))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":               res.ID,
		"token_expired_at": tokenExpiredAt,
	})

	tokenString, _ := token.SignedString([]byte("ydhnwb"))

	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":                       res.ID,
		"refresh_token_expired_at": refreshTokenExpiredAt,
	})
	tokenStringRefresh, _ := tokenRefresh.SignedString([]byte("ydhnwb"))

	loginResponse := login_dto.LoginResponse{
		ID:                    res.ID,
		Name:                  res.Name,
		Email:                 res.Email,
		AccessToken:           tokenString,
		AccesssTokenExpiredAt: tokenExpiredAt,
		RefreshToken:          tokenStringRefresh,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
	}

	return &loginResponse, nil
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (c *authService) validaterefreshToken(refreshToken string) (*user_dto.UserResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Refresh token is not valid")
		}
		return []byte("ydhnwb"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("Refresh Token check: OK! -> %v", claims["id"])
		tokenExpiredAt := claims["refresh_token_expired_at"]
		expired := tokenExpiredAt.(string)
		asDate, _ := time.Parse(time.RFC3339, expired)

		if time.Now().After(asDate) {
			return nil, errors.New("Refresh token is expired. Please relogin from the app.")
		}

		id := claims["id"].(float64)
		// u64, _ := strconv.ParseUint(id, 10, 32)

		uid := uint(id)
		user, err := c.userRepo.FindUserByID(uid)
		if err != nil {
			return nil, err
		}

		userRes := user_dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		return &userRes, nil

	}
	return nil, errors.New("refresh token is not provided or is not valid")
}
