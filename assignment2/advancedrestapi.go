package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "docs" // Swagger docs

	httpSwagger "github.com/swaggo/http-swagger" // Swagger middleware
)

var (
	sqlDB  *sql.DB  // for direct SQL queries
	gormDB *gorm.DB // for GORM queries
)

// @title           GoLang REST API by Bakytzhan
// @version         1.0
// @description     Advanced REST API with MySQL, SQL, GORM and Swagger documentation.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:8080
// @BasePath        /

type User struct {
	ID      uint    `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name" gorm:"unique;not null"`
	Age     int     `json:"age" gorm:"not null"`
	Profile Profile `json:"profile" gorm:"foreignKey:UserID"`
}

type Profile struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	UserID            uint   `json:"user_id" gorm:"unique;not null"`
	Bio               string `json:"bio"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

// Connects to MySQL using sql.DB
func connectSQL() {
	var err error

	connStr := "root:password@tcp(127.0.0.1:3306)/gocon"
	sqlDB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	sqlDB.SetMaxOpenConns(25) // Connection pooling
	sqlDB.SetMaxIdleConns(25)
	fmt.Println("Connected to MySQL using sql.DB!")
}

// Connects to MySQL using GORM
func connectGORM() {
	var err error
	dsn := "root:password@tcp(127.0.0.1:3306)/gocon?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) // Use gormDB
	if err != nil {
		log.Fatal("Failed to connect to MySQL database:", err)
	}
	fmt.Println("Connected to MySQL using GORM!")
}

// @Summary Get Users with optional filtering and pagination (SQL)
// @Description Retrieve a list of users from MySQL using SQL queries with optional filtering by age and pagination.
// @Tags Users
// @Produce json
// @Param age query string false "Filter by age"
// @Param sort query string false "Sort by name (asc or desc)"
// @Param page query string false "Pagination page number"
// @Success 200 {array} User
// @Failure 500 {object} map[string]string
// @Router /sql/users [get]
func getUsersSQL(w http.ResponseWriter, r *http.Request) {
	ageFilter := r.URL.Query().Get("age")
	sortOrder := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	limit := 10
	offset := (page - 1) * limit
	query := "SELECT * FROM users"
	if ageFilter != "" {
		query += " WHERE age = " + ageFilter
	}
	if sortOrder == "asc" {
		query += " ORDER BY name ASC"
	} else if sortOrder == "desc" {
		query += " ORDER BY name DESC"
	}
	query += " LIMIT ? OFFSET ?"

	rows, err := sqlDB.Query(query, limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, "Failed to scan user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @Summary Create a new User (SQL)
// @Description Insert a new user into MySQL using SQL queries with name uniqueness validation.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body User true "User"
// @Success 201 {object} User
// @Failure 500 {object} map[string]string
// @Router /sql/users [post]
func createUserSQL(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err := sqlDB.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// @Summary Get Users with optional filtering and pagination (GORM)
// @Description Retrieve a list of users from MySQL using GORM with optional filtering by age and pagination.
// @Tags Users
// @Produce json
// @Param age query string false "Filter by age"
// @Param sort query string false "Sort by name (asc or desc)"
// @Param page query string false "Pagination page number"
// @Success 200 {array} User
// @Failure 500 {object} map[string]string
// @Router /gorm/users [get]
func getUsersGORM(w http.ResponseWriter, r *http.Request) {
	ageFilter := r.URL.Query().Get("age")
	sortOrder := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	var users []User
	query := gormDB.Limit(10).Offset((page - 1) * 10)
	if ageFilter != "" {
		query = query.Where("age = ?", ageFilter)
	}
	if sortOrder == "asc" {
		query = query.Order("name asc")
	} else if sortOrder == "desc" {
		query = query.Order("name desc")
	}
	query.Find(&users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @Summary Create a new User (GORM)
// @Description Insert a new user into MySQL using GORM with name uniqueness validation.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body User true "User"
// @Success 201 {object} User
// @Failure 500 {object} map[string]string
// @Router /gorm/users [post]
func createUserGORM(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := gormDB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func main() {
	// Connect to both SQL and GORM databases
	connectSQL()
	connectGORM()

	// Set up Swagger documentation
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Set up routes
	http.HandleFunc("/sql/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsersSQL(w, r)
		} else if r.Method == http.MethodPost {
			createUserSQL(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/gorm/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsersGORM(w, r)
		} else if r.Method == http.MethodPost {
			createUserGORM(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
