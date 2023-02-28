package main

import (
	"Turn_on_PC/internal/clinet"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 5 {
		panic("You must pass url http/https, login, password, name")
	}
	JWT, err := clinet.GetJWT(args[1], args[2], args[3], "ClientPC")
	if err != nil {
		panic(err)
	}
	clinet.StartWSPC(args[1], JWT, args[4])
}
