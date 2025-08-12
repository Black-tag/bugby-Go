package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)





func TestCreatUserHandler(t *testing.T) {
	cfg, mock := setupTest(t)
	defer cfg.SQLDB.Close()
	
	logger := slog.Default().With(
		"test", "TestCreateBugHandler",
	)
	logger.Info("testing started")

	testEmail := "testuser@gmail.com"
	testPassword := "mySecret"
	logger = logger.With(
		"testEmail", testEmail,
		"testPassword", testPassword,
	)
	hashed_password, err := utils.HashPassword(testPassword)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	logger = logger.With("hashedpassword in test", hashed_password)
	
	dbParams := database.CreateUserParams{
        Email:          testEmail,
        HashedPassword: hashed_password,  // This is what your DB layer expects
    }
	expectedUser := database.User{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: testEmail,
		HashedPassword: hashed_password,
		Role: "user",
	}

	expectedQuery := `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, role`

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email", "hashed_password", "role"}).AddRow(
		expectedUser.ID,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt, 
		expectedUser.Email, 
		expectedUser.HashedPassword, 
		expectedUser.Role)

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(dbParams.Email, sqlmock.AnyArg()).WillReturnRows(rows)

	testRequest := struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }{
        Email:    testEmail,
        Password: testPassword,
    } 
	requestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("failed to marshal json request body: %v", err)
	}
	logger = logger.With("requestBody", requestBody)

	req := httptest.NewRequest("POST", "/api/users/", bytes.NewBuffer(requestBody))

	w := httptest.NewRecorder()
	t.Logf("Making request to: %s", req.URL.Path)
	cfg.CreateUserHandler(w, req)
	
	t.Logf("Response status: %d", w.Code)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status code 201, got: %d.", w.Code)
	}


	var response database.User
		err = json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
	assert.Equal(t, testEmail, response.Email)
	// assert.Equal(t, hashed_password, response.HashedPassword)
	assert.NoError(t, mock.ExpectationsWereMet())
	logger.Info("test ended")



}