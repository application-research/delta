package core

type WebsocketService struct {
	DeltaNode *DeltaNode
}

func NewWebsocketService(dn *DeltaNode) *WebsocketService {
	return &WebsocketService{
		DeltaNode: dn,
	}
}

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
