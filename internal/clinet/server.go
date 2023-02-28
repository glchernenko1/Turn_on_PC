package clinet

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func StartWSServer(host string, JWT string, name string, path string) {
	urlWS := url.URL{Scheme: "WS", Host: host, Path: "ws"}

	con, resp, err := websocket.DefaultDialer.Dial(urlWS.String(),
		http.Header{"Authorization": {JWT}, "name": {name}})

	if resp.StatusCode != 101 {
		message, _ := io.ReadAll(resp.Body)
		panic(string(message))
	}
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	defer con.Close()
	users, err := ReadJsonUser(path)
	for {
		msgType, msg, err := con.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			panic(err)
		}
		comand := strings.Split(string(msg), " ")
		switch comand[0] {
		case "get_list_PCUser_by_Server":
			{
				go getListPCUserByServer(con, users)
			}
		case "turnOn":
			{
				turnOn(users, comand[1])
			}
		default:
			fmt.Printf("Comand: \"%s\" not found \n", comand)
		}
	}
}

func turnOn(users map[string]string, name string) {
	user, ok := users[name]
	if !ok {
		fmt.Printf("User: \"%s\" not found\n", name)
	} else {
		macAddr, err := net.ParseMAC(user)
		if err != nil {
			panic(err)
		}

		// Формируем пакет пробуждения
		packet := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
		for i := 0; i < 16; i++ {
			packet = append(packet, macAddr...)
		}

		addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
		if err != nil {
			panic(err)
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		_, err = conn.Write(packet)
		if err != nil {
			panic(err)
		}
	}

}

func getListPCUserByServer(conn *websocket.Conn, users map[string]string) {
	message := "list_PCUser_by_Server"
	for key, _ := range users {
		message += fmt.Sprintf(" %s", key)
	}
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func ReadJsonUser(path string) (map[string]string, error) {
	user := make(map[string]string)
	file, err := os.ReadFile(path)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(file), &user)
	if err != nil {
		return user, err
	}
	if len(user) == 0 {
		panic("User is empty or not correctly JSON")
	}
	return user, err
}
