package services

import (
	"fmt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
)

type RoleManagementService struct {
	repo ports.RoleRepository
}

// Constructor
func NewRoleManagementService(repo ports.RoleRepository) *RoleManagementService {
	return &RoleManagementService{
		repo: repo,
	}
}

// Implement RoleService interface

func (s *RoleManagementService) CreateRole(role domain.Role) (domain.Role, error) {
	// Check if role already exists
	existingRole, err := s.GetRoleByName(role.Name)
	if err == nil && existingRole.ID != 0 {
		return domain.Role{}, fmt.Errorf("role '%s' already exists", role.Name)
	}

	// Proceed to create role if it does not exist
	return s.repo.CreateRole(role)
}

func (s *RoleManagementService) GetRoleByName(roleName string) (domain.Role, error) {
	return s.repo.GetRoleByName(roleName)
}

func (s *RoleManagementService) ListRoles() ([]domain.Role, error) {
	return s.repo.ListRoles()
}

func (s *RoleManagementService) UpdateRole(roleID int, updates domain.Role) error {
	return s.repo.UpdateRole(roleID, updates)
}

func (s *RoleManagementService) DeleteRole(roleID int) error {
	return s.repo.DeleteRole(roleID)
}

func (s *RoleManagementService) AssignRoleToUser(userID string, roleName string) error {
	return s.repo.AssignRoleToUser(userID, roleName)
}
