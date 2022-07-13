package rpcclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RequestRPC struct { // handler for localhost:7777/create
	UserId string // :7775 rpc서버로 UserId 전달
	Port   string
}
type userId struct { // for unmarshal request body
	UserId string
}

func (rpc RequestRPC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Send("super Wtreasure", rpc.Port)
	user := userId{}
	body, err := ioutil.ReadAll(r.Body) // []byte return
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &user)
	userinfo := Send(user.UserId, rpc.Port) // rpc_util.go파일. rpc서버에 등록된 함수호출. *User반환
	req := ToJson(userinfo)
	w.Write(req)
}
func Start() {
	mux := http.NewServeMux()
	mux.Handle("/create", RequestRPC{Port: ":7775"})
	log.Print("Listening port 7777...")
	http.ListenAndServe(":7777", mux)
}
