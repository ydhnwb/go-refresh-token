package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ydhnwb/go-refresh-token-example/config"
	"github.com/ydhnwb/go-refresh-token-example/handler"
	"github.com/ydhnwb/go-refresh-token-example/middleware"
	"github.com/ydhnwb/go-refresh-token-example/repo"
	"github.com/ydhnwb/go-refresh-token-example/service"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB                     = config.SetupDatabaseConnection()
	userRepo    repo.UserRepoInterface       = repo.NewUserRepo(db)
	authService service.AuthServiceInterface = service.NewAuthService(db, userRepo)
	userService service.UserServiceInterface = service.NewUserService(userRepo)
	jwtService  service.JWTServiceInterface  = service.NewJWTService(db, userRepo)
	authHandler handler.AuthHandlerInterface = handler.NewAuthHandler(userService, authService)
	userHandler handler.UserHandlerInterface = handler.NewUserHandler(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	server := gin.Default()

	authRoutes := server.Group("api/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/generate-new-access-token", authHandler.RefreshToken)
	}

	// user must be logged in
	userRoutes := server.Group("api/user", middleware.ValidateJWTToken())
	{
		userRoutes.GET("/profile", userHandler.MyProfile)
	}

	server.Run()

}
