package services

import (
	"errors"
	"fmt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-user-service/pkg"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserManagementService struct {
	repo        ports.UserRepository
	roleService ports.RoleService
}

// Constructor
func NewUserManagementService(repo ports.UserRepository, roleService ports.RoleService) *UserManagementService {
	return &UserManagementService{
		repo:        repo,
		roleService: roleService,
	}
}

// Implement UserService interface
func (s *UserManagementService) RegisterUser(user domain.User) (domain.User, error) {
	// Check if user already exists by email
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err == nil && existingUser.ID != "" {
		return domain.User{}, fmt.Errorf("user with email '%s' already exists", user.Email)
	}

	// Get role by name
	role, err := s.roleService.GetRoleByName(user.RoleName)
	if err != nil || role.ID == 0 {
		return domain.User{}, fmt.Errorf("role '%s' does not exist", user.RoleName)
	}

	// Assign RoleID instead of Role struct
	userId := uuid.New().String()
	user.ID = userId
	user.RoleID = role.ID

	return s.repo.RegisterUser(user)
}

func (s *UserManagementService) AuthenticateUser(email, password string) (domain.User, error) {
	return s.repo.AuthenticateUser(email, password)
}

func (s *UserManagementService) GetUserByID(userID string) (domain.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *UserManagementService) GetUserByEmail(email string) (domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserManagementService) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}

func (s *UserManagementService) UpdateUser(userID string, updates domain.User) error {
	return s.repo.UpdateUser(userID, updates)
}

func (s *UserManagementService) DeleteUser(adminID, userID string) error {
	adminUser, err := s.GetUserByID(adminID)
	if err != nil {
		return fmt.Errorf("failed to retrieve admin user: %w", err)
	}

	// Ensure the requesting user is an admin
	if adminUser.RoleName != "User Admin" {
		return fmt.Errorf("unauthorized: only admins can delete users")
	}

	// Check if the user to be deleted exists
	_, err = s.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with ID %s not found", userID)
		}
		return fmt.Errorf("failed to retrieve user for deletion: %w", err)
	}

	return s.repo.DeleteUser(userID)
}

func (s *UserManagementService) Deactivate(userID string) error {
	return s.repo.Deactivate(userID)
}

func (s *UserManagementService) Activate(userID string) error {
	return s.repo.Activate(userID)
}

func (s *UserManagementService) RefreshToken(refreshToken string) (string, error) {
	email, err := pkg.ValidateRefreshToken(refreshToken)

	if err != nil {
		return "", err
	}

	return email, nil
}

func (s *UserManagementService) AssignRole(userID string, roleID string) error {
	return s.repo.AssignRole(userID, roleID)
}

func (s *UserManagementService) ChangePassword(userID string, oldPassword, newPassword string) error {
	return s.repo.ChangePassword(userID, oldPassword, newPassword)
}

func (s *UserManagementService) ForgotPassword(email string) error {
	return s.repo.ForgotPassword(email)
}

func (s *UserManagementService) ResetPassword(token, newPassword string) error {
	email, err := pkg.ValidateRefreshToken(token)

	if err != nil {
		return err
	}

	return s.repo.ResetPassword(email, newPassword)
}

func (s *UserManagementService) VerifyEmail(token string) error {
	return s.repo.VerifyEmail(token)
}
