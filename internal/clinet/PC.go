package clinet

import (
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

func StartWSPC(host string, JWT string, name string) {
	urlWS := url.URL{Scheme: "WS", Host: host, Path: "ws"}

	con, resp, err := websocket.DefaultDialer.Dial(urlWS.String(),
		http.Header{"Authorization": {JWT}, "name": {name}})

	if resp.StatusCode != 200 {
		message, _ := io.ReadAll(resp.Body)
		panic(string(message))
	}
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	defer con.Close()
	for {
		msgType, msg, err := con.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			panic(err)
		}
		comand := string(msg)
		if comand == "turnOff" {
			tornOnPC()
		} else {
			println("command not found")
		}
	}
}

func tornOnPC() {
	switch runtime.GOOS {
	case "linux":
		{
			cmd := exec.Command("systemctl", "suspend", "-i")
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		}
	case "windows":
		{
			cmd := exec.Command("powercfg", "/hibernate", "off")
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
			cmd = exec.Command("psshutdown.exe", "-d", "-t", "0", "-accepteula")
			err = cmd.Run()
			if err != nil {
				panic(err)
			}
		}
	default:
		println("add your operating system")
	}
}
