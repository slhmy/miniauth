package service

import (
	"errors"
	"miniauth/database"

	"gorm.io/gorm"
)

type OrgService struct {
	db *gorm.DB
}

func NewOrgService(db *gorm.DB) *OrgService {
	return &OrgService{db: db}
}

// CreateOrg creates a new organization
func (s *OrgService) CreateOrg(org *database.Org) error {
	if org.Name == "" {
		return errors.New("organization name is required")
	}
	if org.Slug == "" {
		return errors.New("organization slug is required")
	}

	return s.db.Create(org).Error
}

// GetOrgByID retrieves an organization by ID with optional preloading
func (s *OrgService) GetOrgByID(id uint, preload ...string) (*database.Org, error) {
	var org database.Org
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.First(&org, id).Error
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// GetOrgBySlug retrieves an organization by slug with optional preloading
func (s *OrgService) GetOrgBySlug(slug string, preload ...string) (*database.Org, error) {
	var org database.Org
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Where("slug = ?", slug).First(&org).Error
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// UpdateOrg updates an existing organization
func (s *OrgService) UpdateOrg(org *database.Org) error {
	if org.ID == 0 {
		return errors.New("organization ID is required")
	}

	return s.db.Save(org).Error
}

// DeleteOrg soft deletes an organization by ID
func (s *OrgService) DeleteOrg(id uint) error {
	// Start a transaction to delete organization and its members
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete all org members first
		if err := tx.Where("org_id = ?", id).Delete(&database.OrgMember{}).Error; err != nil {
			return err
		}

		// Delete the organization
		return tx.Delete(&database.Org{}, id).Error
	})
}

// ListOrgs retrieves organizations with pagination and optional filtering
func (s *OrgService) ListOrgs(offset, limit int, preload ...string) ([]database.Org, error) {
	var orgs []database.Org
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Offset(offset).Limit(limit).Find(&orgs).Error
	return orgs, err
}

// CountOrgs returns the total number of organizations
func (s *OrgService) CountOrgs() (int64, error) {
	var count int64
	err := s.db.Model(&database.Org{}).Count(&count).Error
	return count, err
}

// GetOrgMembers retrieves all members of an organization
func (s *OrgService) GetOrgMembers(orgID uint) ([]database.OrgMember, error) {
	var members []database.OrgMember
	err := s.db.Where("org_id = ?", orgID).Find(&members).Error
	return members, err
}

// GetOrgMemberByUserID retrieves a specific member of an organization
func (s *OrgService) GetOrgMemberByUserID(orgID, userID uint) (*database.OrgMember, error) {
	var member database.OrgMember
	err := s.db.Where("org_id = ? AND user_id = ?", orgID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// IsUserInOrg checks if a user is a member of an organization
func (s *OrgService) IsUserInOrg(userID, orgID uint) (bool, error) {
	var count int64
	err := s.db.Model(&database.OrgMember{}).
		Where("user_id = ? AND org_id = ?", userID, orgID).
		Count(&count).Error
	return count > 0, err
}

// GetUserOrgs retrieves all organizations that a user belongs to
func (s *OrgService) GetUserOrgs(userID uint, preload ...string) ([]database.Org, error) {
	var orgs []database.Org
	query := s.db.Model(&database.Org{})

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Joins("JOIN org_members ON orgs.id = org_members.org_id").
		Where("org_members.user_id = ?", userID).
		Find(&orgs).Error

	return orgs, err
}

// GetOrgsByRole retrieves organizations where user has a specific role
func (s *OrgService) GetOrgsByRole(userID uint, role database.OrgMemberRole, preload ...string) ([]database.Org, error) {
	var orgs []database.Org
	query := s.db.Model(&database.Org{})

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Joins("JOIN org_members ON orgs.id = org_members.org_id").
		Where("org_members.user_id = ? AND org_members.role = ?", userID, role).
		Find(&orgs).Error

	return orgs, err
}

// OrgWithRole represents an organization with the user's role in it
type OrgWithRole struct {
	database.Org
	Role database.OrgMemberRole `json:"role"`
}

// GetUserOrgsWithRoles retrieves all organizations that a user belongs to with their roles
func (s *OrgService) GetUserOrgsWithRoles(userID uint) ([]OrgWithRole, error) {
	var results []OrgWithRole

	err := s.db.Table("orgs").
		Select("orgs.*, org_members.role").
		Joins("JOIN org_members ON orgs.id = org_members.org_id").
		Where("org_members.user_id = ? AND org_members.deleted_at IS NULL", userID).
		Scan(&results).Error

	return results, err
}
