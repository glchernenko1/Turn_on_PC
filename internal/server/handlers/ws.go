package handlers

import (
	"Turn_on_PC/internal/DTO"
	"Turn_on_PC/internal/server/middleware"
	"Turn_on_PC/internal/server/servis"
	"Turn_on_PC/pkg/logging"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

const (
	WsUrl = "/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type handlerWs struct {
	logger   *logging.Logger
	dbSocket servis.DBClient
}

func NewHandlerWs(logger *logging.Logger) Handler {
	var dbSocket servis.DBClient
	servis.NewDBClient(&dbSocket)
	return &handlerWs{
		logger:   logger,
		dbSocket: dbSocket,
	}
}

func (h *handlerWs) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, WsUrl, middleware.MiddlewareAuth(h.WS, "WebUser", "ServerUser", "ClientPC"))
}

func (h *handlerWs) WS(w http.ResponseWriter, r *http.Request, token *DTO.JWTUser) error {

	name := r.Header.Get("name") // todo подумать где лучше разместить информацию о имине
	conn, _ := upgrader.Upgrade(w, r, nil)
	client := servis.NewClient(token, name, conn)
	err := h.dbSocket.AddClient(client)
	if err != nil {
		conn.Close()
		return err
	}
	defer h.dbSocket.DeleteClient(client)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			return nil
		}

		if client.Token.Scope == "WebUser" {
			commands := strings.Split(string(msg), " ")
			switch commands[0] {
			case "turnOn":
				{
					go h.dbSocket.MessageToServerUser(client, string(msg))
				}
			case "get_list_ServerUser":
				{
					go client.Socket.WriteJSON(h.dbSocket.GetListServerUser(client))
				}
			case "get_list_PCUserConnection":
				{
					go client.Socket.WriteJSON(h.dbSocket.GetListPCUserConnection(client))
				}
			case "get_list_PCUser_by_Server":
				{
					go h.dbSocket.MessageToServerUser(client, "get_list_PCUser_by_Server")
				}
			case "turnOff":
				{
					go h.dbSocket.TurnOff(client, commands[1])
				}
			default:
				go client.Socket.WriteMessage(websocket.TextMessage, []byte("command not Found"))
			}

		}
		if client.Token.Scope == "ServerUser" {
			commands := strings.Split(string(msg), " ")
			if commands[0] == "list_PCUser_by_Server" {
				go h.dbSocket.ListPCUserByServer(client, commands[1:])
			} else {
				go client.Socket.WriteMessage(websocket.TextMessage, []byte("command not Found"))
			}
		}

	}
}
