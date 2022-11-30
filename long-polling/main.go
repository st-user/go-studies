package main

import (
	"long-polling/pkg/chat"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type JoinRequest struct {
	RoomId string `json:"roomID"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

func main() {

	e := echo.New()

	chatRoom := chat.NewChatRoom()
	e.POST("/enter", func(c echo.Context) error {
		body := new(JoinRequest)
		if err := c.Bind(body); err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Invalid Request Body")
		}
		roomID := chat.KeyRoomID(body.RoomId)
		clientId, err := chatRoom.Enter(roomID)

		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return c.String(http.StatusOK, clientId.String())
	})

	e.POST("/message", func(c echo.Context) error {
		paramClientID, err := uuid.Parse(c.QueryParam("client_id"))
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Invalid client_id")
		}
		body := new(MessageRequest)
		if err := c.Bind(body); err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Invalid Request Body")
		}

		clientId := chat.KeyClientID(paramClientID)
		err = chatRoom.SendMessage(clientId, body.Message)

		if handled, err := handleError(c, err); handled {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.GET("/message", func(c echo.Context) error {
		paramClientID, err := uuid.Parse(c.QueryParam("client_id"))
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Invalid client_id")
		}

		clientId := chat.KeyClientID(paramClientID)
		message, err := chatRoom.ReceiveMessage(clientId)

		if handled, err := handleError(c, err); handled {
			return err
		}

		return c.String(http.StatusOK, message)
	})

	e.DELETE("/leave", func(c echo.Context) error {
		paramClientID, err := uuid.Parse(c.QueryParam("client_id"))
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Invalid client_id")
		}
		clientId := chat.KeyClientID(paramClientID)
		chatRoom.Leave(clientId)

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func handleError(c echo.Context, err error) (bool, error) {
	if err, ok := err.(chat.ClientIDNotFoundError); ok {
		log.Error(err)
		return true, c.String(http.StatusNotFound, err.Error())
	}
	if err != nil {
		log.Error(err)
		return true, c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return false, nil
}
