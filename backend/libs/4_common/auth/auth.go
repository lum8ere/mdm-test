package auth

import (
	"net/http"
	"strings"

	"mdm/libs/4_common/smart_context"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var JWTSecret = []byte("your-secret") // по умолчанию, замените на значение из ENV

// SetJWTSecret позволяет установить секрет из внешнего источника (например, из переменной окружения).
func SetJWTSecret(secret string) {
	JWTSecret = []byte(secret)
}

// JWTMiddleware проверяет заголовок Authorization и валидирует JWT-токен.
func JWTMiddleware(next http.Handler, sctx smart_context.ISmartContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			sctx.GetLogger().Error("Invalid JWT token", zap.Error(err))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// При необходимости можно извлечь claims и положить их в контекст запроса.
		next.ServeHTTP(w, r)
	})
}
