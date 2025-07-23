package api

import (
	"github.com/blacktag/bugby-Go/internal/database"
)

type APIConfig struct {
	DB *database.Queries
	SECRET string

} 