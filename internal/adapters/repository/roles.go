package repository

import (
	"fmt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"gorm.io/gorm"
)

type RoleRepository struct {
	dbClient *PostgresDBClient
}

func NewRoleRepository(dbClient *PostgresDBClient) *RoleRepository {
	return &RoleRepository{dbClient: dbClient}
}

func (r *RoleRepository) CreateRole(role domain.Role) (domain.Role, error) {
	if err := r.dbClient.DB.Create(&role).Error; err != nil {
		return domain.Role{}, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) GetRoleByName(roleName string) (domain.Role, error) {
	var role domain.Role
	err := r.dbClient.DB.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Role{}, fmt.Errorf("role '%s' not found", roleName)
		}
		return domain.Role{}, fmt.Errorf("failed to retrieve role: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) ListRoles() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.dbClient.DB.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	return roles, nil
}

func (r *RoleRepository) UpdateRole(roleID int, updates domain.Role) error {
	result := r.dbClient.DB.Model(&domain.Role{}).Where("id = ?", roleID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no role found with ID %d", roleID)
	}
	return nil
}

func (r *RoleRepository) DeleteRole(roleID int) error {
	result := r.dbClient.DB.Where("id = ?", roleID).Delete(&domain.Role{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no role found with ID %d", roleID)
	}
	return nil
}

func (r *RoleRepository) AssignRoleToUser(userID string, roleName string) error {
	db := r.dbClient.DB

	// Fetch role
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("could not get role by name: %w", err)
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var userRole domain.UserRole
	err = tx.Where("user_id = ?", userID).First(&userRole).Error

	switch {
	case err == nil:
		// Update existing role
		if err := tx.Model(&userRole).Update("role_id", role.ID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update user role: %w", err)
		}
	case err == gorm.ErrRecordNotFound:
		// Assign new role
		newUserRole := domain.UserRole{UserID: userID, RoleID: role.ID}
		if err := tx.Create(&newUserRole).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign role: %w", err)
		}
	default:
		tx.Rollback()
		return fmt.Errorf("error checking user role: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
