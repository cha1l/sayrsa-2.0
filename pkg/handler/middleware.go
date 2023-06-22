package handler

import (
	"context"
	"net/http"
	"strings"
)

func (h *Handler) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type usernameKey struct{}

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
			NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		rcopy := r.WithContext(context.WithValue(r.Context(), usernameKey{}, username))

		next.ServeHTTP(w, rcopy)
	})
}

func GetParams(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	username, ok := ctx.Value(usernameKey{}).(string)
	if ok {
		return username
	}

	return ""
}
