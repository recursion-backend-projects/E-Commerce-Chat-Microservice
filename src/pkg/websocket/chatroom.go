package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn       *websocket.Conn
	customerID int64
	adminID    int64
}

type ChatRoom struct {
	customer *Client
	admin    *Client
}

// チャットルームを保持するハッシュテーブル
var chatRooms = make(map[int64]*ChatRoom)
var mu sync.Mutex // Protects chatRooms

// チャットルームにクライアントを追加する
func AddClient(client *Client, connectionType string) bool {
	mu.Lock()
	defer mu.Unlock()

	chatRoom, exists := chatRooms[client.customerID]
	if !exists && connectionType == "create" {
		chatRoom = &ChatRoom{customer: client}
		chatRooms[client.customerID] = chatRoom
		log.Printf("Created chat room for customer %d", client.customerID)
		return true
	} else if exists && connectionType == "join" {
		if chatRoom.admin == nil && client.adminID != 0 {
			chatRoom.admin = client
			log.Printf("Admin %d joined chat room for customer %d", client.adminID, client.customerID)
			return true
		}
		log.Printf("Admin %d is already connected to chat room for customer %d or invalid adminID", client.adminID, client.customerID)
		return false
	}
	return false
}

// チャットルームからクライアントを削除する
func RemoveClient(client *Client, connectionType string) {
	mu.Lock()
	defer mu.Unlock()

	chatRoom, exists := chatRooms[client.customerID]
	if !exists {
		return
	}

	if connectionType == "create" {
		delete(chatRooms, client.customerID)
		log.Printf("Deleted chat room for customer %d", client.customerID)
	} else if connectionType == "join" {
		chatRoom.admin = nil
		if chatRoom.customer == nil {
			delete(chatRooms, client.customerID)
			log.Printf("Deleted chat room for customer %d", client.customerID)
		}
	}
}

// チャットルーム内の他のクライアントにメッセージを中継する
func RelayMessage(client *Client, messageType int, message []byte) {
	mu.Lock()
	defer mu.Unlock()

	chatRoom, exists := chatRooms[client.customerID]
	if !exists {
		return
	}

	var targetConn *websocket.Conn
	if client.adminID == 0 && chatRoom.admin != nil { // Customer sending message to admin
		targetConn = chatRoom.admin.conn
	} else if client.adminID != 0 && chatRoom.customer != nil { // Admin sending message to customer
		targetConn = chatRoom.customer.conn
	}

	if targetConn != nil {
		err := targetConn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("write error: %v", err)
		}
		log.Printf("send: %s", message)
	} else {
		log.Println("No target connection available")
	}
}
