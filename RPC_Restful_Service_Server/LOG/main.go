package main // rest

import (
	"blockwiki/LOG/rpc_client"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var logs = map[string]*Log{}

type Log struct {
	UserId   string `json:"UserId"`
	LogDb    int    `json:"LogDb"` // 제목 0, 사진 1, 설명 2
	RId      int    `json:"RId"`
	Content  string `json:"Content"`
	PrivateK string `json:"PrivateK"`
	// Sign    []byte `json:"Sign"`   //Signature <= UserID 를 SHA256() 해시화 하여 개인키로 (ecdsa.Sign() 함수 이용 ) 암호화한 값
	// HashId  []byte `json:"HashId"` //Id를 sha256 돌린값
}

type Current struct {
	LogDb int `json:"LogDb"` // 제목 0, 사진 1, 설명 2
	RId   int `json:"RId"`
}

type Currents struct { // 레스토랑 여러개 담는 구조체
	Rts []*Current `json:"Rts"`
}

type Tx struct {
	TxID      []byte `json:"TxID"`      //sha256(Data + TimeStamp + Nonce)
	TimeStamp int64  `json:"TimeStamp"` //Tx 생성시간
	UserId    string `json:"UserId"`    //Tx 발생 시킨 유저 ID
	LogDb     int    `json:"LogDb"`     //LogDB 의 정보
	Content   string `json:"Content"`   //Tx 내용
	RId       int    `json:"RId"`       //Content Type
}

type Txs struct {
	Txs []*Tx `json:"Txs"`
}

// func InitRtList() []Current {
// 	rtList := []Current{}
// 	return rtList
// }
// func appendList(rtList *[]Current) {
// 	rt := Current{0, 1} // DB에서 조회해서 넘기기
// 	*rtList = append(*rtList, rt)
// }

// 4
func jsonContentTypeMiddleware(next http.Handler) http.Handler {

	// 들어오는 요청의 Response Header에 Content-Type을 추가
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(rw, r)

	})
}

func main() {

	mux := http.NewServeMux()

	createHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodOptions:
			wr.WriteHeader(204)
			return
		case http.MethodPost: // 등록 - 하는 함수를 호출(서버에 POST하기)
			// parameter 받아서 출력해보기
			r.ParseForm()
			fmt.Println(r.Form)

			log := &Log{}

			log.UserId = r.Form.Get("id")
			log.LogDb, _ = strconv.Atoi(r.Form.Get("log")) // 항목별로 구분해줄 것
			log.RId, _ = strconv.Atoi(r.Form.Get("l_id"))
			log.Content = r.Form.Get("edit")
			log.PrivateK = r.Form.Get("privk")

			fmt.Println(log.Content)

			// 받은 걸 넣어주기. 여기서 포스팅이 들어가면 될것 같음
			pbytes, _ := json.Marshal(log)
			buff := bytes.NewBuffer(pbytes)

			resp, err := http.Post("http://localhost:8081/logs", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}
	})

	currentHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			wr.WriteHeader(204)
			return
		case http.MethodPost: // 등록 - 하는 함수를 호출(서버에 POST하기)
			// parameter 받아서 출력해보기
			r.ParseForm()
			fmt.Println(r.Form)
			var current Current

			current.LogDb, _ = strconv.Atoi(r.Form.Get("log_db"))
			current.RId, _ = strconv.Atoi(r.Form.Get("r_id"))

			json.NewDecoder(r.Body).Decode(&current)

			//logs[log.UserId] = &log

			// json.NewEncoder(wr).Encode(current)

			//fmt.Println(log)

			// 받은 걸 넣어주기. 여기서 포스팅이 들어가면 될것 같음
			pbytes, _ := json.Marshal(current)
			buff := bytes.NewBuffer(pbytes)

			// resp, err := http.Post("http://localhost:8081/current", "application/json", buff)
			resp, err := http.Post("http://localhost:81/current_tx", "application/json", buff)
			if err != nil {
				panic(err)
			}
			// fmt.Println(resp) // 출력확인용 이걸 write로 써주기

			tx := &Tx{}
			response, err := ioutil.ReadAll(resp.Body)

			json.Unmarshal(response, tx)
			fmt.Println(tx.Content)
			fmt.Println(tx.LogDb)
			fmt.Println(tx.RId)
			fmt.Println(tx.TimeStamp)
			fmt.Println(tx.TxID)
			fmt.Println(tx.UserId)

			wr.Write(response)
			// for _, v := range tx.Txs {
			// 	fmt.Println(v.Content)
			// }
			defer resp.Body.Close()

		}
	})

	currentsHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			wr.WriteHeader(204)
			return
		case http.MethodPost: // 등록 - 하는 함수를 호출(서버에 POST하기)
			// parameter 받아서 출력해보기

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			// rt := []*Current{}
			rts := &Currents{}

			// // 확인
			json.Unmarshal(body, rts)
			for _, v := range rts.Rts {
				fmt.Printf("LogDb %d :", v.LogDb)
				fmt.Printf("RId %d :", v.RId)
			}
			fmt.Println()

			// 받은 걸 넣어주기. 여기서 포스팅이 들어가면 될것 같음
			// pbytes, _ := json.Marshal(rts)
			buff := bytes.NewBuffer(body)

			// resp, err := http.Post("http://localhost:8081/current", "application/json", buff)
			resp, err := http.Post("http://localhost:81/rid_txs", "application/json", buff)
			if err != nil {
				panic(err)
			}

			// body, err := ioutil.ReadAll(resp.Body)
			txs := &Txs{}
			response, err := ioutil.ReadAll(resp.Body)
			json.Unmarshal(response, txs)

			wr.Write(response)
			for _, v := range txs.Txs {
				fmt.Println(v.Content)
			}

			defer resp.Body.Close()

		}
	})

	logsHandler := http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodOptions:
			wr.WriteHeader(204)
			return
		case http.MethodPost: // 등록 - 하는 함수를 호출(서버에 POST하기)
			// parameter 받아서 출력해보기

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			// rt := []*Current{}
			rts := &Currents{}

			// // 확인
			json.Unmarshal(body, rts)
			for _, v := range rts.Rts {
				fmt.Printf("LogDb %d :", v.LogDb)
				fmt.Printf("RId %d :", v.RId)
			}
			fmt.Println()

			// 받은 걸 넣어주기. 여기서 포스팅이 들어가면 될것 같음
			// pbytes, _ := json.Marshal(rts)
			buff := bytes.NewBuffer(body)

			// resp, err := http.Post("http://localhost:8081/current", "application/json", buff)
			resp, err := http.Post("http://localhost:81/rid_txs", "application/json", buff)
			if err != nil {
				panic(err)
			}

			// body, err := ioutil.ReadAll(resp.Body)
			txs := &Txs{}
			response, err := ioutil.ReadAll(resp.Body)
			json.Unmarshal(response, txs)

			wr.Write(response)
			for _, v := range txs.Txs {
				fmt.Println(v.Content)
			}

			defer resp.Body.Close()

		}
	})

	// 3
	mux.Handle("/create", jsonContentTypeMiddleware(createHandler))     // 트랜잭션 생성
	mux.Handle("/current", jsonContentTypeMiddleware(currentHandler))   // 트랜잭션 조회
	mux.Handle("/currents", jsonContentTypeMiddleware(currentsHandler)) // 식당 리스트 조회
	mux.Handle("/logs", jsonContentTypeMiddleware(logsHandler))         // 생성 내역 조회

	mux.Handle("/create_wallet", rpc_client.UserInfoReq{}) // 지갑 생성
	http.ListenAndServe(":7779", mux)
}
