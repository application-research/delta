package core

import model "github.com/application-research/delta-db/db_models"

type WebsocketService struct {
	DeltaNode *DeltaNode
}

func NewWebsocketService(dn *DeltaNode) *WebsocketService {
	return &WebsocketService{
		DeltaNode: dn,
	}
}

func (ws *WebsocketService) SendContentMessage(content model.Content) error {
	ws.DeltaNode.Websocket.WriteMessage(1, []byte("Hello, Client!"))
	return nil
}

func (ws *WebsocketService) SendMessage(message string) error {
	ws.DeltaNode.Websocket.WriteMessage(1, []byte(message))
	return nil
}
