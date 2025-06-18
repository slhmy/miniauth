package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"miniauth/database"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OAuthService handles OAuth 2.0 operations
type OAuthService struct {
	db *gorm.DB
}

// NewOAuthService creates a new OAuth service instance
func NewOAuthService(db *gorm.DB) *OAuthService {
	return &OAuthService{db: db}
}

// OAuth request/response structures
type AuthorizeRequest struct {
	ResponseType        string `json:"response_type" validate:"required"`
	ClientID            string `json:"client_id" validate:"required"`
	RedirectURI         string `json:"redirect_uri" validate:"required"`
	Scope               string `json:"scope"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret"`
	CodeVerifier string `json:"code_verifier"`
	RefreshToken string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

type ApplicationCreateRequest struct {
	Name         string   `json:"name" validate:"required"`
	RedirectURIs []string `json:"redirect_uris" validate:"required"`
	Scopes       []string `json:"scopes"`
	Description  string   `json:"description"`
	Website      string   `json:"website"`
	Trusted      bool     `json:"trusted"`
}

// InternalApplicationCreateRequest allows specifying custom client_id and secret for internal use
type InternalApplicationCreateRequest struct {
	Name         string   `json:"name" validate:"required"`
	RedirectURIs []string `json:"redirect_uris" validate:"required"`
	Scopes       []string `json:"scopes"`
	Description  string   `json:"description"`
	Website      string   `json:"website"`
	Trusted      bool     `json:"trusted"`
	ClientID     string   `json:"client_id" validate:"required"`     // Custom client ID
	ClientSecret string   `json:"client_secret" validate:"required"` // Custom client secret
}

type ApplicationResponse struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Website      string   `json:"website"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       []string `json:"scopes"`
	Trusted      bool     `json:"trusted"`
	Active       bool     `json:"active"`
	CreatedBy    string   `json:"created_by"`
	CreatedAt    string   `json:"created_at"`
}

// CreateApplication creates a new OAuth application (admin only)
func (s *OAuthService) CreateApplication(adminUserID uint, req ApplicationCreateRequest) (*ApplicationResponse, error) {
	// Generate client ID and secret
	clientID := uuid.New().String()
	clientSecret := s.generateClientSecret()

	// Convert redirect URIs and scopes to JSON strings
	redirectURIsJSON, err := json.Marshal(req.RedirectURIs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal redirect URIs: %w", err)
	}

	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"read"}
	}

	app := &database.OAuthApplication{
		Name:         req.Name,
		Description:  req.Description,
		Website:      req.Website,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURIs: string(redirectURIsJSON),
		Scopes:       strings.Join(scopes, " "),
		CreatedByID:  adminUserID,
		Trusted:      req.Trusted,
		Active:       true,
	}

	if err := s.db.Preload("CreatedBy").Create(app).Error; err != nil {
		return nil, fmt.Errorf("failed to create OAuth application: %w", err)
	}

	return &ApplicationResponse{
		ID:           app.ID,
		Name:         app.Name,
		Description:  app.Description,
		Website:      app.Website,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		RedirectURIs: req.RedirectURIs,
		Scopes:       scopes,
		Trusted:      app.Trusted,
		Active:       app.Active,
		CreatedBy:    app.CreatedBy.Username,
		CreatedAt:    app.CreatedAt.Format(time.RFC3339),
	}, nil
}

// CreateApplicationWithCustomCredentials creates a new OAuth application with specified client_id and secret (internal use)
func (s *OAuthService) CreateApplicationWithCustomCredentials(adminUserID uint, req InternalApplicationCreateRequest) (*ApplicationResponse, error) {
	// Check if client_id already exists
	var existingApp database.OAuthApplication
	if err := s.db.Where("client_id = ?", req.ClientID).First(&existingApp).Error; err == nil {
		return nil, fmt.Errorf("client_id already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check client_id uniqueness: %w", err)
	}

	// Convert redirect URIs and scopes to JSON strings
	redirectURIsJSON, err := json.Marshal(req.RedirectURIs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal redirect URIs: %w", err)
	}

	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"read"}
	}

	app := &database.OAuthApplication{
		Name:         req.Name,
		Description:  req.Description,
		Website:      req.Website,
		ClientID:     req.ClientID,     // Use provided client ID
		ClientSecret: req.ClientSecret, // Use provided client secret
		RedirectURIs: string(redirectURIsJSON),
		Scopes:       strings.Join(scopes, " "),
		CreatedByID:  adminUserID,
		Trusted:      req.Trusted,
		Active:       true,
	}

	if err := s.db.Preload("CreatedBy").Create(app).Error; err != nil {
		return nil, fmt.Errorf("failed to create OAuth application: %w", err)
	}

	return &ApplicationResponse{
		ID:           app.ID,
		Name:         app.Name,
		Description:  app.Description,
		Website:      app.Website,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		RedirectURIs: req.RedirectURIs,
		Scopes:       scopes,
		Trusted:      app.Trusted,
		Active:       app.Active,
		CreatedBy:    app.CreatedBy.Username,
		CreatedAt:    app.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetAllApplications retrieves all OAuth applications (admin only)
func (s *OAuthService) GetAllApplications() ([]*ApplicationResponse, error) {
	var apps []database.OAuthApplication
	if err := s.db.Preload("CreatedBy").Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	result := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		var redirectURIs []string
		if err := json.Unmarshal([]byte(app.RedirectURIs), &redirectURIs); err != nil {
			redirectURIs = []string{}
		}

		scopes := strings.Split(app.Scopes, " ")
		if len(scopes) == 1 && scopes[0] == "" {
			scopes = []string{}
		}

		result[i] = &ApplicationResponse{
			ID:           app.ID,
			Name:         app.Name,
			Description:  app.Description,
			Website:      app.Website,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			RedirectURIs: redirectURIs,
			Scopes:       scopes,
			Trusted:      app.Trusted,
			Active:       app.Active,
			CreatedBy:    app.CreatedBy.Username,
			CreatedAt:    app.CreatedAt.Format(time.RFC3339),
		}
	}

	return result, nil
}

// GetActiveApplications retrieves all active OAuth applications
func (s *OAuthService) GetActiveApplications() ([]*ApplicationResponse, error) {
	var apps []database.OAuthApplication
	if err := s.db.Preload("CreatedBy").Where("active = ?", true).Find(&apps).Error; err != nil {
		return nil, fmt.Errorf("failed to get active applications: %w", err)
	}

	result := make([]*ApplicationResponse, len(apps))
	for i, app := range apps {
		var redirectURIs []string
		if err := json.Unmarshal([]byte(app.RedirectURIs), &redirectURIs); err != nil {
			redirectURIs = []string{}
		}

		scopes := strings.Split(app.Scopes, " ")
		if len(scopes) == 1 && scopes[0] == "" {
			scopes = []string{}
		}

		result[i] = &ApplicationResponse{
			ID:           app.ID,
			Name:         app.Name,
			Description:  app.Description,
			Website:      app.Website,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			RedirectURIs: redirectURIs,
			Scopes:       scopes,
			Trusted:      app.Trusted,
			Active:       app.Active,
			CreatedBy:    app.CreatedBy.Username,
			CreatedAt:    app.CreatedAt.Format(time.RFC3339),
		}
	}

	return result, nil
}

// ValidateAuthorizationRequest validates an OAuth authorization request
func (s *OAuthService) ValidateAuthorizationRequest(req AuthorizeRequest) (*database.OAuthApplication, error) {
	// Check if response_type is supported
	if req.ResponseType != "code" {
		return nil, fmt.Errorf("unsupported response_type: %s", req.ResponseType)
	}

	// Get the OAuth application
	var app database.OAuthApplication
	if err := s.db.Where("client_id = ? AND active = ?", req.ClientID, true).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or inactive client_id")
		}
		return nil, fmt.Errorf("failed to get OAuth application: %w", err)
	}

	// Validate redirect URI
	var redirectURIs []string
	if err := json.Unmarshal([]byte(app.RedirectURIs), &redirectURIs); err != nil {
		return nil, fmt.Errorf("failed to parse redirect URIs: %w", err)
	}

	isValidRedirectURI := false
	for _, uri := range redirectURIs {
		if uri == req.RedirectURI {
			isValidRedirectURI = true
			break
		}
	}

	if !isValidRedirectURI {
		return nil, fmt.Errorf("invalid redirect_uri")
	}

	// Validate PKCE parameters if present
	if req.CodeChallenge != "" {
		if req.CodeChallengeMethod != "" && req.CodeChallengeMethod != "plain" && req.CodeChallengeMethod != "S256" {
			return nil, fmt.Errorf("unsupported code_challenge_method: %s", req.CodeChallengeMethod)
		}
	}

	return &app, nil
}

// CreateAuthorizationCode creates an authorization code for a user
func (s *OAuthService) CreateAuthorizationCode(userID uint, app *database.OAuthApplication, req AuthorizeRequest) (string, error) {
	// Generate authorization code
	code := s.generateAuthorizationCode()

	// Determine scopes (intersection of requested and allowed)
	requestedScopes := strings.Split(req.Scope, " ")
	if len(requestedScopes) == 1 && requestedScopes[0] == "" {
		requestedScopes = []string{"read"} // Default scope
	}

	allowedScopes := strings.Split(app.Scopes, " ")
	grantedScopes := s.intersectScopes(requestedScopes, allowedScopes)

	authCode := &database.OAuthAuthorizationCode{
		Code:                code,
		ClientID:            app.ClientID,
		UserID:              userID,
		RedirectURI:         req.RedirectURI,
		Scopes:              strings.Join(grantedScopes, " "),
		ExpiresAt:           time.Now().Add(10 * time.Minute), // Authorization codes expire in 10 minutes
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	if err := s.db.Create(authCode).Error; err != nil {
		return "", fmt.Errorf("failed to create authorization code: %w", err)
	}

	return code, nil
}

// ExchangeCodeForToken exchanges an authorization code for access token
func (s *OAuthService) ExchangeCodeForToken(req TokenRequest) (*TokenResponse, error) {
	// Validate grant type
	if req.GrantType != "authorization_code" {
		return nil, fmt.Errorf("unsupported grant_type: %s", req.GrantType)
	}

	// Get authorization code
	var authCode database.OAuthAuthorizationCode
	if err := s.db.Where("code = ? AND used = false", req.Code).First(&authCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid or expired authorization code")
		}
		return nil, fmt.Errorf("failed to get authorization code: %w", err)
	}

	// Check if code is expired
	if time.Now().After(authCode.ExpiresAt) {
		return nil, fmt.Errorf("authorization code expired")
	}

	// Validate client
	if authCode.ClientID != req.ClientID {
		return nil, fmt.Errorf("client_id mismatch")
	}

	// Validate redirect URI
	if authCode.RedirectURI != req.RedirectURI {
		return nil, fmt.Errorf("redirect_uri mismatch")
	}

	// Validate PKCE if present
	if authCode.CodeChallenge != "" {
		if req.CodeVerifier == "" {
			return nil, fmt.Errorf("code_verifier required for PKCE")
		}

		if !s.validatePKCE(authCode.CodeChallenge, authCode.CodeChallengeMethod, req.CodeVerifier) {
			return nil, fmt.Errorf("invalid code_verifier")
		}
	}

	// Get client secret for validation (optional for public clients)
	if req.ClientSecret != "" {
		var app database.OAuthApplication
		if err := s.db.Where("client_id = ?", req.ClientID).First(&app).Error; err != nil {
			return nil, fmt.Errorf("failed to get OAuth application: %w", err)
		}

		if app.ClientSecret != req.ClientSecret {
			return nil, fmt.Errorf("invalid client_secret")
		}
	}

	// Mark authorization code as used
	if err := s.db.Model(&authCode).Update("used", true).Error; err != nil {
		return nil, fmt.Errorf("failed to mark authorization code as used: %w", err)
	}

	// Generate access token and refresh token
	accessToken := s.generateAccessToken()
	refreshToken := s.generateRefreshToken()

	// Store access token
	accessTokenRecord := &database.OAuthAccessToken{
		Token:     accessToken,
		ClientID:  req.ClientID,
		UserID:    authCode.UserID,
		Scopes:    authCode.Scopes,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Access tokens expire in 1 hour
	}

	if err := s.db.Create(accessTokenRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Store refresh token
	refreshTokenRecord := &database.OAuthRefreshToken{
		Token:       refreshToken,
		ClientID:    req.ClientID,
		UserID:      authCode.UserID,
		AccessToken: accessToken,
		ExpiresAt:   time.Now().Add(24 * 30 * time.Hour), // Refresh tokens expire in 30 days
	}

	if err := s.db.Create(refreshTokenRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
		RefreshToken: refreshToken,
		Scope:        authCode.Scopes,
	}, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OAuthService) RefreshAccessToken(req TokenRequest) (*TokenResponse, error) {
	// Validate grant type
	if req.GrantType != "refresh_token" {
		return nil, fmt.Errorf("unsupported grant_type: %s", req.GrantType)
	}

	// Get refresh token
	var refreshTokenRecord database.OAuthRefreshToken
	if err := s.db.Where("token = ? AND revoked = false", req.RefreshToken).First(&refreshTokenRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if refresh token is expired
	if time.Now().After(refreshTokenRecord.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Validate client
	if refreshTokenRecord.ClientID != req.ClientID {
		return nil, fmt.Errorf("client_id mismatch")
	}

	// Get the associated access token to get scopes
	var oldAccessToken database.OAuthAccessToken
	if err := s.db.Where("token = ?", refreshTokenRecord.AccessToken).First(&oldAccessToken).Error; err != nil {
		return nil, fmt.Errorf("failed to get associated access token: %w", err)
	}

	// Revoke old access token
	if err := s.db.Model(&oldAccessToken).Update("revoked", true).Error; err != nil {
		return nil, fmt.Errorf("failed to revoke old access token: %w", err)
	}

	// Generate new access token
	newAccessToken := s.generateAccessToken()

	// Store new access token
	newAccessTokenRecord := &database.OAuthAccessToken{
		Token:     newAccessToken,
		ClientID:  req.ClientID,
		UserID:    refreshTokenRecord.UserID,
		Scopes:    oldAccessToken.Scopes,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Access tokens expire in 1 hour
	}

	if err := s.db.Create(newAccessTokenRecord).Error; err != nil {
		return nil, fmt.Errorf("failed to create new access token: %w", err)
	}

	// Update refresh token's associated access token
	if err := s.db.Model(&refreshTokenRecord).Update("access_token", newAccessToken).Error; err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  newAccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour in seconds
		RefreshToken: req.RefreshToken,
		Scope:        oldAccessToken.Scopes,
	}, nil
}

// ValidateAccessToken validates an access token and returns user info
func (s *OAuthService) ValidateAccessToken(token string) (*database.User, []string, error) {
	var accessToken database.OAuthAccessToken
	if err := s.db.Preload("User").Where("token = ? AND revoked = false", token).First(&accessToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("invalid access token")
		}
		return nil, nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(accessToken.ExpiresAt) {
		return nil, nil, fmt.Errorf("access token expired")
	}

	scopes := strings.Split(accessToken.Scopes, " ")
	return &accessToken.User, scopes, nil
}

// Helper functions

func (s *OAuthService) generateClientSecret() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (s *OAuthService) generateAuthorizationCode() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (s *OAuthService) generateAccessToken() string {
	return uuid.New().String()
}

func (s *OAuthService) generateRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (s *OAuthService) intersectScopes(requested, allowed []string) []string {
	allowedMap := make(map[string]bool)
	for _, scope := range allowed {
		allowedMap[scope] = true
	}

	var result []string
	for _, scope := range requested {
		if allowedMap[scope] {
			result = append(result, scope)
		}
	}

	return result
}

func (s *OAuthService) validatePKCE(codeChallenge, codeChallengeMethod, codeVerifier string) bool {
	if codeChallengeMethod == "" || codeChallengeMethod == "plain" {
		return codeChallenge == codeVerifier
	}

	if codeChallengeMethod == "S256" {
		hash := sha256.Sum256([]byte(codeVerifier))
		computed := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])
		return codeChallenge == computed
	}

	return false
}

// JWT token generation for OAuth (optional, for JWT-based tokens)
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

func (s *OAuthService) GenerateJWTToken(userID uint, clientID string, scopes []string, expiresAt time.Time) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		ClientID: clientID,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "miniauth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-jwt-secret-key")) // Should be from environment
}

// DeleteApplication deletes an OAuth application and all associated tokens (admin only)
func (s *OAuthService) DeleteApplication(appID uint) error {
	// Check if application exists
	var app database.OAuthApplication
	if err := s.db.Where("id = ?", appID).First(&app).Error; err != nil {
		return err
	}

	// Delete all associated tokens and codes in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete authorization codes
		if err := tx.Where("client_id = ?", app.ClientID).Delete(&database.OAuthAuthorizationCode{}).Error; err != nil {
			return err
		}

		// Delete access tokens
		if err := tx.Where("client_id = ?", app.ClientID).Delete(&database.OAuthAccessToken{}).Error; err != nil {
			return err
		}

		// Delete refresh tokens
		if err := tx.Where("client_id = ?", app.ClientID).Delete(&database.OAuthRefreshToken{}).Error; err != nil {
			return err
		}

		// Delete the application
		if err := tx.Delete(&app).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateApplication updates an OAuth application (admin only)
func (s *OAuthService) UpdateApplication(appID uint, req ApplicationCreateRequest) (*ApplicationResponse, error) {
	var app database.OAuthApplication
	if err := s.db.Preload("CreatedBy").Where("id = ?", appID).First(&app).Error; err != nil {
		return nil, err
	}

	// Convert redirect URIs to JSON string
	redirectURIsJSON, err := json.Marshal(req.RedirectURIs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal redirect URIs: %w", err)
	}

	scopes := req.Scopes
	if len(scopes) == 0 {
		scopes = []string{"read"}
	}

	// Update application fields
	app.Name = req.Name
	app.Description = req.Description
	app.Website = req.Website
	app.RedirectURIs = string(redirectURIsJSON)
	app.Scopes = strings.Join(scopes, " ")
	app.Trusted = req.Trusted

	if err := s.db.Save(&app).Error; err != nil {
		return nil, fmt.Errorf("failed to update OAuth application: %w", err)
	}

	return &ApplicationResponse{
		ID:           app.ID,
		Name:         app.Name,
		Description:  app.Description,
		Website:      app.Website,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		RedirectURIs: req.RedirectURIs,
		Scopes:       scopes,
		Trusted:      app.Trusted,
		Active:       app.Active,
		CreatedBy:    app.CreatedBy.Username,
		CreatedAt:    app.CreatedAt.Format(time.RFC3339),
	}, nil
}

// ToggleApplicationStatus toggles the active status of an OAuth application
func (s *OAuthService) ToggleApplicationStatus(appID uint) error {
	var app database.OAuthApplication
	if err := s.db.Where("id = ?", appID).First(&app).Error; err != nil {
		return err
	}

	app.Active = !app.Active
	return s.db.Save(&app).Error
}

// ToggleApplicationTrustedStatus toggles the trusted status of an OAuth application
func (s *OAuthService) ToggleApplicationTrustedStatus(appID uint) error {
	var app database.OAuthApplication
	if err := s.db.Where("id = ?", appID).First(&app).Error; err != nil {
		return err
	}

	app.Trusted = !app.Trusted
	return s.db.Save(&app).Error
}
