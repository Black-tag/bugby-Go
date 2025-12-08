package middleware

type contextKey string

const (
	UserIDKey          contextKey = "userID"
	TokenStringKey     contextKey = "tokenString"
	RoleKey            contextKey = "role"
	RefreshTokenString contextKey = "refreshTokenString"
	UserKey            contextKey = "user"
)
