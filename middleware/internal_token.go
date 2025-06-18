package middleware

import (
	"crypto/subtle"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// InternalTokenManager manages internal API token authentication
type InternalTokenManager struct {
	token string
}

// NewInternalTokenManager creates a new internal token manager
func NewInternalTokenManager() *InternalTokenManager {
	// Get internal token from environment variable
	token := os.Getenv("MINIAUTH_INTERNAL_TOKEN")
	if token == "" {
		// Use a default token for development (should be changed in production)
		token = "miniauth-internal-default-token-change-in-production"
	}

	return &InternalTokenManager{
		token: token,
	}
}

// ValidateInternalToken validates the internal token from request
func (itm *InternalTokenManager) ValidateInternalToken(c echo.Context) bool {
	// Check Authorization header first
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		// Support "Bearer <token>" format
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := authHeader[7:] // Remove "Bearer " prefix
			return subtle.ConstantTimeCompare([]byte(token), []byte(itm.token)) == 1
		}
		// Support "Internal <token>" format
		if strings.HasPrefix(authHeader, "Internal ") {
			token := authHeader[9:] // Remove "Internal " prefix
			return subtle.ConstantTimeCompare([]byte(token), []byte(itm.token)) == 1
		}
	}

	// Check X-Internal-Token header
	internalToken := c.Request().Header.Get("X-Internal-Token")
	if internalToken != "" {
		return subtle.ConstantTimeCompare([]byte(internalToken), []byte(itm.token)) == 1
	}

	// Check internal_token query parameter
	queryToken := c.QueryParam("internal_token")
	if queryToken != "" {
		return subtle.ConstantTimeCompare([]byte(queryToken), []byte(itm.token)) == 1
	}

	return false
}

// RequireInternalToken middleware requires internal token authentication
func (itm *InternalTokenManager) RequireInternalToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if itm.ValidateInternalToken(c) {
			return next(c)
		}
		return c.JSON(401, map[string]string{
			"error":             "unauthorized",
			"error_description": "Valid internal token required",
		})
	}
}

// RequireInternalTokenOrAdmin middleware requires either internal token or admin session
func (itm *InternalTokenManager) RequireInternalTokenOrAdmin(sessionManager *SessionManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First check internal token
			if itm.ValidateInternalToken(c) {
				// Set a virtual admin user for internal token requests
				c.Set("internal_auth", true)
				c.Set("virtual_admin_id", uint(1)) // Use admin user ID 1 as default
				return next(c)
			}

			// Fall back to admin session authentication
			user, err := sessionManager.GetCurrentUser(c)
			if err != nil {
				return c.JSON(401, map[string]string{
					"error":             "unauthorized",
					"error_description": "Valid internal token or admin session required",
				})
			}

			if user.Role != "admin" {
				return c.JSON(403, map[string]string{
					"error":             "forbidden",
					"error_description": "Admin privileges required",
				})
			}

			c.Set("internal_auth", false)
			c.Set("admin_user", user)
			return next(c)
		}
	}
}
