package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	

	// "github.com/casbin/casbin/v2/log"
	"github.com/google/uuid"
)


type CreateBugResponse struct {
		ID           uuid.UUID `json:"bug_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		PostedBy     uuid.UUID `json:"posted_by"`
		CreatedBy    time.Time `json:"created_at"`
		Updated_at   time.Time `json:"updated_at"`

	}

type CreateBugRequest struct {
	Title        string    `json:"title" example:"This is the bug needed"`
	Description  string    `json:"description" example:"this is descrption"`
	PostedBy     uuid.UUID `json:"posted_by" example:"9b733930-ef6f-4b01-add2-f410962ec695"`
}

type UpdateBugRequest struct {
	Title       *string    `json:"title" example:"This is the bug needed"`
	Description *string    `json:"description" example:"this is descrption"`
		
		
}

// @Summary Create bugs
// @Description Existing users can create bugs 
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateBugRequest true "bug creation data" 
// @Success 201 {object} CreateBugResponse
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 404 {object} utils.ErrorResponse "Not Found - Resource doesn't exist"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /bugs [post]
// @Security BearerAuth
func (cfg *APIConfig) CreateBugHandler (w http.ResponseWriter, r *http.Request){
	slog.Info("handler entered")
	logger := slog.Default().With(
		"handler", "CreateBugHandler",
		"method", r.Method,
		"path", r.URL.Path,
	)
	w.Header().Set("Content-Type", "application/json")
	userIDValue := r.Context().Value("userID")
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		logger.Error("user id missing in context")
		utils.RespondWithError(w, http.StatusNotFound, "cannot find user ID")
		return
	}
	logger = logger.With("user_id", userID)
	

	var req CreateBugRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error("failed to decode json", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "error decoding json")
		return
	}
	logger = logger.With("bug_title", req.Title)
	bug, err := cfg.DB.CreateBug(r.Context(),database.CreateBugParams{
		Title: req.Title,
		Description: req.Description,
		PostedBy: userID,
	})
	if err != nil {
		logger.Error("database operation failed", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create bug")
		return
	}

	logger.Info("bug created successfully", "bug_id", bug.ID)
	utils.RespondWithJSON(w, http.StatusCreated, CreateBugResponse{
		ID: bug.ID,
		Title: bug.Title,
		Description: bug.Description,
		PostedBy: bug.PostedBy,
		CreatedBy: bug.CreatedAt,
		Updated_at: bug.UpdatedAt,
	})
	slog.Info("about to respond")
	
}
// @Summary Get existing  bugs
// @Description  users can get all existing bugs
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} database.Bug
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /bugs [get]
// @Security BearerAuth
func (cfg *APIConfig) GetBugsHandler (w http.ResponseWriter, r *http.Request) {
	bugs, err := cfg.DB.GetAllBugs(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldnt fetch bugs")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, bugs)
}

// @Summary GET bug by id
// @Description Existing users can update their info using email and password
// @Tags bugs
// @Accept json
// @Produce json
// @Param bugid path string true "Bug ID" example:"87f0ea02-7b24-41bd-8418-0831a019fc87"  
// @Success 200 {object} database.Bug
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /bugs/{bugid} [get]
// @Security BearerAuth
func (cfg *APIConfig) GetBugByIDHandler (w http.ResponseWriter, r *http.Request) {
	slog.Info("Handler started")
	logger := slog.Default().With(
		"handler", "CreateBugHandler",
		"methid", r.Method,
		"path", r.URL.Path,
	)
	bugIDParam := r.PathValue("bugid")
	

	if bugIDParam == "" {
		logger.Info("no id given by user")
		utils.RespondWithError(w, http.StatusBadRequest, " no id given ")
		return
	}
	bugID, err := uuid.Parse(bugIDParam)
	if err != nil {
		logger.Error("invalid id format", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "wrong Id format ")
		return
	}
	logger.Info("parsed uuid", "bugID", bugID)

	logger.Info("calling database")
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil {
		logger.Error("database error", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, " bug not found ")
		return
	}
	logger.Info("response ready", "bug", bug)
	utils.RespondWithJSON(w, http.StatusOK, bug)
}

// @Summary Update an existing  bug
// @Description Existing users can update their bug
// @Tags users
// @Accept json
// @Produce json
// @Param bugid path string true "Bug ID" example:"87f0ea02-7b24-41bd-8418-0831a019fc87"  
// @Param request body UpdateBugRequest true "bug updation data" 
// @Success 200 {object} database.Bug
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Missing/invalid credentials"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /bug/{bugid} [put]
// @Security BearerAuth
func (cfg *APIConfig) UpdateBugHandler (w http.ResponseWriter, r *http.Request) {
	logger := slog.With(
		"handler", "UpdateBugHandler",
		"method", r.Method, 
		"path", r.URL.Path,
	)
	logger.Info("Entered handler")
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	logger = logger.With("userID", userID)
	if !ok {
		logger.Error(" user id not given or invalid")
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}
	bugParam := r.PathValue("bugid")
	
	if bugParam == "" {
		logger.Info("no bugID in request")
		utils.RespondWithError(w, http.StatusBadRequest, " no bugID given ")
		return
	}
	bugID, err := uuid.Parse(bugParam)
	if err != nil {
		logger.Error("given id format is wrong", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "wrong format Id")
		return
	}
	logger = logger.With("bugId", bugID)
	
	var req UpdateBugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("given request body in wrong format", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	logger.Info("doing database operation")
	if err != nil { 
		logger.Error("databse error bug not found in database", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "no bug  found with the id")
		return

	}
	logger = logger.With("bug", bug)
	if userID != bug.PostedBy {
		logger.Error("unauthorised to edit the bug, not owned by user", "error", err)
		utils.RespondWithError(w, http.StatusUnauthorized, "only author can delete the bug")
		return

	}
	params := database.UpdateBugByIDParams{
		ID: bugID,
		Title: toNullString(req.Title),
		Description: toNullString(req.Description),
		
	}
	logger = logger.With("params", params)
	
	err = cfg.DB.UpdateBugByID(r.Context(), params)
	logger.Info("doing database updation")
	if err != nil {
		logger.Error("updating bug in databse failed", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot update bug")
		return
	}
	updatedbug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot fetch updated bug")
		return
	}
	logger = logger.With("updatedbug", updatedbug)

	utils.RespondWithJSON(w, http.StatusOK, updatedbug)
	logger.Info("completed updation")
}
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
// @Summary Delete an existing  user
// @Description admin can delete bugs using their id
// @Tags bugs
// @Accept json
// @Produce json
// @Success 204 {string} string "No content" 
// @Failure 400 {object} utils.ErrorResponse "Bad Request - Invalid input"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - Missing/invalid credentials"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Param bugid path string true "Bug ID" example:"87f0ea02-7b24-41bd-8418-0831a019fc87" 
// @Router /bugs/{bugid} [delete]
// @Security BearerAuth
func (cfg *APIConfig) DeleteBugByIDHandler (w http.ResponseWriter, r *http.Request) {
	logger := slog.Default().With(
		"handler", "DeleteBugByIDHandler",
		"method", r.Method,
		"path", r.URL.Path,
	)
	logger.Info("Entered the deleteHandler")
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}
	logger = logger.With("userID", userID)
	userVal := r.Context().Value("user")
	user, ok := userVal.(database.User)
	if !ok {
		logger.Error("user data not found in contest")
    	utils.RespondWithError(w, http.StatusUnauthorized, "user not in context")
    	return
	}

	if user.Role != "admin" {
		logger.Error("user has no admin status, cannot delete")
    	utils.RespondWithError(w, http.StatusForbidden, "admin access required")
    	return
	}

	bugParam := r.PathValue("bugid")
	if bugParam == "" {
		logger.Error("bug id not found in path")
		utils.RespondWithError(w, http.StatusBadRequest, "no id in the request")
		return
	}
	bugID, err := uuid.Parse(bugParam)
	if err != nil {
		logger.Error("bugID in wrong format", "error", err)
		utils.RespondWithError(w, http.StatusBadRequest, "wrong format id")
		return
	}
	logger = logger.With("bugId", bugID)
	logger.Info("started querying to get existing bug with bugID")
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil { 
		logger.Error("database operation failed", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "no bug  found with the id")
		return

	}
	logger = logger.With("bug", bug)
	if userID != bug.PostedBy {
		logger.Error("user is not the author of the bug")
		utils.RespondWithError(w, http.StatusUnauthorized, "only author can delete the bug")
		return

	}
	logger.Info("starting databse operationto delete bug")
	err = cfg.DB.DeleteBugByID(r.Context(), bugID)
	if err != nil {
		logger.Error("databse operation failed, cannot delete", "error", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot delete bug")
		return
	}
	logger.Info("completed handler ")
	w.WriteHeader(http.StatusNoContent)
}


