package controllers

import (
	"net/http"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service ports.UserService
}

func NewUserController(service ports.UserService) *UserController {
	return &UserController{service: service}
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

// func (uc *UserController) VerifyEmail(ctx *gin.Context) {
// 	token := ctx.Param("token")
// 	if err := uc.service.VerifyEmail(token); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
// }
