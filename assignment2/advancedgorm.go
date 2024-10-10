package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// User model
type User struct {
	ID      uint    `gorm:"primaryKey"`
	Name    string  `gorm:"unique;not null"`
	Age     int     `gorm:"not null"`
	Profile Profile `gorm:"foreignKey:UserID"`
}

// Profile model (one-to-one relationship with User)
type Profile struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint `gorm:"unique;not null"`
	Bio               string
	ProfilePictureURL string
}

// Connect to MySQL using GORM
func ConnectGORM() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/gocon?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}
	fmt.Println("Connected to MySQL using GORM!")
}

// AutoMigrate models
func AutoMigrateModels() {
	db.AutoMigrate(&User{}, &Profile{})
	fmt.Println("Database migrated successfully!")
}

// Insert user and profile with transaction
func InsertUserWithProfile() {
	user := User{Name: "John Doe", Age: 28, Profile: Profile{Bio: "Software Engineer", ProfilePictureURL: "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_960_720.png"}}
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println("Failed to insert user:", result.Error)
		return
	}
	fmt.Println("User and profile inserted successfully!")
}

// Query users with profiles (eager loading)
func QueryUsersWithProfile() {
	var users []User
	db.Preload("Profile").Find(&users)

	for _, user := range users {
		fmt.Printf("User: %s, Bio: %s, Profile Picture: %s\n", user.Name, user.Profile.Bio, user.Profile.ProfilePictureURL)
	}
}

// Update user's profile
func UpdateUserProfile(userID uint, newBio string) {
	result := db.Model(&Profile{}).Where("user_id = ?", userID).Update("Bio", newBio)
	if result.Error != nil {
		fmt.Println("Failed to update profile:", result.Error)
		return
	}
	fmt.Println("Profile updated successfully!")
}

// Delete user and associated profile
func DeleteUserWithProfile(userID uint) {
	result := db.Delete(&User{}, userID)
	if result.Error != nil {
		fmt.Println("Failed to delete user:", result.Error)
		return
	}
	fmt.Println("User and profile deleted successfully!")
}

func main() {
	ConnectGORM()
	AutoMigrateModels()
	InsertUserWithProfile()
	QueryUsersWithProfile()

	fmt.Println("Updating user's profile...")
	UpdateUserProfile(1, "Senior Developer")

	fmt.Println("Deleting user...")
	DeleteUserWithProfile(1)
}
