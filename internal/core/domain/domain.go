package domain

import (
	"time"

	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")
)

type User struct {
	ID           string    `json:"id" db:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email"  db:"email"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	RoleID       int       `json:"role_id" db:"role_id"`
	Role         Role      `json:"role" db:"-"`
	RoleName     string    `json:"role_name" db:"role_name"`
	IsActive     bool      `json:"is_active" db:"is_active" gorm:"default:true"`
	Token        string    `json:"token" db:"token"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserToken struct {
	UserID string `json:"user_id" db:"user_id"`
	Token  string `json:"token"  db:"token"`
}

type Role struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"not null;index"`
	RoleID int    `gorm:"not null;index"`
}

type LogMessage struct {
	LogLevel string `json:"log_level"`
	Message  string `json:"message"`
	Service  string `json:"service"`
}

type Session struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Success    bool        `json:"success"`
	Results    interface{} `json:"results"`
}

type ResetPasswordToken struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
