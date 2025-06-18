package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	gorm.Model
	Username     string
	Email        string   `gorm:"uniqueIndex;not null"`
	PasswordHash string   `gorm:"not null"`
	Orgs         []Org    `gorm:"many2many:user_orgs;"`
	Role         UserRole `gorm:"not null;default:'user'"`
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type Org struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Slug    string `gorm:"uniqueIndex;not null"`
	Members []OrgMember
}

type OrgMemberRole string

const (
	OrgMemberRoleOwner  OrgMemberRole = "owner"
	OrgMemberRoleAdmin  OrgMemberRole = "admin"
	OrgMemberRoleMember OrgMemberRole = "member"
	OrgMemberRoleGuest  OrgMemberRole = "guest"
)

type OrgMember struct {
	OrgID     uint `gorm:"primaryKey"`
	UserID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
	Role      OrgMemberRole `gorm:"not null;default:'member'"`
}

// OAuth Application represents a registered OAuth client application (system-level)
type OAuthApplication struct {
	gorm.Model
	Name         string `gorm:"not null"`
	ClientID     string `gorm:"uniqueIndex;not null"`
	ClientSecret string `gorm:"not null"`
	RedirectURIs string `gorm:"type:text"`      // JSON array of allowed redirect URIs
	Scopes       string `gorm:"default:'read'"` // Space-separated scopes
	Description  string `gorm:"type:text"`      // Application description
	Website      string // Application website URL
	CreatedByID  uint   `gorm:"not null"` // Admin who created the application
	CreatedBy    User   `gorm:"foreignKey:CreatedByID"`
	Trusted      bool   `gorm:"default:false"` // Whether this app can skip user consent
	Active       bool   `gorm:"default:true"`  // Whether this app is active
}

// OAuth Authorization Code
type OAuthAuthorizationCode struct {
	gorm.Model
	Code                string    `gorm:"uniqueIndex;not null"`
	ClientID            string    `gorm:"not null"`
	UserID              uint      `gorm:"not null"`
	User                User      `gorm:"foreignKey:UserID"`
	RedirectURI         string    `gorm:"not null"`
	Scopes              string    // Space-separated scopes
	ExpiresAt           time.Time `gorm:"not null"`
	Used                bool      `gorm:"default:false"`
	CodeChallenge       string    // For PKCE
	CodeChallengeMethod string    // For PKCE (plain or S256)
}

// OAuth Access Token
type OAuthAccessToken struct {
	gorm.Model
	Token     string    `gorm:"uniqueIndex;not null"`
	ClientID  string    `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Scopes    string    // Space-separated scopes
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
}

// OAuth Refresh Token
type OAuthRefreshToken struct {
	gorm.Model
	Token       string    `gorm:"uniqueIndex;not null"`
	ClientID    string    `gorm:"not null"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	AccessToken string    `gorm:"not null"` // Associated access token
	ExpiresAt   time.Time `gorm:"not null"`
	Revoked     bool      `gorm:"default:false"`
}

// OAuth Scope represents available OAuth scopes
type OAuthScope struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Default     bool `gorm:"default:false"` // Whether this scope is granted by default
}
