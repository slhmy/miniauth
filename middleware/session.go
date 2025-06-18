package middleware

import (
	"encoding/gob"
	"miniauth/database"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// SessionManager manages user sessions
type SessionManager struct {
	store sessions.Store
}

// SessionData represents the data stored in a session
type SessionData struct {
	UserID   uint              `json:"user_id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Role     database.UserRole `json:"role"`
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionKey string) *SessionManager {
	// Register SessionData type for gob encoding
	gob.Register(SessionData{})
	gob.Register(database.UserRole(""))

	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return &SessionManager{
		store: store,
	}
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(ctx echo.Context, user *database.User) error {
	session, err := sm.store.Get(ctx.Request(), "user-session")
	if err != nil {
		return err
	}

	sessionData := SessionData{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	session.Values["user"] = sessionData
	return session.Save(ctx.Request(), ctx.Response())
}

// GetSession retrieves session data for the current user
func (sm *SessionManager) GetSession(ctx echo.Context) (*SessionData, error) {
	session, err := sm.store.Get(ctx.Request(), "user-session")
	if err != nil {
		return nil, err
	}

	userData, ok := session.Values["user"]
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "No session found")
	}

	sessionData, ok := userData.(SessionData)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid session data")
	}

	return &sessionData, nil
}

// UpdateSession updates session data for the current user
func (sm *SessionManager) UpdateSession(ctx echo.Context, user *database.User) error {
	session, err := sm.store.Get(ctx.Request(), "user-session")
	if err != nil {
		return err
	}

	sessionData := SessionData{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	session.Values["user"] = sessionData
	return session.Save(ctx.Request(), ctx.Response())
}

// DestroySession destroys the current user session
func (sm *SessionManager) DestroySession(ctx echo.Context) error {
	session, err := sm.store.Get(ctx.Request(), "user-session")
	if err != nil {
		return err
	}

	session.Values["user"] = nil
	session.Options.MaxAge = -1
	return session.Save(ctx.Request(), ctx.Response())
}

// RequireAuth is a middleware that requires authentication
func (sm *SessionManager) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionData, err := sm.GetSession(ctx)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authentication required",
			})
		}

		// Store session data in context for use in handlers
		ctx.Set("currentUser", sessionData)
		return next(ctx)
	}
}

// RequireAdmin is a middleware that requires admin role
func (sm *SessionManager) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionData, err := sm.GetSession(ctx)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authentication required",
			})
		}

		if sessionData.Role != database.UserRoleAdmin {
			return ctx.JSON(http.StatusForbidden, map[string]string{
				"error": "Admin access required",
			})
		}

		// Store session data in context for use in handlers
		ctx.Set("currentUser", sessionData)
		return next(ctx)
	}
}

// GetCurrentUser retrieves the current user from session and returns a User object
func (sm *SessionManager) GetCurrentUser(ctx echo.Context) (*database.User, error) {
	sessionData, err := sm.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	// Convert session data to User object
	user := &database.User{
		Username: sessionData.Username,
		Email:    sessionData.Email,
		Role:     sessionData.Role,
	}
	user.ID = sessionData.UserID

	return user, nil
}
