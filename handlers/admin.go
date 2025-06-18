package handlers

import (
	"fmt"
	"miniauth/database"
	"miniauth/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Admin API structures
type AdminListUsersResponse struct {
	Users []AdminUserInfo `json:"users"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

type AdminUserInfo struct {
	ID        uint              `json:"id"`
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	OrgCount  int               `json:"org_count"`
	Role      database.UserRole `json:"role"`
}

type AdminUpdateUserRequest struct {
	Username string             `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string             `json:"email" validate:"omitempty,email"`
	Role     *database.UserRole `json:"role" validate:"omitempty,oneof=admin user"`
}

type AdminCreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type AdminUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

// AdminListUsers lists all users with pagination
//
//	@Summary		List all users (Admin)
//	@Description	Get a paginated list of all users in the system
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (default: 1)"
//	@Param			size	query		int	false	"Page size (default: 10)"
//	@Success		200		{object}	AdminListUsersResponse
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/admin/users [get]
func AdminListUsers(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Parse pagination parameters
	page := 1
	size := 10

	if p := ctx.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if s := ctx.QueryParam("size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 && parsed <= 100 {
			size = parsed
		}
	}

	offset := (page - 1) * size

	// Get users and total count
	users, err := serviceManager.User.ListUsers(offset, size)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get users",
		})
	}

	total, err := serviceManager.User.CountUsers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to count users",
		})
	}

	// Convert to admin user info
	var adminUsers []AdminUserInfo
	for _, user := range users {
		// Get organization count for each user
		orgs, _ := serviceManager.Org.GetUserOrgs(user.ID)

		adminUsers = append(adminUsers, AdminUserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			OrgCount:  len(orgs),
			Role:      user.Role,
		})
	}

	return ctx.JSON(http.StatusOK, AdminListUsersResponse{
		Users: adminUsers,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// AdminGetUser gets a specific user by ID
//
//	@Summary		Get user by ID (Admin)
//	@Description	Get detailed information about a specific user
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	GetUserResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/admin/users/{id} [get]
func AdminGetUser(ctx echo.Context) error {
	// Reuse the existing GetUser logic
	return GetUser(ctx)
}

// AdminCreateUser creates a new user (admin operation)
//
//	@Summary		Create user (Admin)
//	@Description	Create a new user account (admin operation)
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		AdminCreateUserRequest	true	"User creation request"
//	@Success		201		{object}	AdminUserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/admin/users [post]
func AdminCreateUser(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	var req AdminCreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Create user
	user := &database.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to set password",
		})
	}

	if err := serviceManager.User.CreateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, AdminUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User created successfully",
	})
}

// AdminUpdateUser updates an existing user
//
//	@Summary		Update user (Admin)
//	@Description	Update user information (admin operation)
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"User ID"
//	@Param			user	body		AdminUpdateUserRequest	true	"User update request"
//	@Success		200		{object}	AdminUserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/admin/users/{id} [put]
func AdminUpdateUser(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user ID from path
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Parse request
	var req AdminUpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Get existing user
	user, err := serviceManager.User.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := serviceManager.User.UpdateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, AdminUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User updated successfully",
	})
}

// AdminDeleteUser permanently deletes a user
//
//	@Summary		Delete user (Admin)
//	@Description	Permanently delete a user account and all related data (admin operation)
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/admin/users/{id} [delete]
func AdminDeleteUser(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user ID from path
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Check if user exists
	_, err := serviceManager.User.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Delete user
	if err := serviceManager.User.DeleteUser(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

// AdminResetUserPassword resets a user's password
//
//	@Summary		Reset user password (Admin)
//	@Description	Reset a user's password to a new value (admin operation)
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int					true	"User ID"
//	@Param			password	body		map[string]string	true	"New password"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/admin/users/{id}/reset-password [post]
func AdminResetUserPassword(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user ID from path
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Parse request
	var req struct {
		Password string `json:"password" validate:"required,min=6,max=100"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Check if user exists and reset password
	if err := serviceManager.User.UpdateUserPassword(id, req.Password); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password reset successfully",
	})
}

// AdminUpdateUserRoleRequest represents the request structure for updating user role
type AdminUpdateUserRoleRequest struct {
	Role database.UserRole `json:"role" validate:"required,oneof=admin user"`
}

// AdminUpdateUserRole updates a user's role
//
//	@Summary		Update user role (Admin)
//	@Description	Update a user's role (admin or user)
//	@Tags			admin
//	@Security		BasicAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"User ID"
//	@Param			role	body		AdminUpdateUserRoleRequest	true	"User role update request"
//	@Success		200		{object}	AdminUserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/admin/users/{id}/role [put]
func AdminUpdateUserRole(ctx echo.Context) error {
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user ID from path
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Parse request
	var req AdminUpdateUserRoleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Update user role using service
	if err := serviceManager.User.UpdateUserRole(id, req.Role); err != nil {
		if err.Error() == "user not found" {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Get updated user info
	user, err := serviceManager.User.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get updated user info",
		})
	}

	return ctx.JSON(http.StatusOK, AdminUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User role updated successfully",
	})
}
