package handler

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == "" {
			NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		splitToken := strings.Split(reqToken, " ")
		if len(splitToken) != 2 {
			NewErrorResponse(w, http.StatusUnauthorized, "wrong length of header")
			return
		}

		reqToken = splitToken[1]
		username, err := h.service.Authorization.GetUsernameByToken(reqToken)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rcopy := r.WithContext(context.WithValue(r.Context(), "username", username))

		next.ServeHTTP(w, rcopy)
	})
}

func GetParams(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	username, ok := ctx.Value("username").(string)
	if ok {
		return username
	}

	return ""
}
