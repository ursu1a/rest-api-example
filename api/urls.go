package api

import "os"

var urls = map[string]string{
	"GOOGLE_OAUTH_PATH":          os.Getenv("GOOGLE_OAUTH_PATH"),
	"GOOGLE_OAUTH_CALLBACK_PATH": os.Getenv("GOOGLE_OAUTH_CALLBACK_PATH"),
	"TOKEN_REFRESH_PATH":         "/api/auth/refresh",
	"LOGIN_PATH":                 "/api/login",
	"REGISTER_PATH":              "/api/register",
	"GET_USER_PATH":              "/api/me",
	"VERIFY_EMAIL_PATH":          "/api/verify-email",
	"USERS_ALL_PATH":             "/api/users/all",
	"USERS_CREATE_PATH":          "/users/new",
	"USERS_DETAIL_PATH":          "/users/{id}",
	"USERS_DELETE_PATH":          "/users/{id}/delete",
	"SECURE_TEST_PATH":           "/secure/test",
}

var Urls = func(key string) string {
	v, found := urls[key]
	if found {
		return v
	}
	return ""
}
