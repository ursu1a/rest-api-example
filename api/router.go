package api

import (
	"backend/handlers"
	"backend/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NotFound(w http.ResponseWriter, req *http.Request) {
	http.Error(w, fmt.Sprintf(`{"error": "Route not found: %s"}`, req.URL.Path), http.StatusNotFound)
}

func securedHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Entered to secure path"))
}

func InitRouter() *mux.Router {
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
	router.HandleFunc(Urls("VERIFY_EMAIL_PATH"), handlers.HandleVerifyEmail)
	router.HandleFunc(Urls("REQUEST_RESET_PASSWORD_PATH"), handlers.HandleRequestResetPassword).Methods("POST")
	router.HandleFunc(Urls("UPDATE_PASSWORD_PATH"), handlers.HandleUpdatePassword).Methods("POST")

	// Users routes
	router.HandleFunc(Urls("USERS_ALL_PATH"), handlers.GetAllUsers)
	router.HandleFunc(Urls("USERS_CREATE_PATH"), handlers.CreateUser).Methods("POST")
	router.HandleFunc(Urls("USERS_DETAIL_PATH"), handlers.GetUser)
	router.HandleFunc(Urls("USERS_DELETE_PATH"), handlers.DeleteUser).Methods("DELETE")

	router.Handle(Urls("SECURE_TEST_PATH"), protectedRoute(http.HandlerFunc(securedHandler)))

	// Not found
	router.NotFoundHandler = http.HandlerFunc(NotFound)

	http.Handle("/", router)
	return router
}
