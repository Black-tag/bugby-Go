package api

import (
	// "context"
	// "context"
	// "context"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"regexp"

	"log/slog"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)


func setupTest (t *testing.T) (*APIConfig, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	return &APIConfig{
		DB: database.New(db),
		SQLDB: db,

	}, mock
}





func TestGetBugHandler (t *testing.T) {
	
	cfg, mock :=setupTest(t)
	defer cfg.SQLDB.Close()

	
	expectedBugs := []database.Bug{
        {
            ID:          uuid.New(),
            Title:       "Test bug 1",
            Description: "test description 1",
            PostedBy:    uuid.New(),
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        },
        {
            ID:          uuid.New(),
            Title:       "Test bug 2",
            Description: "test description 2",
            PostedBy:    uuid.New(),
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        },
    }
	rows := sqlmock.NewRows([]string{"id", "title", "description", "posted_by", "created_at", "updated_at" })
	for _, bug := range expectedBugs {
		rows.AddRow(bug.ID, bug.Title, bug.Description, bug.PostedBy, bug.CreatedAt, bug.UpdatedAt)
	}
	mock.ExpectQuery("SELECT (.+) FROM bugs").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/api/bugs/", nil)
	
	w := httptest.NewRecorder()

	cfg.GetBugsHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response []database.Bug
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, len(expectedBugs))
	for i := range expectedBugs {
        assert.Equal(t, expectedBugs[i].ID, response[i].ID)
        assert.Equal(t, expectedBugs[i].Title, response[i].Title)
        assert.Equal(t, expectedBugs[i].Description, response[i].Description)
        
    }
    
    
    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestGetBugbyIDHandler(t *testing.T) {
	
	slog.Info("handler entered")
	logger := slog.Default().With(
		"handler", "TesBugIDHandler",
		
	)
	cfg, mock := setupTest(t)
	defer cfg.SQLDB.Close()


	testbug := database.Bug{
		ID: uuid.New(),
		Title: "testing bugbyID function",
		Description: "hope it works",
		PostedBy: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "title", "description", "posted_by",
	 "created_at", "updated_at"}).AddRow(testbug.ID, testbug.Title, testbug.Description, testbug.PostedBy,
		 testbug.CreatedAt, testbug.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta("-- name: GetBugsByID :one SELECT id, title, description, posted_by, created_at, updated_at FROM bugs WHERE Id = $1")).WithArgs(testbug.ID).WillReturnRows(rows)
	logger = logger.With("rows", rows)
	
	
    
	logger = logger.With("tetsbugId", testbug.ID.String())
	mux := http.NewServeMux()
    mux.HandleFunc("/api/bugs/{bugid}", cfg.GetBugByIDHandler)
	req := httptest.NewRequest("GET", "/api/bugs/"+testbug.ID.String(), nil)
	
    
	
	
	w := httptest.NewRecorder()
	t.Logf("Making request to: %s", req.URL.Path)
    t.Logf("Expected bug ID: %s", testbug.ID)
	mux.ServeHTTP(w, req)
	t.Logf("Response status: %d", w.Code)
    body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, string(body))

	}
    t.Logf("Response body: %s", string(body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got: %d. Body: %s", w.Code, string(body))
	}



	var response database.Bug
    err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	assert.Equal(t, testbug.Title, response.Title)
	assert.Equal(t, testbug.Description, response.Description)
	assert.NoError(t, mock.ExpectationsWereMet())
	logger.Info("test ended")
}



func TestCreateBugHandler (t *testing.T) {
	cfg, mock := setupTest(t)
	defer cfg.SQLDB.Close()
	slog.Info("started CreateBugHandler Test")
	logger := slog.Default().With(
		"test", "TestcreateBugHandler",
		
	)
	userID := uuid.New()
	logger = logger.With("userID", userID)


	testbug := database.CreateBugParams{
		Title: "testing CreateBugHandler",
		Description: "it should work",
		
	}
	expectedBug := database.Bug{
		ID: uuid.New(),
		Title: testbug.Title,
		Description: testbug.Description,
		PostedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{ "id", "title", "description", "posted_by", "created_at", "updated_at",
	}).AddRow(expectedBug.ID, expectedBug.Title, expectedBug.Description, expectedBug.PostedBy, expectedBug.CreatedAt, expectedBug.UpdatedAt)
	expectedQuery := `-- name: CreateBug :one INSERT INTO bugs (id, title, description, posted_by, created_at, updated_at) VALUES ( gen_random_uuid(), $1, $2, $3, NOW(), NOW() ) RETURNING id, title, description, posted_by, created_at, updated_at`
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(testbug.Title, testbug.Description, userID).WillReturnRows(rows)
	logger = logger.With("rows", rows)

	requestBody, err := json.Marshal(testbug)
    if err != nil {
        t.Fatalf("failed to marshal request body: %v", err)
    }

	logger = logger.With("requestBody", requestBody)
	req := httptest.NewRequest("POST", "/api/bugs", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	ctx := context.WithValue(req.Context(), "userID", userID)

	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	t.Logf("Making request to: %s", req.URL.Path)
    cfg.CreateBugHandler(w, req)
	
	t.Logf("Response status: %d", w.Code)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status code 201, got: %d.", w.Code)
	}

	var response database.Bug
	err = json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	assert.Equal(t, testbug.Title, response.Title)
	assert.Equal(t, testbug.Description, response.Description)
	assert.Equal(t, testbug.PostedBy, response.PostedBy)

}

