package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/golang-jwt/jwt/v5"
	"E-Commerce-Chat-Microservice/pkg/auth"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	// コンテキストからクレームを取得
	claims, ok := r.Context().Value(auth.ClaimsKey{}).(jwt.MapClaims)
	if !ok {
		log.Println("Failed to get claims from context")
		return
	}

	customerID := int64(claims["customer_id"].(float64))
	adminID := int64(0)
	if claims["admin_id"] != nil {
		adminID = int64(claims["admin_id"].(float64))
	}
	connectionType := claims["type"].(string)

	log.Printf("customerID: %d, adminID: %d, connectionType: %s", customerID, adminID, connectionType)

	client := &Client{conn: conn, customerID: customerID, adminID: adminID}

	added := AddClient(client, connectionType)
	if !added {
		log.Printf("Failed to add client to chat room: customerID=%d, adminID=%d, connectionType=%s", customerID, adminID, connectionType)
		return
	}

	defer RemoveClient(client, connectionType)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}
		log.Printf("recv: %s", message)

		RelayMessage(client, messageType, message)
	}
}