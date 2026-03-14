package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Definindo uma chave de contexto para evitar colisões
type contextKey string

const UserIDKey contextKey = "userID"
const TenantIDKey contextKey = "tenantID"

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// AuthMiddleware é o middleware que valida o token JWT.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Cabeçalho de autorização ausente", http.StatusUnauthorized)
			return
		}

		// O cabeçalho deve estar no formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Cabeçalho de autorização mal formatado", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse e validação do token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token com claims inválidas", http.StatusUnauthorized)
			return
		}

		// Extrai o ID do usuário e o ID do tenant e os adiciona ao contexto da requisição
		userID := claims["sub"].(string)
		tenantID := claims["tid"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, TenantIDKey, tenantID)

		// Chama o próximo handler com o contexto modificado
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
