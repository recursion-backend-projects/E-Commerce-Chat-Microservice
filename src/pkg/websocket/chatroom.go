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
var mu sync.Mutex

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
	} else if exists {
		if connectionType == "create" {
			chatRoom.customer = client
			log.Printf("Customer %d reconnected to chat room", client.customerID)
			return true
		} else if connectionType == "join" {
			if chatRoom.admin == nil && client.adminID != 0 {
				chatRoom.admin = client
				log.Printf("Admin %d joined chat room for customer %d", client.adminID, client.customerID)
				return true
			}
			log.Printf("Admin %d is already connected to chat room for customer %d or invalid adminID", client.adminID, client.customerID)
		}
	}
	return false
}

// チャットルームからクライアントを削除する
func RemoveClient(client *Client, connectionType string) {
	mu.Lock()
	defer mu.Unlock()

	chatRoom, exists := chatRooms[client.customerID]
	if exists {
		if connectionType == "create" {
			chatRoom.customer = nil
			log.Printf("Customer %d disconnected", client.customerID)
		} else if connectionType == "join" {
			chatRoom.admin = nil
			log.Printf("Admin %d disconnected from chat room for customer %d", client.adminID, client.customerID)
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
	if client.adminID == 0 && chatRoom.admin != nil {
		targetConn = chatRoom.admin.conn
	} else if client.adminID != 0 && chatRoom.customer != nil {
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
