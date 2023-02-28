package clinet

import (
	"Turn_on_PC/internal/DTO"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetJWT(host string, login string, password string, scope string) (string, error) {
	url := url.URL{Scheme: "HTTP", Host: host, Path: "oauth"}
	user := DTO.UserSingIn{Login: login, Password: password, Scope: scope}
	userJson, err := json.Marshal(&user)
	if err != nil {
		return "", err
	}
	r, err := http.Post(url.String(), "application/json", bytes.NewBuffer(userJson))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	MyJWT := string(b)
	return fmt.Sprintf("Bearer %s", MyJWT), nil
}
