package app

import (
	"fmt"
	"log"
	"time"

	"github.com/AntonyIS-chain/lost-found-user-service/config"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/adapters/app/controllers"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinRoutes(userSvc ports.UserService, roleSvc ports.RoleService, config *config.Config) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Controllers
	userController := controllers.NewUserController(userSvc)
	roleController := controllers.NewRoleController(roleSvc, userSvc)

	// User Routes
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", userController.RegisterUser)
		authRoutes.POST("/login", userController.AuthenticateUser)
		authRoutes.POST("/refresh-token", userController.RefreshToken)
		authRoutes.POST("/:id/change-password", userController.ChangePassword)
		authRoutes.POST("/forgot-password", userController.ForgotPassword)
		authRoutes.POST("/reset-password", userController.ResetPassword)
		authRoutes.GET("/verify-email/:token", userController.VerifyEmail)
	}

	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/", userController.ListUsers)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.GET("/deactivate/:id", userController.DeactivateUser)
		userRoutes.GET("/activate/:id", userController.ActivateUser)
		userRoutes.DELETE("/delete/:adminID/:userID", userController.DeleteUser)
		userRoutes.POST("/roles/assign-role", userController.AssignRole)
	}

	// Role Routes
	roleRoutes := router.Group("/api/v1/users/roles")
	{
		roleRoutes.POST("/", roleController.CreateRole)
		roleRoutes.GET("/:role_name", roleController.GetRoleByName)
		roleRoutes.GET("/", roleController.ListRoles)
		roleRoutes.PUT("/:role_id", roleController.UpdateRole)
		roleRoutes.DELETE("/:role_id", roleController.DeleteRole)
		roleRoutes.POST("/assign", roleController.AssignRoleToUser)
	}

	// Start server
	log.Println("Starting server on port", config.SERVER_PORT)
	router.Run(fmt.Sprintf(":%s", config.SERVER_PORT))
}
