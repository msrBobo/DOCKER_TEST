package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

var db *sql.DB

func main() {
	// Initialize the database connection
	var err error
	db, err = sql.Open("postgres", "postgres://bobo:1234@localhost/userdb?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize the Gin router
	router := gin.Default()

	// Define routes
	router.POST("/register", registerHandler)
	router.POST("/login", loginHandler)

	// Run the server
	router.Run(":8080")
}

func registerHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the username or email already exists
	if userExists(user.Username, user.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.Status(http.StatusCreated)
}

func loginHandler(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Retrieve the user from the database
	row := db.QueryRow("SELECT id, password FROM users WHERE username = $1", user.Username)
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare the hashed password with the provided password
	if !checkPasswordHash(user.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token or set session cookie for authentication
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func userExists(username, email string) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2", username, email)
	if err := row.Scan(&count); err != nil {
		return false
	}
	return count > 0
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	hashedPassword := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return hashedPassword, nil
}

func checkPasswordHash(hashedPassword, password string) bool {
	decodedHashedPassword, err := base64.URLEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(decodedHashedPassword[:16]) // Extract salt from hashed password
	return bcrypt.CompareHashAndPassword(decodedHashedPassword[16:], hash.Sum(nil)) == nil
}
