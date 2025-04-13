package repository

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
)

type AuthManagementRepository struct {
	dbClient *PostgresDBClient
}

func AuthRepository(dbClient *PostgresDBClient) *AuthManagementRepository {
	return &AuthManagementRepository{dbClient: dbClient}
}

func (a *AuthManagementRepository) SignUp(user domain.User) (domain.User, error) {
	var existingUser domain.User
	fmt.Println("User", user.Email, user.ID)

	// Check by email
	if err := a.dbClient.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return domain.User{}, errors.New("user already exists")
	}


	err := a.dbClient.DB.Create(&user).Error
	if err != nil {
		// Gracefully handle duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return domain.User{}, errors.New("user with this ID or email already exists")
		}
		return domain.User{}, errors.New("user already exists")
	}

	return user, nil
}


func (a *AuthManagementRepository) SignIn(email, password string) (domain.User, error) {
	var user domain.User
	if err := a.dbClient.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.User{}, errors.New("invalid credentials")
	}

	return user, nil

}

func (a *AuthManagementRepository) ForgotPassword(email string) (domain.User, error) {
	var user domain.User
	if err := a.dbClient.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, errors.New("user not found")
	}

	return user, nil

}

func (a *AuthManagementRepository) ResetPassword(email, oldPassword, newPassword string) (domain.User, error) {
	var user domain.User

	// Find user by email
	if err := a.dbClient.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return domain.User{}, fmt.Errorf("incorrect old password: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := a.dbClient.DB.Save(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("failed to update password: %w", err)
	}

	return user, nil
}
