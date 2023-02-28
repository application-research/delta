package api

import (
	"delta/core"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

// ConfigureWebsocketRouter It creates a new route that listens for GET requests to the `/ws` endpoint and calls the `handleWebsocketConn` function
// when it receives one
func ConfigureWebsocketRouter(e *echo.Group, ln *core.DeltaNode) {
	e.GET("/ws", handleWebsocketConn(ln))
}

// Handling a websocket connection.
func handleWebsocketConn(ln *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		ln.Websocket = ws
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// Write
			err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}
}
