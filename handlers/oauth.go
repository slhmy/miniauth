package handlers

import (
	"fmt"
	"miniauth/database"
	"miniauth/middleware"
	"miniauth/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// OAuth Authorization endpoint
//
//	@Summary		OAuth Authorization
//	@Description	Start OAuth authorization flow
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			response_type			query		string	true	"Response type (must be 'code')"
//	@Param			client_id				query		string	true	"OAuth client ID"
//	@Param			redirect_uri			query		string	true	"Redirect URI"
//	@Param			scope					query		string	false	"Requested scopes (space-separated)"
//	@Param			state					query		string	false	"State parameter for CSRF protection"
//	@Param			code_challenge			query		string	false	"PKCE code challenge"
//	@Param			code_challenge_method	query		string	false	"PKCE code challenge method"
//	@Success		302						{string}	string	"Redirect to authorization page or back to client"
//	@Failure		400						{object}	map[string]string
//	@Router			/oauth/authorize [get]
func OAuthAuthorize(c echo.Context) error {
	// Parse query parameters
	req := service.AuthorizeRequest{
		ResponseType:        c.QueryParam("response_type"),
		ClientID:            c.QueryParam("client_id"),
		RedirectURI:         c.QueryParam("redirect_uri"),
		Scope:               c.QueryParam("scope"),
		State:               c.QueryParam("state"),
		CodeChallenge:       c.QueryParam("code_challenge"),
		CodeChallengeMethod: c.QueryParam("code_challenge_method"),
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Validate the authorization request
	app, err := oauthService.ValidateAuthorizationRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Check if user is authenticated
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		// Redirect to login with OAuth parameters preserved
		loginURL := fmt.Sprintf("/login?oauth_redirect=true&client_id=%s&redirect_uri=%s&scope=%s&state=%s&response_type=%s",
			req.ClientID, req.RedirectURI, req.Scope, req.State, req.ResponseType)
		if req.CodeChallenge != "" {
			loginURL += fmt.Sprintf("&code_challenge=%s&code_challenge_method=%s", req.CodeChallenge, req.CodeChallengeMethod)
		}
		return c.Redirect(http.StatusFound, loginURL)
	}

	// For trusted applications, automatically grant authorization
	if app.Trusted {
		code, err := oauthService.CreateAuthorizationCode(user.ID, app, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error", "error_description": err.Error()})
		}

		// Redirect back to client with authorization code
		redirectURL := fmt.Sprintf("%s?code=%s", req.RedirectURI, code)
		if req.State != "" {
			redirectURL += fmt.Sprintf("&state=%s", req.State)
		}
		return c.Redirect(http.StatusFound, redirectURL)
	}

	// Return authorization page data (frontend will handle the consent UI)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"client_name":           app.Name,
		"client_id":             app.ClientID,
		"redirect_uri":          req.RedirectURI,
		"scope":                 req.Scope,
		"state":                 req.State,
		"response_type":         req.ResponseType,
		"code_challenge":        req.CodeChallenge,
		"code_challenge_method": req.CodeChallengeMethod,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// OAuth Authorization Decision endpoint
//
//	@Summary		OAuth Authorization Decision
//	@Description	Handle user's authorization decision
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		map[string]interface{}	true	"Authorization decision"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/oauth/authorize [post]
func OAuthAuthorizeDecision(c echo.Context) error {
	// Check if user is authenticated
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Parse request body
	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request"})
	}

	// Check authorization decision
	authorized, ok := requestBody["authorized"].(bool)
	if !ok || !authorized {
		// User denied authorization
		redirectURI, _ := requestBody["redirect_uri"].(string)
		state, _ := requestBody["state"].(string)

		redirectURL := fmt.Sprintf("%s?error=access_denied", redirectURI)
		if state != "" {
			redirectURL += fmt.Sprintf("&state=%s", state)
		}
		return c.JSON(http.StatusOK, map[string]string{"redirect_url": redirectURL})
	}

	// Build authorization request from the stored parameters
	req := service.AuthorizeRequest{
		ResponseType: requestBody["response_type"].(string),
		ClientID:     requestBody["client_id"].(string),
		RedirectURI:  requestBody["redirect_uri"].(string),
		Scope:        requestBody["scope"].(string),
		State:        requestBody["state"].(string),
	}

	if codeChallenge, ok := requestBody["code_challenge"].(string); ok {
		req.CodeChallenge = codeChallenge
	}
	if codeChallengeMethod, ok := requestBody["code_challenge_method"].(string); ok {
		req.CodeChallengeMethod = codeChallengeMethod
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Validate the authorization request again
	app, err := oauthService.ValidateAuthorizationRequest(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Create authorization code
	code, err := oauthService.CreateAuthorizationCode(user.ID, app, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error", "error_description": err.Error()})
	}

	// Return redirect URL with authorization code
	redirectURL := fmt.Sprintf("%s?code=%s", req.RedirectURI, code)
	if req.State != "" {
		redirectURL += fmt.Sprintf("&state=%s", req.State)
	}

	return c.JSON(http.StatusOK, map[string]string{"redirect_url": redirectURL})
}

// OAuth Token endpoint
//
//	@Summary		OAuth Token
//	@Description	Exchange authorization code or refresh token for access token
//	@Tags			OAuth
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			grant_type		formData	string	true	"Grant type (authorization_code or refresh_token)"
//	@Param			code			formData	string	false	"Authorization code (required for authorization_code grant)"
//	@Param			redirect_uri	formData	string	false	"Redirect URI (required for authorization_code grant)"
//	@Param			client_id		formData	string	true	"OAuth client ID"
//	@Param			client_secret	formData	string	false	"OAuth client secret"
//	@Param			code_verifier	formData	string	false	"PKCE code verifier"
//	@Param			refresh_token	formData	string	false	"Refresh token (required for refresh_token grant)"
//	@Success		200				{object}	service.TokenResponse
//	@Failure		400				{object}	map[string]string
//	@Router			/oauth/token [post]
func OAuthToken(c echo.Context) error {
	// Parse form data
	req := service.TokenRequest{
		GrantType:    c.FormValue("grant_type"),
		Code:         c.FormValue("code"),
		RedirectURI:  c.FormValue("redirect_uri"),
		ClientID:     c.FormValue("client_id"),
		ClientSecret: c.FormValue("client_secret"),
		CodeVerifier: c.FormValue("code_verifier"),
		RefreshToken: c.FormValue("refresh_token"),
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	var tokenResponse *service.TokenResponse
	var err error

	switch req.GrantType {
	case "authorization_code":
		tokenResponse, err = oauthService.ExchangeCodeForToken(req)
	case "refresh_token":
		tokenResponse, err = oauthService.RefreshAccessToken(req)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported_grant_type"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	return c.JSON(http.StatusOK, tokenResponse)
}

// OAuth User Info endpoint
//
//	@Summary		OAuth User Info
//	@Description	Get user information using access token
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]string
//	@Router			/oauth/userinfo [get]
func OAuthUserInfo(c echo.Context) error {
	// Get access token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid_token"})
	}

	// Extract token from "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid_token"})
	}
	token := authHeader[7:]

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Validate access token
	user, scopes, err := oauthService.ValidateAccessToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid_token"})
	}

	// Build user info response based on scopes
	userInfo := map[string]interface{}{
		"sub": user.ID, // Subject (user ID)
	}

	// Add information based on granted scopes
	for _, scope := range scopes {
		switch scope {
		case "profile":
			userInfo["username"] = user.Username
			userInfo["email"] = user.Email
		case "read":
			userInfo["id"] = user.ID
			userInfo["role"] = user.Role

			// Include organization information if user has organizations
			orgService := serviceManager.Org
			orgs, err := orgService.GetUserOrgsWithRoles(user.ID)
			if err == nil && len(orgs) > 0 {
				orgInfo := make([]map[string]interface{}, len(orgs))
				for i, org := range orgs {
					orgInfo[i] = map[string]interface{}{
						"id":   org.ID,
						"name": org.Name,
						"slug": org.Slug,
						"role": org.Role,
					}
				}
				userInfo["organizations"] = orgInfo
			}
		}
	}

	return c.JSON(http.StatusOK, userInfo)
}

// Create OAuth Application endpoint (Admin only)
//
//	@Summary		Create OAuth Application
//	@Description	Create a new OAuth application (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		service.ApplicationCreateRequest	true	"Application data"
//	@Success		201		{object}	service.ApplicationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Router			/admin/oauth/applications [post]
func AdminCreateOAuthApplication(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Parse request body
	var req service.ApplicationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request"})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Create application
	app, err := oauthService.CreateApplication(user.ID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error", "error_description": err.Error()})
	}

	return c.JSON(http.StatusCreated, app)
}

// List OAuth Applications endpoint (Admin only)
//
//	@Summary		List OAuth Applications
//	@Description	Get all OAuth applications (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		service.ApplicationResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Router			/admin/oauth/applications [get]
func AdminListOAuthApplications(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Get all applications
	apps, err := oauthService.GetAllApplications()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error", "error_description": err.Error()})
	}

	return c.JSON(http.StatusOK, apps)
}

