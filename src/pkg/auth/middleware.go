// middleware.go

package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsKey struct{}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		fmt.Printf("token: %v\n", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secretKey := os.Getenv("JWT_SECRET_KEY")
			return []byte(secretKey), nil
		})

		if claims := token.Claims.(jwt.MapClaims); err != nil || !token.Valid {
			http.Error(w, "Forbidden", http.StatusForbidden)
			fmt.Println("invalid token")
			return
		} else {
			fmt.Printf("customer_id: %v\n", int64(claims["customer_id"].(float64)))
			if adminID, ok := claims["admin_id"]; ok {
				fmt.Printf("admin_id: %v\n", adminID)
			}
			fmt.Printf("type: %v\n", claims["type"])
			fmt.Printf("exp: %v\n", int64(claims["exp"].(float64)))
		}

		// クレームをコンテキストに格納
		ctx := context.WithValue(r.Context(), ClaimsKey{}, token.Claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
