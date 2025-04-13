package repository

import (
	"errors"
	"fmt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)



type UserRepository struct {
	dbClient *PostgresDBClient
}

// Constructor
func NewUserRepository(dbClient *PostgresDBClient) *UserRepository {
	return &UserRepository{dbClient: dbClient}
}

// GetUserByID fetches a user by ID.
func (u *UserRepository) GetUserByID(userID string) (domain.User, error) {
	var user domain.User
	if err := u.dbClient.DB.First(&user, "id = ?", userID).Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// GetUserByEmail fetches a user by email.
func (u *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	if err := u.dbClient.DB.First(&user, "email = ?", email).Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// ListUsers returns all users.
func (u *UserRepository) ListUsers() ([]domain.User, error) {
	var users []domain.User
	if err := u.dbClient.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

// UpdateUser updates user fields.
func (u *UserRepository) UpdateUser(userID string, updates domain.User) error {
	result := u.dbClient.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %s", userID)
	}
	return nil
}

// DeleteUser deletes a user permanently.
func (u *UserRepository) DeleteUser(userID string) error {
	result := u.dbClient.DB.Where("id = ?", userID).Delete(&domain.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %s", userID)
	}
	return nil
}

// Deactivate sets the user as inactive.
func (u *UserRepository) Deactivate(userID string) error {
	return u.setActiveStatus(userID, false)
}

// Activate sets the user as active.
func (u *UserRepository) Activate(userID string) error {
	return u.setActiveStatus(userID, true)
}

// Helper for activation toggle
func (u *UserRepository) setActiveStatus(userID string, active bool) error {
	result := u.dbClient.DB.Model(&domain.User{}).Where("id = ?", userID).Update("is_active", active)
	if result.Error != nil {
		return fmt.Errorf("failed to update active status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %s", userID)
	}
	return nil
}

// AssignRole associates a role to the user (assuming user has a role_id foreign key).
func (u *UserRepository) AssignRole(userID string, roleID string) error {
	result := u.dbClient.DB.Model(&domain.User{}).Where("id = ?", userID).Update("role_id", roleID)
	if result.Error != nil {
		return fmt.Errorf("failed to assign role: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %s", userID)
	}
	return nil
}

// ChangePassword validates old password and updates to new hashed password.
func (u *UserRepository) ChangePassword(userID, oldPassword, newPassword string) error {
	var user domain.User
	if err := u.dbClient.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("incorrect old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	if err := u.dbClient.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// // VerifyEmail activates a user's email based on a token (example implementation).
// func (u *UserRepository) VerifyEmail(token string) error {
// 	var user domain.User
// 	if err := u.dbClient.DB.First(&user, "verification_token = ?", token).Error; err != nil {
// 		return fmt.Errorf("invalid or expired token: %w", err)
// 	}

// 	user.IsVerified = true
// 	user.VerificationToken = "" // Clear token after verification
// 	if err := u.dbClient.DB.Save(&user).Error; err != nil {
// 		return fmt.Errorf("failed to verify email: %w", err)
// 	}
// 	return nil
// }
