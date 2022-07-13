package main // service

import (
	"blockwiki/custom/service"
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var logs = map[string]*Log{}

type Log struct {
	UserId   string `json:"UserId"`
	LogDb    int    `json:"LogDb"`
	RId      int    `json:"RId"`
	Content  string `json:"Content"`
	PrivateK string `json:"PrivateK"`
}

type Tx struct {
	UserId  string `json:"UserId"`
	LogDb   int    `json:"LogDb"`
	RId     int    `json:"RId"`
	Content string `json:"Content"`
	Sign    []byte `json:"Sign"`   //Signature <= UserID 를 SHA256() 해시화 하여 개인키로 (ecdsa.Sign() 함수 이용 ) 암호화한 값
	HashId  []byte `json:"HashId"` //Id를 sha256 돌린값
}

// 4
func jsonContentTypeMiddleware(next http.Handler) http.Handler {

	// 들어오는 요청의 Response Header에 Content-Type을 추가
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		// 전달 받은 http.Handler를 호출한다.
		next.ServeHTTP(rw, r)
		// fmt.Printf("%v", rw)
		// http.Post("localhost:8080/logs2", "application/json", r.Body)
	})
}

func main() {

	// 텍스트로 포스팅 해보기
	// log := "{\"UserId\": \"aaa\",\"LogDb\": \"log_r_name\",\"RId\": 5,\"Content\": \"버거킹\"}"

	// reqBody := bytes.NewBufferString(log)

	// resp, err := http.Post("http://127.0.0.1:8080/logs", "text/plain", reqBody)

	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	mux := http.NewServeMux()

	userHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		switch r.Method {
		// case http.MethodGet: // 조회 - 하는 함수를 호출(서버에 요청?)
		// 	var log Log
		// 	err := json.Unmarshal(b, &log)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Println(log)

		case http.MethodPost:
			// 확인용
			var log Log
			err := json.Unmarshal(b, &log)
			if err != nil {
				panic(err)
			}
			fmt.Println(log)
			fmt.Println(len([]byte(log.PrivateK)))

			// var r *ecdsa.PrivateKey
			d, _ := hex.DecodeString(log.PrivateK)
			r, _ := x509.ParseECPrivateKey(d)

			// fmt.Println("------------------")
			// fmt.Println(r)
			tx := &Tx{}

			tx.UserId = log.UserId
			tx.LogDb = log.LogDb
			tx.RId = log.RId
			tx.Content = log.Content
			tx.HashId = service.Hash(log.UserId)
			tx.Sign = service.Sign(*r, tx.HashId)

			b, _ := json.Marshal(tx)

			buff := bytes.NewBuffer(b)
			resp, err := http.Post("http://localhost/create_bc", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}
	})

	currentHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		switch r.Method {
		case http.MethodGet: // 조회 - 하는 함수를 호출(서버에 요청?)
			var log Log
			err := json.Unmarshal(b, &log)
			if err != nil {
				panic(err)
			}
			fmt.Println(log)

		case http.MethodPost:
			// 확인용
			var log Log
			err := json.Unmarshal(b, &log)
			if err != nil {
				panic(err)
			}
			fmt.Println(log)

			// buff := bytes.NewBuffer(b)
			//resp, err := http.Post("http://192.168.10.24/current_tx", "application/json", buff)
			// if err != nil {
			// 	panic(err)
			// }
			// defer resp.Body.Close()
		}
	})

	// 3
	mux.Handle("/logs", jsonContentTypeMiddleware(userHandler))
	mux.Handle("/current", jsonContentTypeMiddleware(currentHandler))
	// mux.HandleFunc("/logs2", func(wr http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("111111")
	// 	fmt.Printf("%v", wr)
	// })
	http.ListenAndServe(":8081", mux)
}
