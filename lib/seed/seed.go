package seed

import (
	"fmt"

	"github.com/horlakz/api.secretariat_repository/lib/database"
	"github.com/horlakz/api.secretariat_repository/lib/helper"
	"github.com/horlakz/api.secretariat_repository/model"
	user_service "github.com/horlakz/api.secretariat_repository/service/user"
)

type SeederInterface interface {
	Seed()
}

type seeder struct {
	dbConn database.DatabaseInterface
}

func NewSeeder(dbConn database.DatabaseInterface) SeederInterface {
	return &seeder{dbConn: dbConn}
}

func (s *seeder) Seed() {
	s.SeedAdmin()
}

func (s *seeder) SeedAdmin() {
	hashing := helper.NewHashing()
	adminEmail := "admin@email.com"
	testEmail := "test@email.com"
	hashedPassword, err := hashing.HashPassword("password")
	if err != nil {
		fmt.Println("Failed to hash password:", err)
		return
	}

	adminExists := s.dbConn.Connection().Where("email = ?", adminEmail).First(&model.User{}).RowsAffected > 0
	testUserExists := s.dbConn.Connection().Where("email = ?", testEmail).First(&model.User{}).RowsAffected > 0

	// Create Admin User
	adminUser := model.User{
		FirstName:       "Repo",
		LastName:        "Admin",
		Email:           adminEmail,
		IsEmailVerified: true,
		Password:        hashedPassword,
		Role:            user_service.UserRoleAdmin,
	}

	adminUser.Prepare()

	// Create test user
	testUser := model.User{
		FirstName:       "Test",
		LastName:        "User",
		Email:           "test@email.com",
		IsEmailVerified: true,
		Password:        hashedPassword,
		Role:            user_service.UserRoleUser,
	}

	testUser.Prepare()

	if adminExists {
		fmt.Println("Admin already exists in the database. Skipping seeding...")
	} else {
		if err := s.dbConn.Connection().Create(&adminUser).Error; err != nil {
			fmt.Println("Failed to create admin user:", err)
		} else {
			fmt.Println("Admin user created successfully.")
		}
	}

	if testUserExists {
		fmt.Println("Test user already exists in the database. Skipping seeding...")
	} else {
		if err := s.dbConn.Connection().Create(&testUser).Error; err != nil {
			fmt.Println("Failed to create test user:", err)
		} else {
			fmt.Println("Test user created successfully.")
		}
	}

}
