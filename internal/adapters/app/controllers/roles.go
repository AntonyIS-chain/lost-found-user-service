package controllers

import (
	"net/http"
	"strconv"

	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-user-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService ports.RoleService
	userService ports.UserService
}

func NewRoleController(roleService ports.RoleService, userService ports.UserService) *RoleController {
	return &RoleController{roleService: roleService, userService: userService}
}

func (rc *RoleController) CreateRole(ctx *gin.Context) {
	var role domain.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdRole, err := rc.roleService.CreateRole(role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdRole)
}

func (rc *RoleController) GetRoleByName(ctx *gin.Context) {
	roleName := ctx.Param("role_name")
	role, err := rc.roleService.GetRoleByName(roleName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

func (rc *RoleController) ListRoles(ctx *gin.Context) {
	roles, err := rc.roleService.ListRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

func (rc *RoleController) UpdateRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var updates domain.Role
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = rc.roleService.UpdateRole(roleID, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func (rc *RoleController) DeleteRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("role_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = rc.roleService.DeleteRole(roleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
func (rc *RoleController) AssignRoleToUser(ctx *gin.Context) {
	var req struct {
		UserID   string `json:"user_id"`
		RoleName string `json:"role_name"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate required fields
	if req.UserID == "" || req.RoleName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Role name is required"})
		return
	}

	// Check if user exists
	user, err := rc.userService.GetUserByID(req.UserID)
	if err != nil || user.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if role exists
	role, err := rc.roleService.GetRoleByName(req.RoleName)
	if err != nil || role.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Assign role to user
	err = rc.roleService.AssignRoleToUser(req.UserID, req.RoleName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user"})
		return
	}

	// Update user record with new role details
	user.RoleID = role.ID
	user.RoleName = role.Name
	err = rc.userService.UpdateUser(user.ID, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role details"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role assigned to user successfully"})
}
