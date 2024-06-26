package main

import (
	"E-Commerce-Chat-Microservice/pkg/auth"
	"E-Commerce-Chat-Microservice/pkg/websocket"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Print("Error loading .env file")
	}

	http.HandleFunc("/chat", auth.Authenticate(websocket.HandleChat))
	http.ListenAndServe(os.Getenv("PORT"), nil)
}
