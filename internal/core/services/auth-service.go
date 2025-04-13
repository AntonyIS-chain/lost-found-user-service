package services

import (
	"fmt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-user-service/pkg"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthManagementService struct {
	repo    ports.AuthRepository
	roleSvc ports.RoleService
}

func NewAuthManagementService(repo ports.AuthRepository, roleSvc ports.RoleService) *AuthManagementService {
	return &AuthManagementService{
		repo: repo,
		roleSvc: roleSvc,
	}
}

func (a AuthManagementService) SignUp(user domain.User) (domain.SessionResponse, error) {

	user.ID = uuid.New().String()

	// Hash the plain password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Save user via repository
	createdUser, err := a.repo.SignUp(user)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens with user ID and role
	accessToken, err := pkg.GenerateToken(createdUser.ID, createdUser.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := pkg.GenerateRefreshToken(createdUser.ID, createdUser.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Assign role after user creation
	err = a.roleSvc.AssignRoleToUser(createdUser.ID, createdUser.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to create role: %w", err)

	}

	// Return session response
	return domain.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: domain.SessionUser{
			ID:   createdUser.ID,
			Role: createdUser.RoleName,
		},
	}, nil
}

func (a AuthManagementService) SignIn(email, password string) (domain.SessionResponse, error) {
	user, err := a.repo.SignIn(email, password)

	if err != nil {
		return domain.SessionResponse{}, err
	}
	// Generate tokens with user ID and role
	accessToken, err := pkg.GenerateToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Return session response
	return domain.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: domain.SessionUser{
			ID:   user.ID,
			Role: user.RoleName,
		},
	}, nil
}

func (a AuthManagementService) ForgotPassword(email string) (domain.SessionResponse, error) {
	user, err := a.repo.ForgotPassword(email)

	if err != nil {
		return domain.SessionResponse{}, err
	}
	// Generate tokens with user ID and role
	accessToken, err := pkg.GenerateToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Return session response
	return domain.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: domain.SessionUser{
			ID:   user.ID,
			Role: user.RoleName,
		},
	}, nil
}

func (a *AuthManagementService) ResetPassword(email, oldPassword, newPassword string) (domain.SessionResponse, error) {
	user, err := a.repo.ResetPassword(email, oldPassword, newPassword)

	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new tokens (access token, refresh token)
	accessToken, err := pkg.GenerateToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.ID, user.RoleName)
	if err != nil {
		return domain.SessionResponse{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Return the session response with new tokens
	return domain.SessionResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		User: domain.SessionUser{
			ID:   user.ID,
			Role: user.RoleName,
		},
	}, nil
}