// Delete OAuth Application endpoint (Admin only)
//
//	@Summary		Delete OAuth Application
//	@Description	Delete an OAuth application (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Application ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/oauth/applications/{id} [delete]
func AdminDeleteOAuthApplication(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Parse application ID
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_application_id"})
	}

	// Get service manager
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Delete the application
	if err := oauthService.DeleteApplication(uint(appID)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "application_not_found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Update OAuth Application endpoint (Admin only)
//
//	@Summary		Update OAuth Application
//	@Description	Update an OAuth application (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int									true	"Application ID"
//	@Param			request	body		service.ApplicationCreateRequest	true	"Application data"
//	@Success		200		{object}	service.ApplicationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/admin/oauth/applications/{id} [put]
func AdminUpdateOAuthApplication(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Parse application ID
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_application_id"})
	}

	// Parse request body
	var req service.ApplicationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request"})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Update application
	app, err := oauthService.UpdateApplication(uint(appID), req)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "application_not_found"})
	}

	return c.JSON(http.StatusOK, app)
}

// Toggle OAuth Application Status endpoint (Admin only)
//
//	@Summary		Toggle OAuth Application Status
//	@Description	Activate or deactivate an OAuth application (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Application ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/oauth/applications/{id}/toggle [post]
func AdminToggleOAuthApplicationStatus(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Parse application ID
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_application_id"})
	}

	// Get service manager
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Toggle application status
	if err := oauthService.ToggleApplicationStatus(uint(appID)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "application_not_found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Toggle OAuth Application Trusted Status endpoint (Admin only)
//
//	@Summary		Toggle OAuth Application Trusted Status
//	@Description	Toggle the trusted status of an OAuth application (admin only)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Application ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/admin/oauth/applications/{id}/toggle-trusted [post]
func AdminToggleOAuthApplicationTrustedStatus(c echo.Context) error {
	// Check if user is authenticated and is admin
	sessionManager := c.Get("sessionManager").(*middleware.SessionManager)
	user, err := sessionManager.GetCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin_required"})
	}

	// Parse application ID
	appIDStr := c.Param("id")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_application_id"})
	}

	// Get service manager
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Toggle application trusted status
	if err := oauthService.ToggleApplicationTrustedStatus(uint(appID)); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "application_not_found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Internal Create OAuth Application endpoint (Internal Token or Admin)
