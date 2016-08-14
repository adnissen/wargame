package gameclient

import (
	"encoding/json"

	"github.com/adnissen/go.uuid"
	"github.com/adnissen/wargame/src/packages/army"
	"github.com/adnissen/wargame/src/packages/userpkg"
	"github.com/adnissen/websocket"
)

type GameClient struct {
	wbs         *websocket.Conn
	Army        army.Army
	CurrentGame uuid.UUID
	User        *userpkg.User
}

type GenericResponse struct {
	MessageType string
	Message     string
}

func (g *GameClient) LoggedIn() bool {
	return g.User != nil
}

func (g *GameClient) SendMessage(m []byte) {
	g.wbs.WriteMessage(1, m)
}

func (g *GameClient) SendMessageOfType(mt string, m []byte) {
	gr := GenericResponse{MessageType: mt, Message: string(m)}
	j, _ := json.Marshal(gr)
	g.wbs.WriteMessage(1, j)
}

func (g *GameClient) CompareWebSocketConn(c *websocket.Conn) bool {
	return c == g.wbs
}

func (g *GameClient) IsStillConnected() bool {
	return g.wbs != nil
}

func (g *GameClient) SetCurrentGame(uid uuid.UUID) {
	g.CurrentGame = uid
}

func (g *GameClient) ResetCurrentGame() {
	g.CurrentGame = uuid.UUID{}
}

func CreateGameClient(c *websocket.Conn) *GameClient {
	ret := GameClient{}
	ret.wbs = c
	ret.ResetCurrentGame()
	return &ret
}
