package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabase creates a new database connection based on the provided config
func NewDatabase(config Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	switch config.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DBName), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			// Enable foreign key constraints for SQLite
			sqlDB, _ := db.DB()
			_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
		}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// SetupDatabase initializes the database with migrations and indexes
func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Org{},
		&OrgMember{},
		&OAuthApplication{},
		&OAuthAuthorizationCode{},
		&OAuthAccessToken{},
		&OAuthRefreshToken{},
		&OAuthScope{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate models: %w", err)
	}
	err = db.SetupJoinTable(&User{}, "Orgs", &OrgMember{})
	if err != nil {
		return fmt.Errorf("failed to setup join table: %w", err)
	}

	// Initialize default OAuth scopes
	err = initializeDefaultOAuthScopes(db)
	if err != nil {
		return fmt.Errorf("failed to initialize default OAuth scopes: %w", err)
	}

	return nil
}

// InitializeDefaultAdmin creates a default admin user if none exists
func InitializeDefaultAdmin(db *gorm.DB) error {
	// Check if default admin creation is disabled
	if getEnv("DISABLE_DEFAULT_ADMIN", "false") == "true" {
		fmt.Println("Default admin creation disabled via DISABLE_DEFAULT_ADMIN environment variable")
		return nil
	}

	// Check if any admin user already exists
	var adminCount int64
	if err := db.Model(&User{}).Where("role = ?", UserRoleAdmin).Count(&adminCount).Error; err != nil {
		return fmt.Errorf("failed to check for existing admin users: %w", err)
	}

	// If admin user already exists, skip initialization
	if adminCount > 0 {
		fmt.Println("Admin user already exists, skipping default admin creation")
		return nil
	}

	// Create default admin user
	defaultAdmin := &User{
		Username: getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		Email:    getEnv("DEFAULT_ADMIN_EMAIL", "admin@example.com"),
		Role:     UserRoleAdmin,
	}

	// Set default password
	defaultPassword := getEnv("DEFAULT_ADMIN_PASSWORD", "admin123")
	if err := defaultAdmin.SetPassword(defaultPassword); err != nil {
		return fmt.Errorf("failed to set default admin password: %w", err)
	}

	// Create the admin user in a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Create the admin user
		if err := tx.Create(defaultAdmin).Error; err != nil {
			return fmt.Errorf("failed to create default admin user: %w", err)
		}

		// Create a default organization for the admin
		defaultOrg := &Org{
			Name: "Admin Organization",
			Slug: "admin-org",
		}

		if err := tx.Create(defaultOrg).Error; err != nil {
			return fmt.Errorf("failed to create default admin organization: %w", err)
		}

		// Add the admin as owner of the organization
		orgMember := &OrgMember{
			UserID: defaultAdmin.ID,
			OrgID:  defaultOrg.ID,
			Role:   OrgMemberRoleOwner,
		}

		if err := tx.Create(orgMember).Error; err != nil {
			return fmt.Errorf("failed to create admin organization membership: %w", err)
		}

		fmt.Printf("Default admin user created successfully:\n")
		fmt.Printf("  Username: %s\n", defaultAdmin.Username)
		fmt.Printf("  Email: %s\n", defaultAdmin.Email)
		fmt.Printf("  Password: %s\n", defaultPassword)
		fmt.Printf("  Please change the default password after first login!\n")

		return nil
	})
}

// GetDatabaseConfig returns database configuration from environment variables
func GetDatabaseConfig() Config {
	return Config{
		Driver:   getEnv("DB_DRIVER", "sqlite"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", ".db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// initializeDefaultOAuthScopes creates default OAuth scopes if they don't exist
func initializeDefaultOAuthScopes(db *gorm.DB) error {
	// Check if any scopes already exist
	var scopeCount int64
	if err := db.Model(&OAuthScope{}).Count(&scopeCount).Error; err != nil {
		return fmt.Errorf("failed to check for existing OAuth scopes: %w", err)
	}

	// If scopes already exist, skip initialization
	if scopeCount > 0 {
		fmt.Println("OAuth scopes already exist, skipping default scope creation")
		return nil
	}

	// Define default OAuth scopes
	defaultScopes := []OAuthScope{
		{
			Name:        "read",
			Description: "Read access to basic user information",
			Default:     true,
		},
		{
			Name:        "write",
			Description: "Write access to user data",
			Default:     false,
		},
		{
			Name:        "admin",
			Description: "Administrative access to all resources",
			Default:     false,
		},
		{
			Name:        "profile",
			Description: "Access to user profile information",
			Default:     true,
		},
		{
			Name:        "organizations",
			Description: "Access to user organization information",
			Default:     false,
		},
	}

	// Create the scopes in a transaction
	return db.Transaction(func(tx *gorm.DB) error {
		for _, scope := range defaultScopes {
			if err := tx.Create(&scope).Error; err != nil {
				return fmt.Errorf("failed to create OAuth scope '%s': %w", scope.Name, err)
			}
		}

		fmt.Printf("Default OAuth scopes created successfully:\n")
		for _, scope := range defaultScopes {
			fmt.Printf("  - %s: %s (default: %t)\n", scope.Name, scope.Description, scope.Default)
		}

		return nil
	})
}
