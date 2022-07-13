package restfulapi

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type UserInfo struct {
	UserId string
	PubKey []byte
}

func SendKey(uid string, pubkey []byte) {
	info := &UserInfo{uid, pubkey}
	body, _ := json.Marshal(info)
	_, err := http.Post("http://localhost:80/save_key", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Print(err)
	}
}
