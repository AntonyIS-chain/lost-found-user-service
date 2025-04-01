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
		{Name: "User Admin", Description: "Administrator with full access to the system"},
		{Name: "Moderator", Description: "Can moderate user content"},
		{Name: "Registered User", Description: "Standard user with basic privileges"},
		{Name: "Guest User", Description: "Limited access user"},
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

func SeedUsers(userService ports.UserService, roleService ports.RoleService) {
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

	var roleMap = make(map[string]int)
	for _, role := range existingRoles {
		roleMap[role.Name] = role.ID
	}

	users := []domain.User{
		{FirstName: "Admin", LastName: "User", Email: "admin@example.com", PasswordHash: "Admin@123", RoleID: roleMap["User Admin"], RoleName: "User Admin"},
		{FirstName: "Moderator", LastName: "User", Email: "moderator@example.com", PasswordHash: "Moderator@123", RoleID: roleMap["Moderator"], RoleName: "Moderator"},
		{FirstName: "John", LastName: "Doe", Email: "user@example.com", PasswordHash: "User@123", RoleID: roleMap["Registered User"], RoleName: "Registered User"},
		{FirstName: "Guest", LastName: "User", Email: "guest@example.com", PasswordHash: "Guest@123", RoleID: roleMap["Guest User"], RoleName: "Guest User"},
	}

	for _, user := range users {
		if existingUserMap[user.Email] {
			// log.Printf("User '%s' already exists. Skipping seeding.\n", user.Email)
			continue
		}
		createdUser, err := userService.RegisterUser(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			continue
		}

		err = roleService.AssignRoleToUser(createdUser.ID, user.RoleName)
		if err != nil {
			log.Printf("Failed to assign role to user %s: %v", user.Email, err)
			continue
		}
	}

	log.Println("Successfully seeded all users and assigned roles.")
}
