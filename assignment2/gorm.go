package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
	Age  int    `gorm:"not null"`
}

// Connect to the database
func connectDatabase() (*gorm.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/gocon?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

// Auto migrate the User model
func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User table migrated.")
}

// Insert a user into the database
func insertUser(db *gorm.DB, name string, age int) {
	user := User{Name: name, Age: age}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	fmt.Printf("Inserted user: %s, Age: %d\n", name, age)
}

// Query all users from the database
func queryUsers(db *gorm.DB) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func main() {
	// Connect to the database
	db, err := connectDatabase()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database.")

	// Migrate the User model
	migrate(db)

	// Insert sample users
	insertUser(db, "Bakytzhan", 21)
	insertUser(db, "Damir", 22)

	// Query users
	queryUsers(db)
}
