package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connect to MySQL
func ConnectMySQL() {
	var err error
	connStr := "root:password@tcp(127.0.0.1:3306)/gocon"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.SetMaxOpenConns(25) // Connection pooling
	db.SetMaxIdleConns(25)
	fmt.Println("Connected to MySQL!")
}

// Create table with constraints
func CreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		age INT NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table created successfully!")
}

// Insert users within a transaction
func InsertUsers() {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction:", err)
	}

	_, err = tx.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Alice", 25)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to insert user Alice:", err)
	}

	_, err = tx.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "Bob", 30)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to insert user Bob:", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}
	fmt.Println("Users inserted successfully!")
}

// Query users with filtering and pagination
func QueryUsers(ageFilter string, page int) {
	limit := 2
	offset := (page - 1) * limit

	query := "SELECT * FROM users"
	if ageFilter != "" {
		query += " WHERE age = " + ageFilter
	}
	query += " LIMIT ? OFFSET ?"

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Fatal("Failed to query users:", err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var id int
		var name string
		var age int
		rows.Scan(&id, &name, &age)
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}
}

// Update user details by ID
func UpdateUser(id int, name string, age int) {
	_, err := db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", name, age, id)
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	fmt.Println("User updated successfully!")
}

// Delete user by ID
func DeleteUser(id int) {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal("Failed to delete user:", err)
	}
	fmt.Println("User deleted successfully!")
}

func main() {
	ConnectMySQL()
	CreateTable()
	InsertUsers()

	fmt.Println("Querying users with age filter 25 and page 1")
	QueryUsers("25", 1)

	fmt.Println("Updating user with ID 1")
	UpdateUser(1, "Alice Updated", 28)

	fmt.Println("Deleting user with ID 2")
	DeleteUser(2)
}