func TestUpdateBugHandler (t *testing.T) {
	logger :=slog.Default().With(
		"test", "TestupdateBugHandler",
	)
	cfg, mock := setupTest(t)
	defer cfg.SQLDB.Close()
	logger.Info("started tests")
	userID := uuid.New()
	logger = logger.With("userID", userID)

	bugID := uuid.New()
	logger = logger.With("bugID", bugID)

	
	 testRequest := struct {
		
        Title       *string `json:"title"`
        Description *string `json:"description"`
    }{
        Title:       stringPtr("this is to update the bug"),
        Description: stringPtr("hope this works"),
    }
	existingBug := database.Bug{
        ID:          bugID,
        Title:       "original title",
        Description: "original description",
        PostedBy:    userID,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
	expectedBug := database.Bug{
		ID: bugID,
		Title: *testRequest.Title,
		Description: *testRequest.Description,
		PostedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

	}
	logger = logger.With("testRequest", testRequest)

	expectedQuery := `-- name: UpdateBugByID :exec UPDATE bugs SET title = COALESCE($2, title), description = COALESCE($3, description), updated_at = Now() WHERE id = $1`

	rows := sqlmock.NewRows([]string{ "id", "title", "description", "posted_by", "created_at", "updated_at"}).AddRow(
	expectedBug.ID, expectedBug.Title, expectedBug.Description, expectedBug.PostedBy, expectedBug.CreatedAt, expectedBug.UpdatedAt,
	)
	    mock.ExpectQuery(regexp.QuoteMeta(
        `SELECT id, title, description, posted_by, created_at, updated_at FROM bugs WHERE Id = $1`,
    )).WithArgs(bugID).
        WillReturnRows(sqlmock.NewRows([]string{
            "id", 
			"title", 
			"description", 
			"posted_by", 
			"created_at", 
			"updated_at",
        }).AddRow(
            existingBug.ID, 
            existingBug.Title, 
            existingBug.Description, 
            existingBug.PostedBy, 
            existingBug.CreatedAt, 
            existingBug.UpdatedAt,
        ))
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
	WithArgs(
		bugID,
		expectedBug.Title, 
		expectedBug.Description).WillReturnResult(sqlmock.NewResult(1, 1))
	
	mock.ExpectQuery(regexp.QuoteMeta(
        `SELECT id, title, description, posted_by, created_at, updated_at FROM bugs WHERE Id = $1`,
    )).WithArgs(bugID).
        WillReturnRows(sqlmock.NewRows([]string{
            "id", 
			"title", 
			"description", 
			"posted_by", 
			"created_at", 
			"updated_at",
        }).AddRow(
            existingBug.ID, 
            expectedBug.Title, 
            expectedBug.Description, 
            existingBug.PostedBy, 
            existingBug.CreatedAt, 
            existingBug.UpdatedAt,
        ))

	logger = logger.With("rows", rows)
	requestBody, err := json.Marshal(testRequest)
	if err != nil {
		t.Fatalf("failed to marshall the requestBody: %v", err)

	}
	logger = logger.With("requestBody", requestBody)

	mux :=http.NewServeMux() 
	mux.HandleFunc("/api/bugs/{bugid}", cfg.UpdateBugHandler)
	logger = logger.With("bugId", bugID)
	req := httptest.NewRequest("POST", "/api/bugs/"+bugID.String(), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	ctx := context.WithValue(req.Context(), "userID", userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	t.Logf("Making request to: %s", req.URL.Path)
	t.Logf("Expected bug ID: %s", bugID)
	mux.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("cannot read body")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected code 200 got: %d, Body: %s", w.Code, string(body))
	}

	var response database.Bug
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	assert.Equal(t, expectedBug.Title, response.Title)
	assert.Equal(t, expectedBug.Description, response.Description)
	
	assert.NoError(t, mock.ExpectationsWereMet())
	logger.Info("test ended")


}
func stringPtr (s string) *string {
	return &s
}


func TestDeleteBugByIDHandler (t *testing.T) {
	logger := slog.Default().With(
		"test", "testDeleteBugByIDHandler",
	)
	logger.Info("satrted delete test")

	userID := uuid.New()
	logger = logger.With("userID", userID)
	bugID := uuid.New()
	logger = logger.With("bugID", bugID)
	cfg, mock := setupTest(t)
	defer cfg.SQLDB.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, description, posted_by, created_at, updated_at FROM bugs WHERE Id = $1`)).
        WithArgs(bugID).
        WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "posted_by", "created_at", "updated_at"}).
            AddRow(bugID, "test bug", "test description", userID, time.Now(), time.Now()))

	expectedQuery := `-- name: DeleteBugByID :exec
DELETE FROM bugs
WHERE id = $1`
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(bugID).WillReturnResult(sqlmock.NewResult(1, 1))
	testUser := database.User{
        ID:    userID,
        Role:  "admin",
    }

	mux := http.NewServeMux()
	mux.HandleFunc("/api/bugs/{bugid}", cfg.DeleteBugByIDHandler)
	req := httptest.NewRequest("DELETE", "/api/bugs/"+bugID.String(), nil)
	
	// type contextKey string
    // const (
    //     userIDKey contextKey = "userID"
    //     userKey   contextKey = "user"
    // )
	
	
	
	ctx := context.WithValue(req.Context(), "userID", userID)
	ctx = context.WithValue(ctx, "userRole", "admin")
	ctx = context.WithValue(ctx, "user", testUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	t.Logf("Making request to: %s", req.URL.Path)
	t.Logf("Expected bug ID: %s", bugID)
	mux.ServeHTTP(w, req)


	if w.Code != http.StatusNoContent {
		t.Logf("Response body: %s", w.Body.String())
		t.Fatalf("expected status code 204 but got: %d", w.Code)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
	logger.Info("test ended")

}