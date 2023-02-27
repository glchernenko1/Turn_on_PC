package main

import (
	"Turn_on_PC/internal/clinet"
)

func main() {
	JWT, err := clinet.GetJWT("127.0.0.1:1234", "122223", "221232131", "ServerUser")
	if err != nil {
		panic(err)
	}
	clinet.StartWSServer("127.0.0.1:1234", JWT, "kek", "/home/google/GolandProjects/Turn_on_PC/cmd/clientServerPC/user.json")

}
