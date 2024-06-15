package main

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
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

func konichiwaHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{Message: "こんにちは、世界！"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func deployTestHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{Message: "Deploy Test Message!"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/goodbye", goodbyeHandler)
    http.HandleFunc("/konichiwa", konichiwaHandler)
    http.HandleFunc("/deployTest", deployTestHandler)
    http.ListenAndServe(":8080", nil)
}
