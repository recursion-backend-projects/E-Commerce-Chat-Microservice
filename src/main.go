package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Response struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		} else {
			log.Printf("send: %s", message)
		}
	}
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
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

		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Goodbye, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/goodbye", goodbyeHandler)
	http.HandleFunc("/chat", authenticate(echo))

	http.ListenAndServe(":8080", nil)
}
