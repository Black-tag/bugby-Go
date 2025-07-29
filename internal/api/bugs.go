package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/database"
	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/google/uuid"
)




func (cfg *APIConfig) CreateBugHandler (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	userIDValue := r.Context().Value("userID")
	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "cannot find user ID")
		return
	}
	type CreateBugRequest struct {
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		PostedBy     uuid.UUID `json:"posted_by"`
	}

	var req CreateBugRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error decoding json")
		return
	}
	bug, err := cfg.DB.CreateBug(r.Context(),database.CreateBugParams{
		Title: req.Title,
		Description: req.Description,
		PostedBy: userID,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot create bug")
		return
	}
	type CreateBugResponse struct {
		ID           uuid.UUID `json:"bug_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		PostedBy     uuid.UUID `json:"posted_by"`
		CreatedBy    time.Time `json:"created_at"`
		Updated_at   time.Time `json:"updated_at"`

	}
	utils.RespondWithJSON(w, http.StatusCreated, CreateBugResponse{
		ID: bug.ID,
		Title: bug.Title,
		Description: bug.Description,
		PostedBy: bug.PostedBy,
		CreatedBy: bug.CreatedAt,
		Updated_at: bug.UpdatedAt,
	})
	
}

func (cfg *APIConfig) GetBugsHandler (w http.ResponseWriter, r *http.Request) {
	bugs, err := cfg.DB.GetAllBugs(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "couldnt fetch bugs")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, bugs)
}


func (cfg *APIConfig) GetBugByIDHandler (w http.ResponseWriter, r *http.Request) {
	
	bugIDParam := r.PathValue("bugid")

	if bugIDParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, " no id given ")
		return
	}
	bugID, err := uuid.Parse(bugIDParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "wrong format Id")
		return
	}
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, " bug not found ")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, bug)
}


func (cfg *APIConfig) UpadteBugHandler (w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}
	bugParam := r.PathValue("bugid")
	if bugParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, " no id given ")
		return
	}
	bugID, err := uuid.Parse(bugParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "wrong format Id")
		return
	}
	type UpdateBugRequest struct {
		Title       *string    `json:"title"`
		Description *string    `json:"description"`
		
		
	}
	var req UpdateBugRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil { 
		utils.RespondWithError(w, http.StatusInternalServerError, "no bug  found with the id")
		return

	}
	if userID != bug.PostedBy {
		utils.RespondWithError(w, http.StatusUnauthorized, "only author can delete the bug")
		return

	}
	params := database.UpdateBugByIDParams{
		ID: bugID,
		Title: toNullString(req.Title),
		Description: toNullString(req.Description),
		
	}
	
	err = cfg.DB.UpdateBugByID(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot update bug")
		return
	}
	updatedbug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot fetch updated bug")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updatedbug)
}
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func (cfg *APIConfig) DeleteBugByIDHandler (w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "invalid or missing user ID")
		return
	}
	userVal := r.Context().Value("user")
	user, ok := userVal.(database.User)
	if !ok {
    	utils.RespondWithError(w, http.StatusUnauthorized, "user not in context")
    	return
	}

	if user.Role != "admin" {
    	utils.RespondWithError(w, http.StatusForbidden, "admin access required")
    	return
	}

	bugParam := r.PathValue("bugid")
	if bugParam == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "no id in the request")
		return
	}
	bugID, err := uuid.Parse(bugParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "wrong format id")
		return
	}
	bug, err := cfg.DB.GetBugsByID(r.Context(), bugID)
	if err != nil { 
		utils.RespondWithError(w, http.StatusInternalServerError, "no bug  found with the id")
		return

	}
	if userID != bug.PostedBy {
		utils.RespondWithError(w, http.StatusUnauthorized, "only author can delete the bug")
		return

	}
	err = cfg.DB.DeleteBugByID(r.Context(), bugID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "cannot delete bug")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


