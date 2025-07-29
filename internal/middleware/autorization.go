package middleware

import (
	"net/http"
	"strings"

	"github.com/blacktag/bugby-Go/internal/utils"
	"github.com/casbin/casbin/v2"
)


func Authorization(enforcer *casbin.Enforcer) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			roleValue := r.Context().Value("role")
			role, ok :=roleValue.(string)
			if !ok {
				utils.RespondWithError(w, http.StatusUnauthorized, "missing or invalid role in context")
				return
			}
			obj := r.URL.Path
			act := r.Method

			allowed, err := enforcer.Enforce(role, obj, strings.ToLower(act))
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "error during enforcement")
				return
			}
			if !allowed {
				utils.RespondWithError(w, http.StatusForbidden,"Access denied")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}