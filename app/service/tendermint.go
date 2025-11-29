package service

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	websocketPath = "/websocket"
)

type MessageHandler func([]byte)

type TendermintService struct {
	rpcUrl string
}

func NewTendermintService(rpcUrl string) (*TendermintService, error) {
	if rpcUrl == "" {
		return nil, fmt.Errorf("invalid dependencies provided to TendermintService")
	}

	return &TendermintService{rpcUrl: rpcUrl}, nil
}

func (s TendermintService) ListenForArticles(handler MessageHandler) error {
	// Define the Tendermint node's WebSocket URL
	u := url.URL{Scheme: "wss", Host: s.rpcUrl, Path: websocketPath}

	// Connect to the WebSocket
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	defer c.Close()

	// Subscribe to new block events
	err = c.WriteJSON(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      0,
		"params": map[string]string{
			"query": "tm.event = 'Tx' AND message.action='/bze.cointrunk.MsgAddArticle'",
		},
	})
	if err != nil {
		log.Fatal("write:", err)
	}

	// Handle incoming messages
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return err
		}
		handler(message)
	}

	return nil
}
