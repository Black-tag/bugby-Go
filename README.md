# 🐞 Bugby-Go
A production-grade bug tracking service built with Go, PostgreSQL, and JWT authentication.

# Description
Bugby-Go is a lightweight, scalable bug tracking system built in Go.  
It features user authentication, bug reporting, rate limiting, CI/CD workflows,  
and is fully Dockerized for easy deployment. Designed to be a backend-heavy project  
that small startups or teams can plug in for tracking and managing issues.

# Tech Stack
- **Language:** Go (net/http, SQLC)
- **Database:** PostgreSQL
- **Auth:** JWT (with bcrypt for password hashing)
- **Caching / Future-ready:** Redis (optional)
- **Observability:** Prometheus + Grafana (planned)
- **Infrastructure:** Docker, CI/CD (GitHub Actions)
- **Deployment:** Railway / Render (free-tier hosting)

# Features
✅ User registration & authentication (JWT)  
✅ Create, read, update, delete (CRUD) bugs  
✅ Role-based access (users & admins)  
✅ Rate limiting (protection against abuse)  
✅ Database migrations (Goose)  
✅ Unit & integration tests (40%+ coverage)  
✅ Dockerized for local and production environments  
🔜 Planned: Redis caching, Prometheus monitoring, Grafana dashboards  

# Architecture 
- Go (net/http) handles requests
- Postgres stores users and bugs
- Redis (optional) will cache frequent queries
- Prometheus collects metrics for Grafana dashboards

# Installation(Local)
1. Clone the repo
```git clone https://github.com/<your-username>/bugby-go.git```
```cd bugby-go```

2. Start services (Postgres + app)
```docker-compose up --build```

3. Run migrations
```./scripts/dev_migrate.sh```

4. Run tests
```go test ./... -cover```


#Environment variables
DB_URL=postgres://user:password@localhost:5432/bugby?sslmode=disable  
JWT_SECRET=supersecret 

# API Endpoints 
```DB_URL=postgres://user:password@localhost:5432/bugby?sslmode=disable```  
```JWT_SECRET=supersecret``` 

# deployment

# Contributing & Roadmap
Contributions are welcome!  
Planned features:  
- OAuth2 login (GitHub/Google)  
- gRPC microservices  
- Prometheus metrics + Grafana dashboards  
- Frontend dashboard



