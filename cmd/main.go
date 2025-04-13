package cmd

import (
	"log"

	"github.com/AntonyIS-chain/lost-found-user-service/config"
	app "github.com/AntonyIS-chain/lost-found-user-service/internal/adapters/app/handlers"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/adapters/repository"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/services"
	"github.com/AntonyIS-chain/lost-found-user-service/pkg"
)

func RunService() {
	// Load configuration
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database client
	dbClient, err := repository.NewPostgresDBClient(conf)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	authRepo := repository.AuthRepository(dbClient)
	roleRepo := repository.NewRoleRepository(dbClient)
	userRepo := repository.NewUserRepository(dbClient)

	// Initialize services

	rolesService := services.NewRoleManagementService(roleRepo)
	authService := services.NewAuthManagementService(authRepo, rolesService)
	usersService := services.NewUserManagementService(userRepo, rolesService)

	pkg.SeedDB(authService, usersService, rolesService)

	// Start HTTP server with initialized services
	app.InitGinRoutes(authService, usersService, rolesService, conf)
}
