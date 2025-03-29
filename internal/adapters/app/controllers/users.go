package controllers

import (
	"net/http"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-user-service/pkg"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service ports.UserService
}

func NewUserController(service ports.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":    "Bad request",
			"success":    false,
			"statusCode": 400,
		})
		return
	}

	createdUser, err := uc.service.RegisterUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (uc *UserController) AuthenticateUser(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.service.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message":    "Invalid email or password",
			"success":    false,
			"statusCode": 401,
		})
		return
	}

	// Generate tokens
	accessToken, err := pkg.GenerateToken(user.Email) // Short-lived access token
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.Email) // Long-lived refresh token
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Return tokens
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) RefreshToken(ctx *gin.Context) {
	var refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := uc.service.RefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message":    "Invalid credentials",
			"success":    false,
			"statusCode": 401,
		})
		return
	}

	// Generate tokens
	accessToken, err := pkg.GenerateToken(email) 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := pkg.GenerateRefreshToken(email) 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Return tokens
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	user, err := uc.service.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) ListUsers(ctx *gin.Context) {
	users, err := uc.service.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	var updates domain.User
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.service.UpdateUser(userID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	adminID := ctx.Param("adminID")
	userID := ctx.Param("userID")
	if err := uc.service.DeleteUser(adminID, userID); err != nil {
		if err.Error() == "failed to retrieve user for deletion: user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uc *UserController) DeactivateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if err := uc.service.Deactivate(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

func (uc *UserController) ActivateUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if err := uc.service.Activate(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User activated successfully"})
}

func (uc *UserController) AssignRole(ctx *gin.Context) {
	var req struct {
		UserID string `json:"user_id"`
		RoleID string `json:"role_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.Param("id")
	if err := uc.service.AssignRole(userID, req.RoleID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

func (uc *UserController) ChangePassword(ctx *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.Param("id")
	if err := uc.service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (uc *UserController) ForgotPassword(ctx *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.service.ForgotPassword(req.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (uc *UserController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Param("token")
	if err := uc.service.VerifyEmail(token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
