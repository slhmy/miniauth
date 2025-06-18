package service

import "gorm.io/gorm"

// ServiceManager holds all service instances
type ServiceManager struct {
	User  *UserService
	Org   *OrgService
	OAuth *OAuthService
}

// NewServiceManager creates a new service manager with all services initialized
func NewServiceManager(db *gorm.DB) *ServiceManager {
	return &ServiceManager{
		User:  NewUserService(db),
		Org:   NewOrgService(db),
		OAuth: NewOAuthService(db),
	}
}
