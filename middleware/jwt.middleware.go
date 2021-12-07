package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	common_dto "github.com/ydhnwb/go-refresh-token-example/dto/common"
)

func ValidateJWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common_dto.BuildErrorResponse(401, "Access Token is not provided"))
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error middleware: %v", "Access Token is not valid")
			}
			return []byte("ydhnwb"), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common_dto.BuildErrorResponse(401, err.Error()))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Printf("Token check: OK! -> %v", claims["id"])
			tokenExpiredAt := claims["token_expired_at"]
			expired := tokenExpiredAt.(string)
			asDate, _ := time.Parse(time.RFC3339, expired)

			if time.Now().After(asDate) {
				response := common_dto.BuildErrorResponse(401, "Your token is expired")
				c.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
			// c.Next()
		} else {
			response := common_dto.BuildErrorResponse(401, "Access Token is not provided.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
	}
}
