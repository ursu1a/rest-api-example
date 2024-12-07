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
	router.HandleFunc(Urls("VERIFY_EMAIL_PATH"), handlers.HandleVerifyEmail)
	router.HandleFunc(Urls("REQUEST_RESET_PASSWORD_PATH"), handlers.HandleRequestResetPassword).Methods("POST")
	router.HandleFunc(Urls("UPDATE_PASSWORD_PATH"), handlers.HandleUpdatePassword).Methods("POST")

	// User routes
	router.Handle(Urls("GET_USER_PATH"), protectedRoute(http.HandlerFunc(handlers.HandleGetUser)))
	router.Handle(Urls("PROFILE_UPDATE_PATH"), protectedRoute(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
	router.Handle(Urls("ACCOUNT_DELETE_PATH"), protectedRoute(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")

	// Not found
	router.NotFoundHandler = http.HandlerFunc(NotFound)

	http.Handle("/", router)
	return router
}
