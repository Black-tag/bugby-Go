module github.com/blacktag/bugby-Go

go 1.24.0

replace github.com/golang-jwt/jwt/v5 => github.com/golang-jwt/jwt/v5 v5.2.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.40.0
)
