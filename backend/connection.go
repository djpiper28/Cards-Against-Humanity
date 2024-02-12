package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const wsBufferSize = 1024

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  wsBufferSize,
	WriteBufferSize: wsBufferSize,
}

func WsUpgrade(w http.ResponseWriter, r *http.Request, playerId, gameId uuid.UUID, cm ConnectionManager) (*WsConnection, error) {
	c, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %s", err)
		return nil, err
	}

	conn := cm.NewConnection(c, playerId, gameId)
	return conn, nil
}

type GameMessage struct {
	Message  string
	GameId   uuid.UUID
	PlayerId uuid.UUID
}

type WsConnection struct {
	Conn         NetworkConnection
	PlayerId     uuid.UUID
	GameID       uuid.UUID
	JoinTime     time.Time
	LastPingTime time.Time
	WsRecieve    chan GameMessage
	WsBroadcast  chan string
	shutdown     chan bool
}

func (gcm *GlobalConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *WsConnection {
	c := &WsConnection{Conn: &WebsocketConnection{Conn: conn},
		PlayerId:     playerId,
		GameID:       gameId,
		JoinTime:     time.Now(),
		LastPingTime: time.Now(),
		WsRecieve:    make(chan GameMessage),
		WsBroadcast:  make(chan string),
		shutdown:     make(chan bool),
	}
	go c.Process()

	gcm.RegisterConnection(gameId, playerId, c)
	return c
}

func (c *WsConnection) Process() {
	go func() {
		for {
			select {
			case <-c.shutdown:
				return
			case msg := <-c.WsBroadcast:
				err := c.Conn.Send([]byte(msg))
				if err != nil {
					log.Printf("Player %s had a network error %s", c.PlayerId, err)
					globalConnectionManager.Close(c.GameID, c.PlayerId)
					return
				}
			}
		}
	}()

	for {
		select {
		case <-c.shutdown:
			return
		default:
			msg, err := c.Conn.Receive()
			if err != nil {
				log.Println(err)
				globalConnectionManager.Close(c.GameID, c.PlayerId)
				return
			}

			c.WsRecieve <- GameMessage{Message: string(msg), GameId: c.GameID, PlayerId: c.PlayerId}
		}
	}
}

func (c *WsConnection) Close() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("Closing the websocket caused an issue %s", err)
			}
		}()

		c.Conn.Close()
		c.shutdown <- true

		close(c.WsBroadcast)
		close(c.WsRecieve)
	}()
}
