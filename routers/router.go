package routers

import (
	"miniauth/handlers"
	"miniauth/middleware"
	"miniauth/service"
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(e *echo.Echo, serviceManager *service.ServiceManager) {
	// Create session manager with a secret key (should be from environment in production)
	sessionManager := middleware.NewSessionManager("your-secret-key-change-in-production")

	// Create internal token manager
	internalTokenManager := middleware.NewInternalTokenManager()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Add service manager and session manager to context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("serviceManager", serviceManager)
			c.Set("sessionManager", sessionManager)
			return next(c)
		}
	})

	// Serve static files from website/dist directory
	e.Static("/", "website/dist")

	// Handle client-side routing - serve index.html for unmatched routes
	e.RouteNotFound("/*", func(c echo.Context) error {
		// Only serve index.html for non-API routes
		if len(c.Request().URL.Path) >= 4 && c.Request().URL.Path[:4] == "/api" {
			return echo.NewHTTPError(http.StatusNotFound, "API endpoint not found")
		}
		return c.File("website/dist/index.html")
	})

	// API routes
	api := e.Group("/api")

	// Auth routes (no authentication required)
	api.POST("/login", handlers.LoginUser)
	api.POST("/logout", handlers.LogoutUser)

	// User routes
	users := api.Group("/users")
	users.POST("", handlers.CreateUser)
	users.GET("/:id", handlers.GetUser)

	// Protected routes (authentication required)
	protected := api.Group("/me")
	protected.Use(sessionManager.RequireAuth)
	protected.GET("", handlers.GetCurrentUser)
	protected.PUT("/change-password", handlers.ChangePassword)
	protected.PUT("/profile", handlers.UpdateProfile)

	// Admin routes (admin authentication required)
	admin := api.Group("/admin")
	admin.Use(sessionManager.RequireAdmin)
	adminUsers := admin.Group("/users")
	adminUsers.GET("", handlers.AdminListUsers)
	adminUsers.POST("", handlers.AdminCreateUser)
	adminUsers.GET("/:id", handlers.AdminGetUser)
	adminUsers.PUT("/:id", handlers.AdminUpdateUser)
	adminUsers.DELETE("/:id", handlers.AdminDeleteUser)
	adminUsers.POST("/:id/reset-password", handlers.AdminResetUserPassword)
	adminUsers.PUT("/:id/role", handlers.AdminUpdateUserRole)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// OAuth routes
	oauth := api.Group("/oauth")
	oauth.GET("/authorize", handlers.OAuthAuthorize)
	oauth.POST("/authorize", handlers.OAuthAuthorizeDecision)
	oauth.POST("/token", handlers.OAuthToken)
	oauth.GET("/userinfo", handlers.OAuthUserInfo)

	// OAuth application management (Admin only)
	adminOAuth := admin.Group("/oauth")
	adminOAuthApps := adminOAuth.Group("/applications")
	adminOAuthApps.GET("", handlers.AdminListOAuthApplications)
	adminOAuthApps.POST("", handlers.AdminCreateOAuthApplication)
	adminOAuthApps.PUT("/:id", handlers.AdminUpdateOAuthApplication)
	adminOAuthApps.DELETE("/:id", handlers.AdminDeleteOAuthApplication)
	adminOAuthApps.POST("/:id/toggle", handlers.AdminToggleOAuthApplicationStatus)
	adminOAuthApps.POST("/:id/toggle-trusted", handlers.AdminToggleOAuthApplicationTrustedStatus)

	// Internal OAuth application management (Internal Token or Admin) - allows custom client_id and secret
	// Note: Create directly under /api to avoid inheriting admin middleware
	internalOAuth := api.Group("/admin/oauth/internal")
	internalOAuthApps := internalOAuth.Group("/applications")

	// Use internal token or admin authentication for internal APIs
	internalOAuthApps.Use(internalTokenManager.RequireInternalTokenOrAdmin(sessionManager))
	internalOAuthApps.POST("", handlers.AdminInternalCreateOAuthApplication)
	internalOAuthApps.POST("/batch", handlers.AdminInternalBatchCreateOAuthApplications)
}
