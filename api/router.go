package api

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NotFound(w http.ResponseWriter, req *http.Request) {
	message := fmt.Sprintf("Route not found: \"%s\"", req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(message))
}

func securedHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Entered to secure path"))
}

func Router() {
	router := mux.NewRouter()

	// Protected routes
	protectedRoute := middleware.Authenticate

	// Authentication routes
	router.HandleFunc(Urls("GOOGLE_OAUTH_PATH"), handlers.HandleGoogleLogin)
	router.HandleFunc(Urls("GOOGLE_OAUTH_CALLBACK_PATH"), handlers.HandleGoogleCallback)
	router.HandleFunc(Urls("LOGIN_PATH"), handlers.HandleLogin).Methods("POST")
	router.HandleFunc(Urls("REGISTER_PATH"), handlers.HandleRegister).Methods("POST")
	router.HandleFunc(Urls("TOKEN_REFRESH_PATH"), handlers.HandleRefreshToken).Methods("POST")
	router.Handle(Urls("GET_USER_PATH"), protectedRoute(http.HandlerFunc(handlers.HandleGetUser)))

	// Users routes
	router.HandleFunc(Urls("USERS_ALL_PATH"), handlers.GetAllUsers)
	router.HandleFunc(Urls("USERS_CREATE_PATH"), handlers.CreateUser).Methods("POST")
	router.HandleFunc(Urls("USERS_DETAIL_PATH"), handlers.GetUser)
	router.HandleFunc(Urls("USERS_DELETE_PATH"), handlers.DeleteUser).Methods("DELETE")

	router.Handle(Urls("SECURE_TEST_PATH"), protectedRoute(http.HandlerFunc(securedHandler)))

	// Not found
	router.NotFoundHandler = http.HandlerFunc(NotFound)

	http.Handle("/", router)
	config.App.Router = router
}
