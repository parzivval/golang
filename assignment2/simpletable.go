package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		age INT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table 'users' created.")
}

func insertData(db *sql.DB, name string, age int) {
	query := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err := db.Exec(query, name, age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted user: %s, Age: %d\n", name, age)
}

func queryData(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var id int
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/gocon")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected.")

	createTable(db)

	insertData(db, "Alice", 30)
	insertData(db, "Bob", 25)

	queryData(db)
}
