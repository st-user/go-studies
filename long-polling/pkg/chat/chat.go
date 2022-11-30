package chat

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type KeyClientID uuid.UUID

func (kc KeyClientID) String() string {
	return uuid.UUID(kc).String()
}

type KeyRoomID string

type ClientIDSet struct {
	innerMap map[KeyClientID]bool
}

func NewClientIDSet() ClientIDSet {
	return ClientIDSet{
		innerMap: make(map[KeyClientID]bool),
	}
}

func (cd ClientIDSet) add(id KeyClientID) {
	cd.innerMap[id] = true
}

func (cd ClientIDSet) keys() []KeyClientID {
	result := []KeyClientID{}
	for k := range cd.innerMap {
		k := k
		result = append(result, k)
	}
	return result
}

type ClientIDNotFoundError struct {
	clientID KeyClientID
}

func (c ClientIDNotFoundError) Error() string {
	return fmt.Sprintf("clientId=%s was not found", c.clientID)
}

type ChatRoom struct {
	mutex             *sync.Mutex
	clientIdToRoomID  map[KeyClientID]KeyRoomID
	roomIDToClientIDs map[KeyRoomID]ClientIDSet
	clientIdToChan    map[KeyClientID]chan string
}

func NewChatRoom() ChatRoom {
	var mutex sync.Mutex
	return ChatRoom{
		mutex:             &mutex,
		clientIdToRoomID:  make(map[KeyClientID]KeyRoomID),
		roomIDToClientIDs: make(map[KeyRoomID]ClientIDSet),
		clientIdToChan:    make(map[KeyClientID]chan string),
	}
}

func (cr ChatRoom) Join(roomID KeyRoomID) (KeyClientID, error) {

	clientIdContent, err := uuid.NewRandom()
	if err != nil {
		return KeyClientID(uuid.Nil), fmt.Errorf("error on joining a room=%s %w", roomID, err)
	}
	clientId := KeyClientID(clientIdContent)

	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cr.clientIdToRoomID[clientId] = roomID
	clientIDs, ok := cr.roomIDToClientIDs[roomID]
	if !ok {
		fmt.Printf("Room = %s is created.\n", roomID)
		clientIDs = NewClientIDSet()
		cr.roomIDToClientIDs[roomID] = clientIDs
	}
	clientIDs.add(clientId)
	ch := make(chan string)
	cr.clientIdToChan[clientId] = ch

	return clientId, nil
}

func (cr ChatRoom) SendMessage(clientID KeyClientID, message string) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	roomID, ok := cr.clientIdToRoomID[clientID]
	if !ok {
		return ClientIDNotFoundError{clientID}
	}

	for _, eachClientID := range cr.roomIDToClientIDs[roomID].keys() {

		if eachClientID == clientID {
			continue
		}

		ch := cr.clientIdToChan[eachClientID]

		go func(ch chan string, msg string) {
			ch <- message
		}(ch, message)
	}

	return nil
}

func (cr ChatRoom) ReceiveMessage(clientID KeyClientID) (string, error) {
	cr.mutex.Lock()
	ch, ok := cr.clientIdToChan[clientID]
	cr.mutex.Unlock()

	if !ok {
		return "", ClientIDNotFoundError{clientID}
	}

	select {
	case <-time.NewTimer(5 * time.Second).C:
		return "", nil
	case message := <-ch:
		return message, nil
	}
}