//
//	@Summary		Internal Create OAuth Application
//	@Description	Create an OAuth application with custom client_id and secret (internal API, requires internal token or admin session)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string										false	"Internal token (Bearer <token> or Internal <token>)"
//	@Param			X-Internal-Token	header		string										false	"Internal token"
//	@Param			internal_token		query		string										false	"Internal token"
//	@Param			request				body		service.InternalApplicationCreateRequest	true	"Application data with custom credentials"
//	@Success		201					{object}	service.ApplicationResponse
//	@Failure		400					{object}	map[string]string
//	@Failure		401					{object}	map[string]string
//	@Failure		403					{object}	map[string]string
//	@Failure		409					{object}	map[string]string
//	@Router			/admin/oauth/internal/applications [post]
func AdminInternalCreateOAuthApplication(c echo.Context) error {
	// Get admin user ID for creation
	var adminUserID uint

	// Check if this is internal token authentication
	if isInternal, ok := c.Get("internal_auth").(bool); ok && isInternal {
		// Use virtual admin ID for internal token requests
		adminUserID = c.Get("virtual_admin_id").(uint)
	} else {
		// Use actual admin user from session
		user := c.Get("admin_user").(*database.User)
		adminUserID = user.ID
	}

	// Parse request body
	var req service.InternalApplicationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request"})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request", "error_description": err.Error()})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	// Create application with custom credentials
	app, err := oauthService.CreateApplicationWithCustomCredentials(adminUserID, req)
	if err != nil {
		if err.Error() == "client_id already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": "client_id_exists", "error_description": "The specified client_id already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server_error", "error_description": err.Error()})
	}

	return c.JSON(http.StatusCreated, app)
}

