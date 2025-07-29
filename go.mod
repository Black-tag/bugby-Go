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

require golang.org/x/time v0.12.0

require (
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/casbin/casbin/v2 v2.110.0 // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
)
