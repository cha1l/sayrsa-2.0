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
			NewErrorResponse(w, http.StatusUnauthorized, "emplty auth header")
			return
		}

		splitToken := strings.Split(reqToken, " ")
		if len(splitToken) != 2 {
			NewErrorResponse(w, http.StatusUnauthorized, "wrong lenght of header")
			return
		}

		reqToken = splitToken[1]
		id, err := h.service.Authorization.GetUserIdByToken(reqToken)
		if err != nil {
			NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rcopy := r.WithContext(context.WithValue(r.Context(), "user_id", id))

		next.ServeHTTP(w, rcopy)
	})
}

func GetParams(ctx context.Context) int {
	if ctx == nil {
		return -1
	}

	id, ok := ctx.Value("user_id").(int)
	if ok {
		return id
	}

	return -1
}
