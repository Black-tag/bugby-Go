# 🐞 Bugby-Go

A production-grade bug tracking service built with Go, PostgreSQL, JWT authentication, and Casbin RBAC.  
**Live on Railway | Automated CI/CD | Dockerized | Observability Ready**

---

## 🚀 Overview

Bugby-Go is a scalable, lightweight bug tracking system designed for modern teams and startups.  
It demonstrates advanced backend patterns, cloud-native deployment, and robust security—all written in idiomatic Go.

---

## 🛠️ Tech Stack

- **Language:** Go (net/http, SQLC)
- **Database:** PostgreSQL
- **Auth:** JWT (with bcrypt for password hashing)
- **Authorization:** Casbin RBAC (configurable model/policies)
- **API Docs:** Swagger/OpenAPI (live at `/swagger/`)
- **Testing:** Go test, testify, sqlmock (unit + integration)
- **Observability:** Prometheus metrics (`/metrics`), structured logging (Go slog)
- **Infrastructure:** Docker, Docker Compose, GitHub Actions (CI/CD)
- **Deployment:** Railway (cloud), local dev via Docker Compose

---

## 🎯 Features

- **User Registration & Authentication:** JWT-based, password hashing with bcrypt.
- **Role-Based Access Control:** User & Admin roles, fine-grained permissions via Casbin.
- **Bug CRUD:** Create, read, update, delete bugs with ownership and role checks.
- **Rate Limiting:** Middleware against abuse.
- **Database Migrations:** Automated with Goose.
- **Automated Testing:** 40%+ coverage (unit/integration), easy to extend.
- **API Documentation:** Live Swagger/OpenAPI (`/swagger/` endpoint).
- **Metrics & Monitoring:** `/metrics` endpoint for Prometheus, ready for Grafana dashboards.
- **Structured Logging:** Go’s slog for JSON logs, context-rich.
- **Middleware:** For logging, authentication, authorization, and metrics.
- **CI/CD:** Automated builds, tests, deploys with GitHub Actions.
- **Containerization:** Dockerfile & docker-compose for local and prod environments.

---

## 📦 Quickstart

### 1. Clone & Setup

```bash
git clone https://github.com/Black-tag/bugby-Go.git
cd bugby-Go
```

### 2. Environment Variables

Copy `.env.example` to `.env` and fill your secrets:

```
DB_URL=postgres://user:password@localhost:5432/bugby?sslmode=disable
JWT_SECRET=supersecret
PORT=8080
```

### 3. Start Services

```bash
docker-compose up --build
```

### 4. DB Migrations

```bash
./scripts/dev_migrate.sh
```

### 5. Run Tests

```bash
go test ./... -cover
```

### 6. View API Docs

Visit [http://localhost:8080/swagger/](http://localhost:8080/swagger/) (after starting the server).

---

## 🔒 Security & Auth

- **JWT Authentication:** Secure endpoints, refresh tokens, blacklisting.
- **Casbin RBAC:** Edit `rbac_model.conf` and `rbac_policy.csv` to customize permissions.
- **Password Hashing:** All credentials securely hashed.

---

## 📊 Observability

- **Metrics:** Prometheus-ready at `/metrics`.
- **Structured Logging:** All requests and errors are logged in structured JSON format.

---

## 🧪 Testing & Coverage

- **Unit/Integration:** Tests for API, DB, and middleware.
- **Coverage:** ≥40%, extendable (see `internal/api/*_test.go`).

---

## 📝 API Endpoints

See Swagger docs at `/swagger/` for live, interactive API documentation.

Example endpoints:

- `POST /api/users` — Register user
- `POST /api/login` — Authenticate
- `POST /api/bugs` — Create bug (JWT required)
- `GET /api/bugs` — List bugs
- `PUT /api/users` — Update user (JWT required)
- ...and more!

---

## 🏗️ Architecture

- **Handlers:** Clean separation between handlers, middleware, and utility logic.
- **Middleware:** Composable for logging, auth, rate limiting, and metrics.
- **DB Access:** SQLC-generated code, safe queries, easy mocking for tests.

---

## 🛡️ Advanced Highlights

- **Casbin RBAC:** Extensible model/policy for robust authorization.
- **Prometheus Metrics:** Custom metrics for HTTP requests, DB queries, bug count.
- **Live Swagger UI:** Always up-to-date with code annotations.
- **CI/CD:** Automated pipeline for build, test, deploy on every commit.

---

## 💡 Planned / Roadmap

- OAuth2 login (GitHub/Google)
- Redis caching
- gRPC microservices
- Advanced monitoring dashboards (Grafana)
- Frontend dashboard

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!  
See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## 📄 License

MIT — see [LICENSE](LICENSE)

---

## 👤 Author

**Anand Unni**  
[GitHub](https://github.com/Black-tag) | unnianandunni007@gmail.com

---

## 🟢 Live Demo

Deployed on Railway (link in repo description or ask maintainer).

---

## 📚 References

- [Go Documentation](https://golang.org/doc/)
- [Casbin RBAC Docs](https://casbin.org/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Swaggo for Go](https://github.com/swaggo/swag)