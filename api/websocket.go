package api

import (
	"delta/core"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"sync"
)

var (
	upgrader = websocket.Upgrader{}
	mutex    = &sync.Mutex{}
)

// ConfigureWebsocketRouter It creates a new websocket handler that will send the content of the node to the client
func ConfigureWebsocketRouter(e *echo.Group, ln *core.DeltaNode) {

	wsGroup := e.Group("/ws")
	// initiate the websocket broadcast
	contentChannel := core.ContentChannel{
		Clients: make(map[*websocket.Conn]bool),
		Channel: make(chan model.Content),
	}

	pieceCommitmentChannel := core.PieceCommitmentChannel{
		Clients: make(map[*websocket.Conn]bool),
		Channel: make(chan model.PieceCommitment),
	}

	contentDealChannel := core.ContentDealChannel{
		Clients: make(map[*websocket.Conn]bool),
		Channel: make(chan model.ContentDeal),
	}

	ln.WebsocketBroadcast.ContentChannel = contentChannel
	ln.WebsocketBroadcast.PieceCommitmentChannel = pieceCommitmentChannel
	ln.WebsocketBroadcast.ContentDealChannel = contentDealChannel

	wsGroup.GET("/contents", handleWebsocketContent(ln))
	wsGroup.GET("/piece-commitments", handleWebsocketPieceCommitment(ln))
	wsGroup.GET("/deals", handleWebsocketContentDeal(ln))

}

// Handling a websocket connection.
func handleWebsocketContent(ln *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Upgrade HTTP connection to WebSocket connection
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return nil
		}

		// Register new client
		mutex.Lock()
		ln.WebsocketBroadcast.ContentChannel.Clients[conn] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.WebsocketBroadcast.ContentChannel.Clients, conn)
			mutex.Unlock()
			conn.Close()
		}()

		for {
			var msg model.Content
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("WebSocket read error:", err)
				break
			}
			// Send message to broadcast channel
			ln.WebsocketBroadcast.ContentChannel.Channel <- msg
		}

		return nil
	}
}

// Handling a websocket connection.
func handleWebsocketPieceCommitment(ln *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Upgrade HTTP connection to WebSocket connection
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return nil
		}

		// Register new client
		mutex.Lock()
		ln.WebsocketBroadcast.PieceCommitmentChannel.Clients[conn] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.WebsocketBroadcast.PieceCommitmentChannel.Clients, conn)
			mutex.Unlock()
			conn.Close()
		}()

		for {
			var msg model.PieceCommitment
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("WebSocket read error:", err)
				break
			}
			// Send message to broadcast channel
			ln.WebsocketBroadcast.PieceCommitmentChannel.Channel <- msg
		}

		return nil
	}
}

// Handling a websocket connection.
func handleWebsocketContentDeal(ln *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Upgrade HTTP connection to WebSocket connection
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return nil
		}

		// Register new client
		mutex.Lock()
		ln.WebsocketBroadcast.ContentDealChannel.Clients[conn] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.WebsocketBroadcast.ContentDealChannel.Clients, conn)
			mutex.Unlock()
			conn.Close()
		}()

		for {
			var msg model.ContentDeal
			err := conn.ReadJSON(&msg)
			if err != nil {
				fmt.Println("WebSocket read error:", err)
				break
			}
			// Send message to broadcast channel
			ln.WebsocketBroadcast.ContentDealChannel.Channel <- msg
		}

		return nil
	}
}
