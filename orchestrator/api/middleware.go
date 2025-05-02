package api

import (
    "context"
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/unethiqual/CALC_PROJ/config"
)

type contextKey string

const userIDKey contextKey = "userID"

func generateJWT(userID int) (string, error) {
    cfg := config.LoadConfig()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": userID,
    })
    return token.SignedString([]byte(cfg.JWTSecret))
}

func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        cfg := config.LoadConfig()
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, http.ErrAbortHandler
            }
            return []byte(cfg.JWTSecret), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        userID, ok := claims["userID"].(float64)
        if !ok {
            http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, int(userID))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}