package api

import "backend/utils"

var urls = map[string]string{
	"GOOGLE_OAUTH_PATH":           utils.GetEnv("GOOGLE_OAUTH_PATH", "/api/auth/google"),
	"GOOGLE_OAUTH_CALLBACK_PATH":  utils.GetEnv("GOOGLE_OAUTH_CALLBACK_PATH", "/api/auth/google/callback"),
	"TOKEN_REFRESH_PATH":          "/api/auth/refresh",
	"LOGIN_PATH":                  "/api/login",
	"REGISTER_PATH":               "/api/register",
	"GET_USER_PATH":               "/api/me",
	"VERIFY_EMAIL_PATH":           "/api/verify-email",
	"REQUEST_RESET_PASSWORD_PATH": "/api/request-reset-password",
	"UPDATE_PASSWORD_PATH":        "/api/update-password",
	"PROFILE_UPDATE_PATH":         "/api/profile/update",
	"ACCOUNT_DELETE_PATH":         "/api/account/delete",
}

var Urls = func(key string) string {
	v, found := urls[key]
	if found {
		return v
	}
	return ""
}
