package handlers

/* User profile handlers */

import (
	"backend/config"
	"backend/db"
	"encoding/json"
	"log"
	"net/http"
)

type UpdatedUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
		responseError(w, err.Error())
		return
	}

	responseJSON(w, user)
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get("UserID")
	log.Printf("Trying to update user: %s", userID)

	DBConn := config.App.DB
	var userData UpdatedUser
	var user db.User

	if err := json.NewDecoder(req.Body).Decode(&userData); err != nil {
		responseError(w, err.Error())
		return
	}

	if err := DBConn.First(&user, userID).Error; err != nil {
		responseError(w, err.Error())
		return
	}

	user.Name = userData.Name
	user.Email = userData.Email

	DBConn.Save(&user)

	responseJSON(w, user)
}

// Delete user by ID
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	DBConn := config.App.DB
	var user db.User

	userID := req.Header.Get("UserID")
	log.Printf("Trying to delete user: %s", userID)

	if err := DBConn.Delete(&db.User{}, userID).Error; err != nil {
		responseError(w, err.Error())
		return
	}

	responseJSON(w, user)
}
