package api

import (
	"delta/core"
	"fmt"
	model "delta/models"
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
		Clients: make(map[*core.ClientChannel]bool),
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

	ln.DeltaEventEmitter.WebsocketBroadcast.ContentChannel = contentChannel
	ln.DeltaEventEmitter.WebsocketBroadcast.PieceCommitmentChannel = pieceCommitmentChannel
	ln.DeltaEventEmitter.WebsocketBroadcast.ContentDealChannel = contentDealChannel

	wsGroup.GET("/contents/:contentId", handleWebsocketContent(ln))
	wsGroup.GET("/piece-commitments/:pieceCommitmentId", handleWebsocketPieceCommitment(ln))
	wsGroup.GET("/deals/by-uuid/:dealUuid", handleWebsocketContentDeal(ln))
	wsGroup.GET("/deals/by-cid/:cid", handleWebsocketContentDeal(ln))

}

// It upgrades the HTTP connection to a WebSocket connection, registers the new client, and then reads messages from the
// client and sends them to the broadcast channel
func handleWebsocketContent(ln *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		// Upgrade HTTP connection to WebSocket connection
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("WebSocket upgrade error:", err)
			return nil
		}

		contentId := c.Param("contentId")

		// Register new client
		mutex.Lock()
		clientChannelForContent := &core.ClientChannel{
			Conn: conn,
			Id:   contentId,
		}
		ln.DeltaEventEmitter.WebsocketBroadcast.ContentChannel.Clients[clientChannelForContent] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.DeltaEventEmitter.WebsocketBroadcast.ContentChannel.Clients, clientChannelForContent)
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
			ln.DeltaEventEmitter.WebsocketBroadcast.ContentChannel.Channel <- msg
		}

		return nil
	}
}

// It upgrades the HTTP connection to a WebSocket connection, registers the new client, and then reads messages from the
// client and sends them to the broadcast channel
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
		ln.DeltaEventEmitter.WebsocketBroadcast.PieceCommitmentChannel.Clients[conn] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.DeltaEventEmitter.WebsocketBroadcast.PieceCommitmentChannel.Clients, conn)
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
			ln.DeltaEventEmitter.WebsocketBroadcast.PieceCommitmentChannel.Channel <- msg
		}

		return nil
	}
}

// It upgrades the HTTP connection to a WebSocket connection, registers the new client, and then reads messages from the
// client and sends them to the broadcast channel
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
		ln.DeltaEventEmitter.WebsocketBroadcast.ContentDealChannel.Clients[conn] = true
		mutex.Unlock()

		// Close WebSocket connection when client disconnects
		defer func() {
			mutex.Lock()
			delete(ln.DeltaEventEmitter.WebsocketBroadcast.ContentDealChannel.Clients, conn)
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
			ln.DeltaEventEmitter.WebsocketBroadcast.ContentDealChannel.Channel <- msg
		}

		return nil
	}
}
