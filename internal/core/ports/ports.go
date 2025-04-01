package ports

import "github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"

type UserService interface {
	RegisterUser(user domain.User) (domain.User, error)
	AuthenticateUser(email, password string) (domain.User, error)
	GetUserByID(userID string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	ListUsers() ([]domain.User, error)
	UpdateUser(userID string, updates domain.User) error
	DeleteUser(adminID, userID string) error
	Deactivate(userID string) error
	Activate(userID string) error
	RefreshToken(refresh_token string) (string, error)
	AssignRole(userID string, roleID string) error
	ChangePassword(userID string, oldPassword, newPassword string) error
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
	VerifyEmail(token string) error
}

type UserRepository interface {
	RegisterUser(user domain.User) (domain.User, error)
	AuthenticateUser(email, password string) (domain.User, error)
	GetUserByID(userID string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	ListUsers() ([]domain.User, error)
	UpdateUser(userID string, updates domain.User) error
	DeleteUser(userID string) error
	Deactivate(userID string) error
	Activate(userID string) error
	AssignRole(userID string, roleID string) error
	ChangePassword(userID string, oldPassword, newPassword string) error
	ForgotPassword(email string) error
	ResetPassword(userID, newPassword string) error
	VerifyEmail(token string) error
}

type RoleRepository interface {
	CreateRole(role domain.Role) (domain.Role, error)
	GetRoleByName(roleName string) (domain.Role, error)
	ListRoles() ([]domain.Role, error)
	UpdateRole(roleID int, updates domain.Role) error
	DeleteRole(roleID int) error
	AssignRoleToUser(userID string, roleName string) error
}

type RoleService interface {
	CreateRole(role domain.Role) (domain.Role, error)
	GetRoleByName(roleName string) (domain.Role, error)
	ListRoles() ([]domain.Role, error)
	UpdateRole(roleID int, updates domain.Role) error
	DeleteRole(roleID int) error
	AssignRoleToUser(userID string, roleName string) error
}

type LoggingService interface {
	SendLog(LogEntry domain.LogMessage)
	LogDebug(LogEntry domain.LogMessage)
	LogInfo(LogEntry domain.LogMessage)
	LogWarning(LogEntry domain.LogMessage)
	LogError(LogEntry domain.LogMessage)
}
