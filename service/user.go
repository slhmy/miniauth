package service

import (
	"errors"
	"miniauth/database"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user and a corresponding organization
func (s *UserService) CreateUser(user *database.User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Username == "" {
		return errors.New("username is required")
	}

	// Start a transaction to ensure both user and org are created together
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create the user first
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Create a same-name organization for the user
		org := &database.Org{
			Name: user.Username,
			Slug: user.Username,
		}

		if err := tx.Create(org).Error; err != nil {
			return err
		}

		// Add the user as owner of the organization
		orgMember := &database.OrgMember{
			UserID: user.ID,
			OrgID:  org.ID,
			Role:   database.OrgMemberRoleOwner,
		}

		if err := tx.Create(orgMember).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetUserByID retrieves a user by ID with optional preloading
func (s *UserService) GetUserByID(id uint, preload ...string) (*database.User, error) {
	var user database.User
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email with optional preloading
func (s *UserService) GetUserByEmail(email string, preload ...string) (*database.User, error) {
	var user database.User
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by username with optional preloading
func (s *UserService) GetUserByUsername(username string, preload ...string) (*database.User, error) {
	var user database.User
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *database.User) error {
	if user.ID == 0 {
		return errors.New("user ID is required")
	}

	return s.db.Save(user).Error
}

// DeleteUser hard deletes a user by ID and cleans up related records
func (s *UserService) DeleteUser(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First, check if user exists
		var user database.User
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		// Delete all organization memberships for this user (hard delete)
		result := tx.Unscoped().Where("user_id = ?", id).Delete(&database.OrgMember{})
		if result.Error != nil {
			return result.Error
		}

		// Log the number of organization memberships deleted
		if result.RowsAffected > 0 {
			// Note: In a real application, you might want to use a proper logger here
			// For now, this comment serves as a placeholder for logging
		}

		// Then hard delete the user (permanent deletion)
		if err := tx.Unscoped().Delete(&database.User{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// ListUsers retrieves users with pagination and optional filtering
func (s *UserService) ListUsers(offset, limit int, preload ...string) ([]database.User, error) {
	var users []database.User
	query := s.db

	for _, p := range preload {
		query = query.Preload(p)
	}

	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// CountUsers returns the total number of users
func (s *UserService) CountUsers() (int64, error) {
	var count int64
	err := s.db.Model(&database.User{}).Count(&count).Error
	return count, err
}

// AddUserToOrg adds a user to an organization with specified role
func (s *UserService) AddUserToOrg(userID, orgID uint, role database.OrgMemberRole) error {
	orgMember := database.OrgMember{
		UserID: userID,
		OrgID:  orgID,
		Role:   role,
	}

	return s.db.Create(&orgMember).Error
}

// RemoveUserFromOrg removes a user from an organization
func (s *UserService) RemoveUserFromOrg(userID, orgID uint) error {
	return s.db.Where("user_id = ? AND org_id = ?", userID, orgID).Delete(&database.OrgMember{}).Error
}

// UpdateUserOrgRole updates a user's role in an organization
func (s *UserService) UpdateUserOrgRole(userID, orgID uint, role database.OrgMemberRole) error {
	return s.db.Model(&database.OrgMember{}).
		Where("user_id = ? AND org_id = ?", userID, orgID).
		Update("role", role).Error
}

// AuthenticateUser verifies user credentials and returns the user if valid
func (s *UserService) AuthenticateUser(email, password string) (*database.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// UpdateUserPassword updates a user's password
func (s *UserService) UpdateUserPassword(userID uint, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	return s.db.Save(user).Error
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(userID uint, role database.UserRole) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Role = role
	return s.db.Save(user).Error
}

// GetUserOrgCount returns the number of organizations a user belongs to
func (s *UserService) GetUserOrgCount(userID uint) (int64, error) {
	var count int64
	err := s.db.Model(&database.OrgMember{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetUserWithOrgCount returns user info along with organization count
func (s *UserService) GetUserWithOrgCount(userID uint) (*database.User, int64, error) {
	var user database.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		return nil, 0, err
	}

	count, err := s.GetUserOrgCount(userID)
	if err != nil {
		return &user, 0, err
	}

	return &user, count, nil
}
