package handlers
/* User control handlers */

import (
	"backend/config"
	"backend/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Get all users
func GetAllUsers(w http.ResponseWriter, req *http.Request) {
	DBConn := config.App.DB
	var users []db.User

	if err := DBConn.Find(&users).Error; err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, users)
}

// Create a new user
func CreateUser(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		responseError(w, "Method not allowed")
		return
	}

	DBConn := config.App.DB
	var user db.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		responseError(w, err.Error())
		return
	}

	if err := DBConn.FirstOrCreate(&user).Error; err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, user)
}

// Retrieve user by ID
func GetUser(w http.ResponseWriter, req *http.Request) {
	DBConn := config.App.DB
	var user db.User

	vars := mux.Vars(req)
	userIdStr := vars["id"]

	if _, err := strconv.Atoi(userIdStr); err != nil {
		responseError(w, fmt.Sprintf("User with ID \"%v\" is not found. %s", userIdStr, err.Error()))
		return
	}

	if err := DBConn.First(&user, userIdStr).Error; err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, user)
}

// Get user information from token
func HandleGetUser(w http.ResponseWriter, req *http.Request) {
	DBConn := config.App.DB
	var user db.User

	userID := req.Header.Get("UserID")
	email := req.Header.Get("Email")
	if userID == "" || email == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if err := DBConn.First(&user, userID).Error; err != nil {
		responseJSON(w, err.Error())
		return
	}

	responseJSON(w, user)
}

// Delete user by ID
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	DBConn := config.App.DB
	var user db.User

	vars := mux.Vars(req)
	userIdStr := vars["id"]

	if _, err := strconv.Atoi(userIdStr); err != nil {
		responseError(w, fmt.Sprintf("User with ID \"%v\" is not found. %s", userIdStr, err.Error()))
		return
	}

	if err := DBConn.Delete(&user, userIdStr).Error; err != nil {
		responseError(w, err.Error())
		return
	} else if DBConn.RowsAffected < 1 {
		responseError(w, fmt.Sprintf("User with ID \"%v\" is not found.", userIdStr))
		return
	}

	responseJSON(w, user)
}

