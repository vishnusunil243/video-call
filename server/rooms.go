package server

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.Mutex
	Map   map[string][]Participant
}

// init initialises the roomMap
func (r *RoomMap) Init() {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	r.Map = make(map[string][]Participant)
}
func (r *RoomMap) Get(roomId string) []Participant {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	return r.Map[roomId]
}

// create a unique roomId return it and insert in hashmap
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	roomId := string(b)
	r.Map[roomId] = []Participant{}
	return roomId
}
func (r *RoomMap) InsertIntoRoom(roomId string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	p := Participant{host, conn}
	log.Println("inserting into toom with roomid: ", roomId)
	r.Map[roomId] = append(r.Map[roomId], p)
}
func (r *RoomMap) DeleteRoom(roomId string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Map, roomId)
}
