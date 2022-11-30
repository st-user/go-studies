package chat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	wait_message_timeout_seconds = 5 * time.Second
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

func (cd ClientIDSet) remove(id KeyClientID) {
	delete(cd.innerMap, id)
}

// forEach
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
	clientIdToChan    map[KeyClientID]messageChannelHandler
}

type messageChannelHandler struct {
	mutex     *sync.Mutex
	messageCh chan string
	leaveCh   chan struct{}
	isClosed  *bool
}

func newMessageChannelHandler() messageChannelHandler {
	var m sync.Mutex
	isClosed := false
	messageCh := make(chan string)
	leaveCh := make(chan struct{})
	return messageChannelHandler{
		mutex:     &m,
		messageCh: messageCh,
		leaveCh:   leaveCh,
		isClosed:  &isClosed,
	}
}

func (m messageChannelHandler) sendWithSuccess(msg string) bool {
	m.mutex.Lock()
	shouldNotOp := *m.isClosed
	m.mutex.Unlock()
	if shouldNotOp {
		return false
	}

	select {
	case m.messageCh <- msg:
		return true
	case <-m.leaveCh:
		close(m.messageCh)
		return false
	}
}

func (m messageChannelHandler) receiveWithSuccess() (string, bool) {
	m.mutex.Lock()
	shouldNotOp := *m.isClosed
	m.mutex.Unlock()
	if shouldNotOp {
		return "", false
	}

	select {
	case <-time.NewTimer(wait_message_timeout_seconds).C:
		return "", true
	case msg := <-m.messageCh:
		return msg, true
	case <-m.leaveCh:
		return "", false
	}
}

func (m messageChannelHandler) stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	*m.isClosed = true
	close(m.leaveCh)
}

func NewChatRoom() ChatRoom {
	var mutex sync.Mutex
	return ChatRoom{
		mutex:             &mutex,
		clientIdToRoomID:  make(map[KeyClientID]KeyRoomID),
		roomIDToClientIDs: make(map[KeyRoomID]ClientIDSet),
		clientIdToChan:    make(map[KeyClientID]messageChannelHandler),
	}
}

func (cr ChatRoom) Enter(roomID KeyRoomID) (KeyClientID, error) {

	clientIdContent, err := uuid.NewRandom()
	if err != nil {
		return KeyClientID(uuid.Nil), fmt.Errorf("error on entering a room=%s %w", roomID, err)
	}
	clientId := KeyClientID(clientIdContent)

	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	if _, ok := cr.clientIdToRoomID[clientId]; ok {
		return KeyClientID(uuid.Nil), fmt.Errorf("this client has already entered a room")
	}

	cr.clientIdToRoomID[clientId] = roomID
	clientIDs, ok := cr.roomIDToClientIDs[roomID]
	if !ok {
		fmt.Printf("Room = %s is created.\n", roomID)
		clientIDs = NewClientIDSet()
		cr.roomIDToClientIDs[roomID] = clientIDs
	}
	clientIDs.add(clientId)
	cr.clientIdToChan[clientId] = newMessageChannelHandler()

	return clientId, nil
}

func (cr ChatRoom) SendMessage(clientID KeyClientID, message string) error {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	roomID, ok := cr.clientIdToRoomID[clientID]
	if !ok {
		return ClientIDNotFoundError{clientID}
	}

	clientCount := 0
	for _, eachClientID := range cr.roomIDToClientIDs[roomID].keys() {

		if eachClientID == clientID {
			continue
		}
		clientCount++

		handler := cr.clientIdToChan[eachClientID]

		go func(handler messageChannelHandler, msg string, eachClientID KeyClientID) {
			if !handler.sendWithSuccess(msg) {
				cr.removeWithLock(eachClientID)
			}
		}(handler, message, eachClientID)
	}
	fmt.Printf("the message has been sent to %d clients\n", clientCount)

	return nil
}

func (cr ChatRoom) ReceiveMessage(clientID KeyClientID) (string, error) {
	cr.mutex.Lock()
	handler, ok := cr.clientIdToChan[clientID]
	cr.mutex.Unlock()

	if !ok {
		return "", ClientIDNotFoundError{clientID}
	}

	if msg, ok := handler.receiveWithSuccess(); !ok {
		return "", errors.New("this client has left the chat room")
	} else {
		return msg, nil
	}

}

func (cr ChatRoom) Leave(clientID KeyClientID) {
	cr.mutex.Lock()
	handler, ok := cr.clientIdToChan[clientID]
	cr.mutex.Unlock()

	if !ok {
		return
	}
	cr.removeWithLock(clientID)
	handler.stop()
}

func (cr ChatRoom) removeWithLock(clientID KeyClientID) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.removeWithoutLock(clientID)
}

func (cr ChatRoom) removeWithoutLock(clientID KeyClientID) {
	if roomID, ok := cr.clientIdToRoomID[clientID]; ok {
		if clientIDs, ok := cr.roomIDToClientIDs[roomID]; ok {
			clientIDs.remove(clientID)
		}
	}
	delete(cr.clientIdToRoomID, clientID)
	delete(cr.clientIdToChan, clientID)
}
