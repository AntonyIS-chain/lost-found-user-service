package postgresDB

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserManagementService struct {
	dbClient *PostgresDBClient
}

// NewPostgresDBClient initializes a new GORM-based Postgres client
func NewUserManagementService(dbClient *PostgresDBClient) *UserManagementService {
	return &UserManagementService{dbClient: dbClient}
}

// RegisterUser inserts a new user with a hashed password
func (c *PostgresDBClient) RegisterUser(user domain.User) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	if err := c.DB.Create(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

// AuthenticateUser checks user credentials
func (c *PostgresDBClient) AuthenticateUser(email, password string) (domain.User, error) {
	var user domain.User

	if err := c.DB.Where("email = ?", email).First(&user).Error; err != nil {

		return domain.User{}, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		fmt.Println("err", err)
		return domain.User{}, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (c *PostgresDBClient) GetUserByID(userID string) (domain.User, error) {
	var user domain.User

	if err := c.DB.Where("ID = ?", userID).First(&user).Error; err != nil {

		return domain.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetUserByEmail retrieves a user by ID
func (c *PostgresDBClient) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	if err := c.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

// ListUsers retrieves all users with their roles
func (c *PostgresDBClient) ListUsers() ([]domain.User, error) {
	var users []domain.User

	if err := c.DB.Preload("Role").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// UpdateUser updates user details
func (c *PostgresDBClient) UpdateUser(userID string, updates domain.User) error {
	if err := c.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser removes a user from the database
func (c *PostgresDBClient) DeleteUser(userID string) error {
	if err := c.DB.Delete(&domain.User{}, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AssignRole assigns a role to a user
func (c *PostgresDBClient) AssignRole(userID string, roleID string) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		var user domain.User
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return fmt.Errorf("user not found")
		}

		var role domain.Role
		if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
			return fmt.Errorf("role not found")
		}

		if err := tx.Model(&user).Update("role_id", roleID).Error; err != nil {
			return fmt.Errorf("failed to assign role: %w", err)
		}

		return nil
	})
}

// ChangePassword allows a user to change their password
func (c *PostgresDBClient) ChangePassword(userID, oldPassword, newPassword string) error {
	var user domain.User

	if err := c.DB.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("incorrect old password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := c.DB.Model(&user).Update("password_hash", string(newHashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// Deactivate sets a user's status to inactive
func (c *PostgresDBClient) Deactivate(userID string) error {
	var user domain.User
	if err := c.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := c.DB.Model(&user).Updates(map[string]interface{}{"is_active": false}).Error; err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

func (c *PostgresDBClient) Activate(userID string) error {
	var user domain.User
	if err := c.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := c.DB.Model(&user).Updates(map[string]interface{}{"is_active": true}).Error; err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// ForgotPassword generates a password reset token (mocked implementation)

func (c *PostgresDBClient) ForgotPassword(email string) error {
	// Check if the user exists
	user, err := c.GetUserByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	// Generate a secure reset token
	resetToken := uuid.New().String()
	expiration := time.Now().Add(15 * time.Minute) // Token expires in 15 minutes

	// Store the reset token in the database
	err = c.StoreResetToken(user.ID, resetToken, expiration)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Construct the reset link
	// resetLink := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken)

	// Send email to user with reset link
	// subject := "Password Reset Request"
	// body := fmt.Sprintf("Click the following link to reset your password: %s\nThis link will expire in 15 minutes.", resetLink)

	// err = email.SendEmail(user.Email, subject, body) // Ensure `SendEmail` is implemented
	// if err != nil {
	// 	log.Printf("Failed to send reset email to %s: %v", user.Email, err)
	// 	return fmt.Errorf("failed to send reset email: %w", err)
	// }

	log.Printf("Password reset link sent to: %s", user.Email)
	return nil
}

// ResetPassword resets a user's password using a token
func (c *PostgresDBClient) ResetPassword(email, newPassword string) error {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	if err := c.DB.Model(&domain.User{}).Where("email = ?", email).Update("password_hash", string(newHashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	fmt.Println("Password reset successful for user ID:", email)
	return nil
}

// VerifyEmail verifies the user's email address
func (c *PostgresDBClient) VerifyEmail(token string) error {
	userID := "someUserID" // Retrieve actual user ID from DB using the token

	if err := c.DB.Model(&domain.User{}).Where("id = ?", userID).Update("email_verified", true).Error; err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	fmt.Println("Email verified for user ID:", userID)
	return nil
}

// VerifyEmail verifies the user's email address
func (c *PostgresDBClient) StoreResetToken(email, resetToken string, expiration time.Time) error {
	query := `INSERT INTO password_reset_tokens (email, token, expires_at) VALUES ($1, $2, $3) 
	          ON CONFLICT (email) DO UPDATE SET token = $2, expires_at = $3`

	err := c.DB.Exec(query, email, resetToken, expiration)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %v", err)
	}

	return nil
}
