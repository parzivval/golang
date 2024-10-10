package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User struct for GORM
type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null;unique"`
	Age  int    `json:"age" gorm:"not null"`
}

// DB variables
var db *gorm.DB
var sqlDB *sql.DB

// Connect to the database using GORM
func connectDatabase() {
	dsn := "root:password@tcp(127.0.0.1:3306)/gocon?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Get the generic database object sql.DB to use it in raw queries
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}
}

// Auto migrate the User model
func migrate() {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User table migrated.")
}

// Handler to fetch all users (using GORM)
func getUsersGORM(c *gin.Context) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Handler to create a user (using GORM)
func createUserGORM(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Handler to update a user (using GORM)
func updateUserGORM(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Model(&User{}).Where("id = ?", id).Updates(user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Handler to delete a user (using GORM)
func deleteUserGORM(c *gin.Context) {
	id := c.Param("id")
	result := db.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// Handler to fetch all users (using direct SQL)
func getUsersSQL(c *gin.Context) {
	rows, err := sqlDB.Query("SELECT id, name, age FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// Handler to create a user (using direct SQL)
func createUserSQL(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := sqlDB.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func main() {
	// Connect to the database
	connectDatabase()
	fmt.Println("Connected to the database.")

	// Migrate the User model
	migrate()

	// Set up Gin router
	router := gin.Default()

	// Routes for GORM
	router.GET("/gorm/users", getUsersGORM)
	router.POST("/gorm/user", createUserGORM)
	router.PUT("/gorm/user/:id", updateUserGORM)
	router.DELETE("/gorm/user/:id", deleteUserGORM)

	// Routes for direct SQL
	router.GET("/sql/users", getUsersSQL)
	router.POST("/sql/user", createUserSQL)

	// Start the server
	router.Run(":8080")
}
