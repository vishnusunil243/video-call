// package server

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var AllRooms RoomMap

// type BroadCasting struct {
// 	Message map[string]interface{}
// 	RoomId  string
// 	Client  *websocket.Conn
// }

// var broadcast = make(chan BroadCasting)

// func Broadcaster() {
// 	for {
// 		msg := <-broadcast

// 		for _, client := range AllRooms.Map[msg.RoomId] {
// 			if client.Conn != msg.Client {
// 				err := client.Conn.WriteJSON(msg.Message)
// 				if err != nil {
// 					log.Fatal(err)
// 					client.Conn.Close()
// 				}
// 			}
// 		}
// 		log.Println(msg.Message)
// 	}
// }

// // create a room and return roomID
// func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
// 	enableCors(&w)
// 	roomId := AllRooms.CreateRoom()
// 	type resp struct {
// 		RoomId string `json:"room_id"`
// 	}
// 	json.NewEncoder(w).Encode(resp{RoomId: roomId})
// }
// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "*")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
// }

// var upgraded = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// // joinrommrequest handler will join client in a particular room
//
//	func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
//		enableCors(&w)
//		queryParams := r.URL.Query()
//		roomId, ok := queryParams["roomId"]
//		if !ok {
//			log.Println("roomId not found")
//			return
//		}
//		ws, err := upgraded.Upgrade(w, r, nil)
//		if err != nil {
//			log.Println("web socket upgrade error ", err)
//		}
//		AllRooms.InsertIntoRoom(roomId[0], false, ws)
//		defer func() {
//			AllRooms.DeleteRoom(roomId[0])
//		}()
//		for {
//			var msg BroadCasting
//			err := ws.ReadJSON(&msg.Message)
//			if err != nil {
//				log.Fatal(err)
//			}
//			msg.Client = ws
//			msg.RoomId = roomId[0]
//			broadcast <- msg
//		}
//	}
// package server

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var AllRooms RoomMap

// func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	roomID := AllRooms.CreateRoom()
// 	type resp struct {
// 		RoomID string `json:"room_id"`
// 	}
// 	log.Println(AllRooms.Map)
// 	json.NewEncoder(w).Encode(resp{RoomID: roomID})
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// type broadcastMsg struct {
// 	Message map[string]interface{}
// 	RoomID  string
// 	Client  *websocket.Conn
// }

// var broadcast = make(chan broadcastMsg)

// func Broadcaster() {
// 	for {
// 		msg := <-broadcast
// 		for _, client := range AllRooms.Map[msg.RoomID] {
// 			if client.Conn != msg.Client {
// 				err := client.Conn.WriteJSON(msg.Message)
// 				if err != nil {
// 					log.Fatal(err)
// 					client.Conn.Close()
// 				}
// 			}
// 		}
// 	}
// }

//	func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
//		roomID, ok := r.URL.Query()["roomId"]
//		if !ok {
//			log.Println("roomID missing in URL Parameters")
//			return
//		}
//		ws, err := upgrader.Upgrade(w, r, nil)
//		if err != nil {
//			log.Fatal("Web Socket Upgrade Error", err)
//		}
//		AllRooms.InsertIntoRoom(roomID[0], false, ws)
//		for {
//			var msg broadcastMsg
//			err := ws.ReadJSON(&msg.Message)
//			if err != nil {
//				log.Fatal("Read Error: ", err)
//			}
//			msg.Client = ws
//			msg.RoomID = roomID[0]
//			log.Println(msg)
//			broadcast <- msg
//		}
//	}
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap
var mutex = &sync.Mutex{} // Mutex for concurrent access to AllRooms.Map
var broadcast = make(chan broadcastMsg)

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()
	type resp struct {
		RoomID string `json:"room_id"`
	}
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Broadcaster() {
	for {
		msg := <-broadcast

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client && client.Conn != nil {
				// Check if connection is still valid before sending
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Println("Error sending message to client:", err)
					// Consider removing the disconnected client from AllRooms
					go func() {
						mutex.Lock()
						defer mutex.Unlock()
						AllRooms.DeleteRoom(msg.RoomID)
					}()
				}
			}
		}
		log.Println(msg) // Optional: log even if no clients received the message
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomId"]
	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	client := wsClient{Conn: ws, RoomID: roomID[0]} // Create a client struct with connection and room ID

	// Add client to AllRooms with a goroutine to handle disconnections gracefully
	go func() {
		defer ws.Close() // Ensure connection is closed on cleanup
		AllRooms.InsertIntoRoom(roomID[0], false, client.Conn)
		for {
			var msg broadcastMsg
			err := ws.ReadJSON(&msg.Message)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Client disconnected:", client.Conn.RemoteAddr())
				} else {
					log.Println("Read Error:", err)
				}
				break // Exit the loop on error or normal closure
			}
			msg.Client = client.Conn
			msg.RoomID = client.RoomID
			broadcast <- msg
		}
		// Remove client from AllRooms upon disconnection
		mutex.Lock()
		defer mutex.Unlock()
		AllRooms.DeleteRoom(roomID[0])
	}()
}

type wsClient struct {
	*websocket.Conn
	RoomID string
}
