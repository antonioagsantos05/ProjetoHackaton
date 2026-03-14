package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware valida se existe o Bearer Token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Token missing", http.StatusUnauthorized)
			return
		}

		// Validação de JWT aqui
		// jwt.Parse...
		// Extração de tenant_id, role e user_id
		next.ServeHTTP(w, r)
	})
}
