package handlers

import (
	"fmt"
	"miniauth/database"
	"miniauth/middleware"
	"miniauth/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type CreateUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

// CreateUser creates a new user with a corresponding organization
//
//	@Summary		Create a new user
//	@Description	Create a new user account with username, email, and password. Also creates a same-name organization with the user as owner.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		CreateUserRequest	true	"User creation request"
//	@Success		201		{object}	CreateUserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users [post]
func CreateUser(ctx echo.Context) error {
	// Get service manager from context
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Parse request body
	var req CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request using echo's validator
	if err := ctx.Validate(&req); err != nil {
		return err // This will return the validation error automatically formatted by CustomValidator
	}

	// Create user
	user := &database.User{
		Username: req.Username,
		Email:    req.Email,
	}

	// Set password
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

	// Return success response
	response := CreateUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User and organization created successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

type GetUserResponse struct {
	ID            uint               `json:"id"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Role          database.UserRole  `json:"role"`
	Organizations []OrganizationInfo `json:"organizations"`
}

type OrganizationInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

// GetUser retrieves user information by ID
//
//	@Summary		Get user by ID
//	@Description	Get user information including associated organizations and roles
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	GetUserResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/users/{id} [get]
func GetUser(ctx echo.Context) error {
	// Get service manager from context
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user ID from URL parameter
	userID := ctx.Param("id")
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}

	// Convert to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Get user
	user, err := serviceManager.User.GetUserByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// Get user's organizations
	orgs, err := serviceManager.Org.GetUserOrgs(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user organizations",
		})
	}

	// Get user's roles in each organization
	var organizations []OrganizationInfo
	for _, org := range orgs {
		member, err := serviceManager.Org.GetOrgMemberByUserID(org.ID, user.ID)
		role := "unknown"
		if err == nil {
			role = string(member.Role)
		}

		organizations = append(organizations, OrganizationInfo{
			ID:   org.ID,
			Name: org.Name,
			Slug: org.Slug,
			Role: role,
		})
	}

	// Return response
	response := GetUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		Organizations: organizations,
	}

	return ctx.JSON(http.StatusOK, response)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=100"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}

// LoginUser authenticates a user with email and password
//
//	@Summary		User login
//	@Description	Authenticate user with email and password and create a session
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		LoginRequest	true	"Login credentials"
//	@Success		200			{object}	LoginResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Router			/login [post]
func LoginUser(ctx echo.Context) error {
	// Get service manager from context
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Parse request body
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request using echo's validator
	if err := ctx.Validate(&req); err != nil {
		return err // This will return the validation error automatically formatted by CustomValidator
	}

	// Authenticate user
	user, err := serviceManager.User.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Create session
	sessionManager := ctx.Get("sessionManager").(*middleware.SessionManager)
	if err := sessionManager.CreateSession(ctx, user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create session",
		})
	}

	// Return success response
	response := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "Login successful",
	}

	return ctx.JSON(http.StatusOK, response)
}

// LogoutUser logs out the current user by destroying the session
//
//	@Summary		User logout
//	@Description	Log out the current user and destroy the session
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/logout [post]
func LogoutUser(ctx echo.Context) error {
	// Get session manager from context
	sessionManager := ctx.Get("sessionManager").(*middleware.SessionManager)

	// Destroy session
	if err := sessionManager.DestroySession(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to logout",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

// GetCurrentUser returns the current authenticated user's information
//
//	@Summary		Get current user
//	@Description	Get the currently authenticated user's information from session
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	GetUserResponse
//	@Failure		401	{object}	map[string]string
//	@Router			/me [get]
func GetCurrentUser(ctx echo.Context) error {
	// Get current user from session (set by middleware)
	currentUser := ctx.Get("currentUser").(*middleware.SessionData)

	// Get service manager to fetch full user details
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Get user's organizations
	orgs, err := serviceManager.Org.GetUserOrgs(currentUser.UserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user organizations",
		})
	}

	// Get user's roles in each organization
	var organizations []OrganizationInfo
	for _, org := range orgs {
		member, err := serviceManager.Org.GetOrgMemberByUserID(org.ID, currentUser.UserID)
		role := "unknown"
		if err == nil {
			role = string(member.Role)
		}

		organizations = append(organizations, OrganizationInfo{
			ID:   org.ID,
			Name: org.Name,
			Slug: org.Slug,
			Role: role,
		})
	}

	// Return response
	response := GetUserResponse{
		ID:            currentUser.UserID,
		Username:      currentUser.Username,
		Email:         currentUser.Email,
		Role:          currentUser.Role,
		Organizations: organizations,
	}

	return ctx.JSON(http.StatusOK, response)
}

// ChangePassword changes the current user's password
//
//	@Summary		Change user password
//	@Description	Change the current user's password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ChangePasswordRequest	true	"Change password request"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/me/change-password [put]
func ChangePassword(ctx echo.Context) error {
	// Get current user from session
	currentUser := ctx.Get("currentUser").(*middleware.SessionData)
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Parse request body
	var req ChangePasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Get user from database to verify current password
	user, err := serviceManager.User.GetUserByID(currentUser.UserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user",
		})
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Current password is incorrect",
		})
	}

	// Set new password
	if err := user.SetPassword(req.NewPassword); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to set new password",
		})
	}

	// Update user in database
	if err := serviceManager.User.UpdateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update password",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

// UpdateProfile updates the current user's profile information
//
//	@Summary		Update user profile
//	@Description	Update the current user's profile information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		UpdateProfileRequest	true	"Update profile request"
//	@Success		200		{object}	GetUserResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/me/profile [put]
func UpdateProfile(ctx echo.Context) error {
	// Get current user from session
	currentUser := ctx.Get("currentUser").(*middleware.SessionData)
	serviceManager := ctx.Get("serviceManager").(*service.ServiceManager)

	// Parse request body
	var req UpdateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := ctx.Validate(&req); err != nil {
		return err
	}

	// Get user from database
	user, err := serviceManager.User.GetUserByID(currentUser.UserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user",
		})
	}

	// Update user information
	user.Username = req.Username

	// Update user in database
	if err := serviceManager.User.UpdateUser(user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update profile",
		})
	}

	// Update session with new user information
	sessionManager := ctx.Get("sessionManager").(*middleware.SessionManager)
	if err := sessionManager.UpdateSession(ctx, user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update session",
		})
	}

	// Get user's organizations for response
	orgs, err := serviceManager.Org.GetUserOrgs(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user organizations",
		})
	}

	// Get user's roles in each organization
	var organizations []OrganizationInfo
	for _, org := range orgs {
		member, err := serviceManager.Org.GetOrgMemberByUserID(org.ID, user.ID)
		role := "unknown"
		if err == nil {
			role = string(member.Role)
		}

		organizations = append(organizations, OrganizationInfo{
			ID:   org.ID,
			Name: org.Name,
			Slug: org.Slug,
			Role: role,
		})
	}

	// Return updated user response
	response := GetUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		Organizations: organizations,
	}

	return ctx.JSON(http.StatusOK, response)
}
