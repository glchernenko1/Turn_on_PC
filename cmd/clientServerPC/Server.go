package main

import (
	"Turn_on_PC/internal/clinet"
	"os"
)

func main() {

	args := os.Args
	if len(args) != 6 {
		panic("You must pass url , login, password, name, path_to_json")
	}

	JWT, err := clinet.GetJWT(args[1], args[2], args[3], "ServerUser")
	if err != nil {
		panic(err)
	}
	clinet.StartWSServer(args[1], JWT, args[4], args[5])

}
