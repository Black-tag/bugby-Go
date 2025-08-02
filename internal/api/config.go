package api

import (
	"database/sql"

	"github.com/blacktag/bugby-Go/internal/database"
)

type APIConfig struct {
	DB *database.Queries
	SECRET string
	SQLDB *sql.DB

} 