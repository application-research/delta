package core

type WebsocketService struct {
	DeltaNode *DeltaNode
}

// `NewWebsocketService` creates a new `WebsocketService` struct and returns a pointer to it
func NewWebsocketService(dn *DeltaNode) *WebsocketService {
	return &WebsocketService{
		DeltaNode: dn,
	}
}

// HandlePieceCommitmentMessages A function that is listening to the channel `ws.DeltaNode.WebsocketBroadcast.PieceCommitmentChannel.Channel` and when it
// receives a message, it broadcasts it to all clients.
func (ws *WebsocketService) HandlePieceCommitmentMessages() error {
	for {
		message := <-ws.DeltaNode.WebsocketBroadcast.PieceCommitmentChannel.Channel
		// Broadcast to all clients
		for client := range ws.DeltaNode.WebsocketBroadcast.PieceCommitmentChannel.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(ws.DeltaNode.WebsocketBroadcast.PieceCommitmentChannel.Clients, client)
			}
		}
	}
	return nil
}

// HandleContentDealMessages Listening to the channel `ws.DeltaNode.WebsocketBroadcast.ContentDealChannel.Channel` and when it
// // receives a message, it broadcasts it to all clients.
func (ws *WebsocketService) HandleContentDealMessages() error {
	for {
		message := <-ws.DeltaNode.WebsocketBroadcast.ContentDealChannel.Channel

		// Broadcast to all clients
		for client := range ws.DeltaNode.WebsocketBroadcast.ContentDealChannel.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(ws.DeltaNode.WebsocketBroadcast.ContentDealChannel.Clients, client)
			}
		}
	}
	return nil
}

// HandleContentMessages Listening to the channel `ws.DeltaNode.WebsocketBroadcast.ContentChannel.Channel` and when it
// // receives a message, it broadcasts it to all clients.
func (ws *WebsocketService) HandleContentMessages() error {
	for {
		message := <-ws.DeltaNode.WebsocketBroadcast.ContentChannel.Channel
		// Broadcast to all clients
		for client := range ws.DeltaNode.WebsocketBroadcast.ContentChannel.Clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(ws.DeltaNode.WebsocketBroadcast.ContentChannel.Clients, client)
			}
		}
	}
	return nil
}
