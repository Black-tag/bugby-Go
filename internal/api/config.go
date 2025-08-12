package api

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blacktag/bugby-Go/internal/database"
)

type APIConfig struct {
	DB     *database.Queries
	SECRET string
	SQLDB  *sql.DB
}

func setupTest(t *testing.T) (*APIConfig, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock DB: %v", err)
	}
	return &APIConfig{
		DB:    database.New(db),
		SQLDB: db,
	}, mock
}
