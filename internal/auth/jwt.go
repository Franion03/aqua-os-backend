package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ClaimsKey  contextKey = "claims"
	RoleKey    contextKey = "role"
	AdminRole  string     = "Admin"
)

var (
	jwtKey    string
	jwtIssuer string
	jwtAud    string
)

func init() {
	jwtKey = os.Getenv("JWT_KEY")
	if jwtKey == "" {
		jwtKey = "CHANGE_ME_TO_A_LONG_RANDOM_SECRET_KEY_AT_LEAST_32_CHARS"
	}
	jwtIssuer = os.Getenv("JWT_ISSUER")
	if jwtIssuer == "" {
		jwtIssuer = "AquaOsCalendar"
	}
	jwtAud = os.Getenv("JWT_AUDIENCE")
	if jwtAud == "" {
		jwtAud = "AquaOsCalendar"
	}
}

// GenerateToken creates a JWT for the given role. Used by the login endpoint.
func GenerateToken(role string) (string, error) {
	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

// Middleware validates the JWT from the Authorization header and injects claims into context.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtKey), nil
		},
			jwt.WithIssuer(jwtIssuer),
			jwt.WithAudience(jwtAud),
		)
		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		if role, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, RoleKey, role)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns middleware that enforces a specific role.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxRole, _ := r.Context().Value(RoleKey).(string)
			if ctxRole != role {
				http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
