package cmd

import (
	"log"

	"github.com/AntonyIS-chain/lost-found-user-service/config"
	app "github.com/AntonyIS-chain/lost-found-user-service/internal/adapters/app/handlers"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/adapters/repository/postgresDB"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/services"
)

func RunService() {
	// Load configuration
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database client
	dbClient, err := postgresDB.NewPostgresDBClient(conf)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	// Initialize services
	rolesService := services.NewRoleManagementService(dbClient)
	usersService := services.NewUserManagementService(dbClient, rolesService)

	// Seed "User Admin" role
	// pkg.SeedRoles(rolesService)
	// pkg.SeedUsers(usersService, rolesService)

	// Start HTTP server with initialized services
	app.InitGinRoutes(usersService, rolesService, conf)
}
