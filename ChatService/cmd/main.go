package main

import (
	"ChatService/config"
	"ChatService/internal/handlers"
	"fmt"
	"log"
	"net/http"
	"time"
)

func websocketWithCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
}

func main() {
	http.HandleFunc("/ws/chat", websocketWithCORS(handlers.WsHandler))

	http.HandleFunc("/chat/history", withCORS(handlers.ChatHistoryHandler))
	http.HandleFunc("/chats", withCORS(handlers.ChatsUserHandler))
	http.HandleFunc("/reviews", withCORS(handlers.ReviewsHandler))
	http.HandleFunc("/api/ratings/average", withCORS(handlers.AverageRatingHandler))

	server := &http.Server{
		Addr:              config.ServiceConfig.CHAT_SERVICE_URL,
		ReadHeaderTimeout: 3 * time.Second,
	}
	fmt.Println("Сервер запущен на :50055")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
