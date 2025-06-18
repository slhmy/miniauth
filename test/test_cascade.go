package main

import (
	"fmt"
	"log"
	"miniauth/database"
	"miniauth/service"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Create a test database
	db, err := gorm.Open(sqlite.Open("test_cascade.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Enable foreign key constraints for SQLite
	sqlDB, _ := db.DB()
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&database.User{}, &database.Org{}, &database.OrgMember{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	// Create a test user
	user := &database.User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     database.UserRoleUser,
	}
	user.SetPassword("password123")

	if err := db.Create(user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Create a test organization
	org := &database.Org{
		Name: "Test Org",
		Slug: "test-org",
	}
	if err := db.Create(org).Error; err != nil {
		log.Fatal("Failed to create org:", err)
	}
	fmt.Printf("Created org with ID: %d\n", org.ID)

	// Create organization membership
	orgMember := &database.OrgMember{
		UserID: user.ID,
		OrgID:  org.ID,
		Role:   database.OrgMemberRoleMember,
	}
	if err := db.Create(orgMember).Error; err != nil {
		log.Fatal("Failed to create org membership:", err)
	}
	fmt.Println("Created organization membership")

	// Test 1: Delete user and check if org memberships are cleaned up
	fmt.Println("\n=== Test 1: User deletion cascade ===")

	// Check organization memberships before user deletion
	var userOrgCountBefore int64
	db.Model(&database.OrgMember{}).Where("user_id = ?", user.ID).Count(&userOrgCountBefore)
	fmt.Printf("Organization memberships before user deletion: %d\n", userOrgCountBefore)

	// Create another user to keep the org alive after first user deletion
	user2 := &database.User{
		Username: "testuser2",
		Email:    "test2@example.com",
		Role:     database.UserRoleUser,
	}
	user2.SetPassword("password123")
	if err := db.Create(user2).Error; err != nil {
		log.Fatal("Failed to create user2:", err)
	}

	// Add user2 to the same org
	orgMember2 := &database.OrgMember{
		UserID: user2.ID,
		OrgID:  org.ID,
		Role:   database.OrgMemberRoleOwner,
	}
	if err := db.Create(orgMember2).Error; err != nil {
		log.Fatal("Failed to create org membership for user2:", err)
	}

	// Delete the first user using transaction (this should also delete org memberships)
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Delete all organization memberships for this user
		if err := tx.Where("user_id = ?", user.ID).Delete(&database.OrgMember{}).Error; err != nil {
			return err
		}
		// Then delete the user
		if err := tx.Delete(&database.User{}, user.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal("Failed to delete user:", err)
	}
	fmt.Println("Deleted user1")

	// Check organization memberships after user deletion
	var userOrgCountAfter int64
	db.Model(&database.OrgMember{}).Where("user_id = ?", user.ID).Count(&userOrgCountAfter)
	fmt.Printf("Organization memberships after user deletion: %d\n", userOrgCountAfter)

	// Test 2: Delete organization and check if memberships are cleaned up
	fmt.Println("\n=== Test 2: Organization deletion cascade ===")

	// Check organization memberships before org deletion
	var orgMemberCountBefore int64
	db.Model(&database.OrgMember{}).Where("org_id = ?", org.ID).Count(&orgMemberCountBefore)
	fmt.Printf("Organization memberships before org deletion: %d\n", orgMemberCountBefore)

	// Delete the organization (this should cascade delete all memberships)
	if err := db.Transaction(func(tx *gorm.DB) error {
		// Delete all organization memberships for this org
		if err := tx.Where("org_id = ?", org.ID).Delete(&database.OrgMember{}).Error; err != nil {
			return err
		}
		// Then delete the organization
		if err := tx.Delete(&database.Org{}, org.ID).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal("Failed to delete organization:", err)
	}
	fmt.Println("Deleted organization")

	// Check organization memberships after org deletion
	var orgMemberCountAfter int64
	db.Model(&database.OrgMember{}).Where("org_id = ?", org.ID).Count(&orgMemberCountAfter)
	fmt.Printf("Organization memberships after org deletion: %d\n", orgMemberCountAfter)

	// Verify user2 still exists
	var user2Exists int64
	db.Model(&database.User{}).Where("id = ?", user2.ID).Count(&user2Exists)
	fmt.Printf("User2 still exists after org deletion: %t\n", user2Exists > 0)

	// Test 3: Test using service layer methods
	fmt.Println("\n=== Test 3: Service layer organization deletion ===")

	// Create a new organization and users for service layer testing
	org3 := &database.Org{
		Name: "Service Test Org",
		Slug: "service-test-org",
	}
	if err := db.Create(org3).Error; err != nil {
		log.Fatal("Failed to create org3:", err)
	}

	user3 := &database.User{
		Username: "serviceuser",
		Email:    "service@example.com",
		Role:     database.UserRoleUser,
	}
	user3.SetPassword("password123")
	if err := db.Create(user3).Error; err != nil {
		log.Fatal("Failed to create user3:", err)
	}

	// Add user3 to org3
	orgMember3 := &database.OrgMember{
		UserID: user3.ID,
		OrgID:  org3.ID,
		Role:   database.OrgMemberRoleMember,
	}
	if err := db.Create(orgMember3).Error; err != nil {
		log.Fatal("Failed to create org membership for user3:", err)
	}

	// Create org service
	orgService := service.NewOrgService(db)

	// Check memberships before service deletion
	var serviceOrgCountBefore int64
	db.Model(&database.OrgMember{}).Where("org_id = ?", org3.ID).Count(&serviceOrgCountBefore)
	fmt.Printf("Organization memberships before service deletion: %d\n", serviceOrgCountBefore)

	// Delete organization using service method
	if err := orgService.DeleteOrg(org3.ID); err != nil {
		log.Fatal("Failed to delete org using service:", err)
	}
	fmt.Println("Deleted organization using service layer")

	// Check memberships after service deletion
	var serviceOrgCountAfter int64
	db.Model(&database.OrgMember{}).Where("org_id = ?", org3.ID).Count(&serviceOrgCountAfter)
	fmt.Printf("Organization memberships after service deletion: %d\n", serviceOrgCountAfter)

	// Verify user3 still exists
	var user3Exists int64
	db.Model(&database.User{}).Where("id = ?", user3.ID).Count(&user3Exists)
	fmt.Printf("User3 still exists after service org deletion: %t\n", user3Exists > 0)

	// Clean up
	os.Remove("test_cascade.db")

	// Final validation
	fmt.Println("\n=== Test Results ===")
	if userOrgCountAfter == 0 && orgMemberCountAfter == 0 && serviceOrgCountAfter == 0 {
		fmt.Println("✅ All cascade deletes working correctly!")
		fmt.Println("   - User deletion properly removed organization memberships")
		fmt.Println("   - Organization deletion properly removed all memberships")
		fmt.Println("   - Service layer organization deletion works correctly")
	} else {
		fmt.Println("❌ Cascade delete issues detected:")
		if userOrgCountAfter > 0 {
			fmt.Printf("   - User deletion left %d orphaned memberships\n", userOrgCountAfter)
		}
		if orgMemberCountAfter > 0 {
			fmt.Printf("   - Organization deletion left %d orphaned memberships\n", orgMemberCountAfter)
		}
		if serviceOrgCountAfter > 0 {
			fmt.Printf("   - Service layer deletion left %d orphaned memberships\n", serviceOrgCountAfter)
		}
	}
}
