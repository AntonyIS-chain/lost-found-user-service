package pkg

import (
	"log"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
)

func SeedRoles(roleService ports.RoleService) {
	existingRoles, err := roleService.ListRoles()
	if err != nil {
		log.Printf("Failed to fetch roles: %v", err)
		return
	}

	roleNames := []domain.Role{
		{Name: "Administrator", Description: "Administrator with full access to the system"},
		{Name: "Moderator", Description: "Can moderate user content"},
		{Name: "Guest", Description: "Limited access user"},
	}

	existingRoleMap := make(map[string]bool)
	for _, role := range existingRoles {
		existingRoleMap[role.Name] = true
	}

	for _, role := range roleNames {
		if existingRoleMap[role.Name] {
			// log.Printf("Role '%s' already exists. Skipping seeding.\n", role.Name)
			continue
		}

		_, err := roleService.CreateRole(role)
		if err != nil {
			log.Printf("Failed to create role '%s': %v", role.Name, err)
			continue
		}
	}

	log.Println("Successfully seeded all roles.")
}

// SeedUsers populates the database with initial users and roles
func SeedDB(authService ports.AuthService, userService ports.UserService, roleService ports.RoleService) {
	SeedRoles(roleService)
	existingUsers, err := userService.ListUsers()
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return
	}

	existingUserMap := make(map[string]bool)
	for _, user := range existingUsers {
		existingUserMap[user.Email] = true
	}

	existingRoles, err := roleService.ListRoles()
	if err != nil {
		log.Printf("Failed to fetch roles: %v", err)
		return
	}

	roleMap := make(map[string]int)
	for _, role := range existingRoles {
		roleMap[role.Name] = role.ID
	}

	users := []domain.User{
		{FirstName: "John", LastName: "Doe", Email: "admin@example.com", RoleName: "Administrator", Phone: "+254723308900"},
		{FirstName: "Mark", LastName: "Tess", Email: "moderator@example.com", RoleName: "Moderator", Phone: "+254723308900"},
		{FirstName: "Mike", LastName: "Tesla", Email: "guest@example.com", RoleName: "Guest", Phone: "+254723308900"},
	}

	for _, user := range users {
		if existingUserMap[user.Email] {
			continue
		}

		// Validate role existence
		roleID, exists := roleMap[user.RoleName]
		if !exists {
			log.Printf("Role '%s' not found. Skipping user %s.\n", user.RoleName, user.Email)
			continue
		}
		user.RoleID = roleID
		user.PasswordHash = "Password@1234"

		_, err := authService.SignUp(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			continue
		}

	}

	log.Println("Successfully seeded all users and assigned roles.")
}
