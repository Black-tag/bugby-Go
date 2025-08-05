package api

import (
	// "context"
	"encoding/json"
	

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


// func TestGetBugbyIDHandler(t *testing.T) {
// 	cfg, mock := setupTest(t)
// 	defer cfg.SQLDB.Close()


// 	testbug := database.Bug{
// 		ID: uuid.New(),
// 		Title: "testing bugbyID function",
// 		Description: "hope it works",
// 		PostedBy: uuid.New(),
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "title", "description", "posted_by",
// 	 "created_at", "updated_at"}).AddRow(testbug.ID, testbug.Title, testbug.Description, testbug.PostedBy,
// 		 testbug.CreatedAt, testbug.UpdatedAt)

// 	mock.ExpectQuery("SELECT (.+) FROM bugs WHERE id = (.+)").WithArgs(testbug.ID).WillReturnRows(rows)
	
	
    

// 	req := httptest.NewRequest("GET", "/api/bugs/"+testbug.ID.String(), nil)
// 	req = req.WithContext(
// 		context.WithValue(req.Context(), http.ServerContextKey, map[string]string{
// 			"bugid": testbug.ID.String(),
// 		}),
// 	)
    
	
	
// 	w := httptest.NewRecorder()


// 	cfg.GetBugByIDHandler(w, req)

// 	res := w.Result()
//     defer res.Body.Close()

	
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response database.Bug
// 	err:= json.NewDecoder(w.Body).Decode(&response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, testbug.Title, response.Title)
// 	assert.Equal(t, testbug.Description, response.Description)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }