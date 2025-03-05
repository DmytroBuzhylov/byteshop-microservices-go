package handlers

import (
	"ChatService/config"
	"ChatService/internal/services"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var clients = make(map[string]*websocket.Conn)
var clientsMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type SendMessage struct {
	Participants []string `json:"participants"`
	Message      Message  `json:"message"`
}

func ChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	otherUserID := r.URL.Query().Get("otherUserID")

	chat, err := services.GetUsersHistoryChat(userID, otherUserID)
	if err != nil {
		http.Error(w, "Ошибка при получении чата", http.StatusInternalServerError)
		log.Println("Ошибка при получении чата:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(chat)
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		log.Println("Ошибка кодирования JSON:", err)
	}

}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка обновления до WebSocket:", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("UserID is missing")
		return
	}

	clientsMutex.Lock()
	clients[userID] = conn
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		delete(clients, userID)
		clientsMutex.Unlock()
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}
		var receivedMessage SendMessage
		if err := json.Unmarshal(msgBytes, &receivedMessage); err != nil {
			log.Println("JSON decode error:", err)
			continue
		}

		chat, err := services.FindOrCreateChat(receivedMessage.Participants[0], receivedMessage.Participants[1])
		if err != nil {
			log.Println("Ошибка чата:", err)
			continue
		}

		err = services.SendMessage(chat.Hex(), receivedMessage.Message.Sender, receivedMessage.Message.Content)
		if err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			return
		}

		for _, participant := range receivedMessage.Participants {
			clientsMutex.Lock()
			clientConn, ok := clients[participant]
			clientsMutex.Unlock()

			if ok {
				response := map[string]interface{}{
					"message": map[string]interface{}{
						"sender":    receivedMessage.Message.Sender,
						"content":   receivedMessage.Message.Content,
						"timestamp": time.Now().Format(time.RFC3339),
					},
				}
				respBytes, _ := json.Marshal(response)
				if err := clientConn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
					log.Println("Ошибка отправки сообщения:", err)
				}
			}
		}
	}

}

func ChatsUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	chats, err := services.GetUserChats(userID)
	if err != nil {
		http.Error(w, "Ошибка при получении чата", http.StatusInternalServerError)
		log.Println("Ошибка при получении чата:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		log.Println("Ошибка кодирования JSON:", err)
	}
}

type reviewData struct {
	ProductID string    `json:"productID"`
	UserID    string    `json:"userID"`
	Grade     int       `json:"grade"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var rew reviewData
		if err = json.Unmarshal(body, &rew); err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
		}

		err = services.FindOrCreateProductReviews(rew.ProductID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		userName := GetName(rew.UserID)
		if userName == "" {
			http.Error(w, "get name error", http.StatusInternalServerError)
		}

		err = services.SendReview(rew.ProductID, userName, rew.Content, rew.Grade)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Данные получены",
		})
	}

	if r.Method == http.MethodGet {
		productID := r.URL.Query().Get("productId")

		reviews, err := services.FindProductReviews(productID)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&reviews)

	}

}

func GetName(userId string) string {
	url := fmt.Sprintf("http://%s/api/name/%s", config.ServiceConfig.AUTH_SERVICE_URL, userId)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var name string
	if err := json.NewDecoder(resp.Body).Decode(&name); err != nil {
		return ""
	}
	return name
}

func AverageRatingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	products, err := GetAllProducts(userID)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	average, err := services.GetAverageRating(products)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&average)
}

func GetAllProducts(userID string) ([]string, error) {
	url := fmt.Sprintf("http://%s/api/give-all-productID/%s", config.ServiceConfig.PRODUCT_SERVICE_URL, userID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка: сервер вернул статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var products []string
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, err
	}

	return products, nil
}