// Batch Internal Create OAuth Applications endpoint (Internal Token or Admin)
//
//	@Summary		Batch Internal Create OAuth Applications
//	@Description	Create multiple OAuth applications with custom client_id and secret (internal API, requires internal token or admin session)
//	@Tags			OAuth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string										false	"Internal token (Bearer <token> or Internal <token>)"
//	@Param			X-Internal-Token	header		string										false	"Internal token"
//	@Param			internal_token		query		string										false	"Internal token"
//	@Param			request				body		[]service.InternalApplicationCreateRequest	true	"Array of application data with custom credentials"
//	@Success		200					{object}	map[string]interface{}
//	@Failure		400					{object}	map[string]string
//	@Failure		401					{object}	map[string]string
//	@Failure		403					{object}	map[string]string
//	@Router			/admin/oauth/internal/applications/batch [post]
func AdminInternalBatchCreateOAuthApplications(c echo.Context) error {
	// Get admin user ID for creation
	var adminUserID uint

	// Check if this is internal token authentication
	if isInternal, ok := c.Get("internal_auth").(bool); ok && isInternal {
		// Use virtual admin ID for internal token requests
		adminUserID = c.Get("virtual_admin_id").(uint)
	} else {
		// Use actual admin user from session
		user := c.Get("admin_user").(*database.User)
		adminUserID = user.ID
	}

	// Parse request body
	var requests []service.InternalApplicationCreateRequest
	if err := c.Bind(&requests); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid_request"})
	}

	if len(requests) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty_request", "error_description": "At least one application must be provided"})
	}

	// Get OAuth service
	serviceManager := c.Get("serviceManager").(*service.ServiceManager)
	oauthService := serviceManager.OAuth

	var successfulApps []*service.ApplicationResponse
	var errors []map[string]interface{}

	// Process each request
	for i, req := range requests {
		// Validate request
		if err := c.Validate(&req); err != nil {
			errors = append(errors, map[string]interface{}{
				"index":             i,
				"client_id":         req.ClientID,
				"error":             "validation_error",
				"error_description": err.Error(),
			})
			continue
		}

		// Create application
		app, err := oauthService.CreateApplicationWithCustomCredentials(adminUserID, req)
		if err != nil {
			errorType := "server_error"
			if err.Error() == "client_id already exists" {
				errorType = "client_id_exists"
			}
			errors = append(errors, map[string]interface{}{
				"index":             i,
				"client_id":         req.ClientID,
				"error":             errorType,
				"error_description": err.Error(),
			})
			continue
		}

		successfulApps = append(successfulApps, app)
	}

	response := map[string]interface{}{
		"success_count": len(successfulApps),
		"error_count":   len(errors),
		"successful":    successfulApps,
		"errors":        errors,
	}

	return c.JSON(http.StatusOK, response)
}
