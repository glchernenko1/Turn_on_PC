package servis

import (
	"github.com/gorilla/websocket"
	"Turn_on_PC/internal/DTO"
	"fmt"
	"Turn_on_PC/internal/server/apperror"
)

type Client struct {
	Name   string
	Token  *DTO.JWTUser
	Socket *websocket.Conn
}

func (c *Client) Close() error {
	return c.Socket.Close()
}

func NewClient(token *DTO.JWTUser, name string, socket *websocket.Conn) Client {
	return Client{Name: name, Token: token, Socket: socket}
}

type DBClient struct {
	WebUser    map[uint]map[Client]bool
	ServerUser map[uint]map[string]Client
	ClientPC   map[uint]map[string]Client
}

func NewDBClient(db *DBClient) {
	if db.WebUser == nil {
		db.WebUser = make(map[uint]map[Client]bool)
	}
	if db.ServerUser == nil {
		db.ServerUser = make(map[uint]map[string]Client)
	}
	if db.ClientPC == nil {
		db.ClientPC = make(map[uint]map[string]Client)
	}
}

func (DB *DBClient) AddClient(client Client) error {

	switch client.Token.Scope {
	case "WebUser":
		{
			id, ok := DB.WebUser[client.Token.UserId]
			if !ok {
				id = make(map[Client]bool)
				DB.WebUser[client.Token.UserId] = id
			}
			DB.WebUser[client.Token.UserId][client] = true
		}
	case "ServerUser":
		{
			id, ok := DB.ServerUser[client.Token.UserId]
			if !ok {
				id = make(map[string]Client)
				DB.ServerUser[client.Token.UserId] = id
			}
			_, ok2 := DB.ServerUser[client.Token.UserId][client.Name]
			if ok2 {
				return apperror.NameTaken // todo добавить что имя уже занято
			}

			DB.ServerUser[client.Token.UserId][client.Name] = client
			go DB.MessageToWebUser(client, fmt.Sprintf("Server: \"%s\" connected", client.Name))
			go DB.GetListPCUserConnection(client)

		}
	case "ClientPC":
		{
			id, ok := DB.ClientPC[client.Token.UserId]
			if !ok {
				id = make(map[string]Client)
				DB.ClientPC[client.Token.UserId] = id
			}
			_, ok2 := DB.ClientPC[client.Token.UserId][client.Name]
			if ok2 {
				return apperror.NameTaken
			}
			DB.ClientPC[client.Token.UserId][client.Name] = client
			go DB.MessageToWebUser(client, fmt.Sprintf("Client \"%s\" turned on", client.Name))

		}
	default:
		{
			return apperror.BadRequest
		}
	}
	return nil
}

func (DB *DBClient) DeleteClient(client Client) {

	switch client.Token.Scope {
	case "WebUser":
		{
			client.Socket.Close()
			delete(DB.WebUser[client.Token.UserId], client)
			if len(DB.WebUser[client.Token.UserId]) == 0 {
				delete(DB.WebUser, client.Token.UserId)
			}
		}
	case "ServerUser":
		{
			go DB.MessageToWebUser(client, fmt.Sprintf("%s disconnected", client.Name))
			client.Socket.Close()
			delete(DB.ServerUser[client.Token.UserId], client.Name)
			if len(DB.ServerUser[client.Token.UserId]) == 0 {
				delete(DB.ServerUser, client.Token.UserId)
			}
		}
	default:
		{
			go DB.MessageToWebUser(client, fmt.Sprintf("%s turned off", client.Name))
			client.Socket.Close()
			delete(DB.ClientPC[client.Token.UserId], client.Name)
			if len(DB.ClientPC[client.Token.UserId]) == 0 {
				delete(DB.ClientPC, client.Token.UserId)
			}
		}

	}
}

func (DB *DBClient) MessageToWebUser(client Client, message string) {
	for conn := range DB.WebUser[client.Token.UserId] {
		conn.Socket.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (DB *DBClient) ListPCUserByServer(client Client, message []string) {
	for conn := range DB.WebUser[client.Token.UserId] {
		conn.Socket.WriteJSON(message)
	}
}

func (DB *DBClient) MessageToServerUser(client Client, message string) {
	for conn := range DB.ServerUser[client.Token.UserId] { //todo подумать нормально ли рассылать этот запрос всем нашим серверам
		DB.ServerUser[client.Token.UserId][conn].Socket.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (DB *DBClient) GetListServerUser(client Client) []string {

	res := make([]string, len(DB.ServerUser[client.Token.UserId]))
	i := 0
	for name, _ := range DB.ServerUser[client.Token.UserId] {
		res[i] = name
		i++
	}
	return res
}

func (DB *DBClient) GetListPCUserConnection(client Client) []string {

	res := make([]string, len(DB.ClientPC[client.Token.UserId]))
	i := 0
	for name, _ := range DB.ClientPC[client.Token.UserId] {
		res[i] = name
		i++
	}
	return res
}

func (DB *DBClient) TurnOff(client Client, name string) {
	coon, ok := DB.ClientPC[client.Token.UserId][name]
	if !ok {
		client.Socket.WriteMessage(websocket.TextMessage, []byte("not Found"))
	}
	coon.Socket.WriteMessage(websocket.TextMessage, []byte("turnOff"))
}
