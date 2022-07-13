package Project

import (
	b "bytes"
	"encoding/json"
	f "fmt"
	"io/ioutil"
	"net/http"
)

type Data struct {
	UserId  string `json:"UserID"`  //Tx 발생 시킨 유저 ID
	LogDb   int    `json:"LogDB"`   //LogDB 의 정보
	Content string `json:"Content"` //Tx 내용
	RId     int    `json:"RId"`     //Restaurant 번호
	Sign    []byte `json:"Sign"`    //Id를 Hash + privKey로 암호화한 값
	HashId  []byte `json:"HashId`   //Id를 Hash 한 값
}

type CurrentTxData struct {
	LogDb int `json:"LogDb"` //LogDB 의 정보
	RId   int `json:"RId"`   //Restaurant 번호
}

type RidTxs struct {
	Rts []*CurrentTxData `json:"Rts"`
}

func StartTxServer() {
	//서버 키면 트랜잭션 구조체 생성
	txs := New_Transaction_Struct()

	// localhsot:81/create_tx 에 접속시
	http.HandleFunc("/create_tx", func(res http.ResponseWriter, req *http.Request) {

		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			bytes := []byte(respBody)
			txs.CreateTx(bytes) //트랜잭션 생성
		}

		//생성된 트랜잭션의 TxID를 인자값으로 전달 -> 최상위 트랜잭션
		Tx := txs.Txs[len(txs.Txs)-1]

		//Block 생성 서버에 TxID Post 전송
		buff := b.NewBuffer(Tx.TxID)

		resp, err := http.Post("http://localhost:80/newblock", "text/plane", buff)

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		respBody, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			str := string(respBody)
			println(str)
		}

	})

	//http://localhost:81/find_tx 접속시 받는값은 UserID
	http.HandleFunc("/find_tx", func(res http.ResponseWriter, req *http.Request) {

		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			UserID := string(respBody)
			UserTxs := txs.Find_tx(UserID) //트랜잭션 조회
			//구조체 마샬
			//bytes, _ := json.Marshal(UserTxs)
			//buff := b.NewBuffer(bytes)

			for _, v := range UserTxs.Txs {
				f.Printf("TxID : %x\n", v.TxID)
				f.Printf("UserId : %s\n", v.UserId)
				f.Printf("Content : %s\n", v.Content)
			}
			//트랜잭션 조회 결과를 Restful Api로 응답
			bytes, _ := json.Marshal(UserTxs)
			res.Write(bytes)
		}
	})

	//http://localhost:81/current_tx 접속시 받는값은 Rid 의 0 1 2 3 값들의 최신 트랜잭션
	http.HandleFunc("/current_tx", func(res http.ResponseWriter, req *http.Request) {

		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			data := &CurrentTxData{}
			err := json.Unmarshal(respBody, data)
			if err != nil {
				f.Println("최신 트랜잭션 조회 언마샬 실패")
			}
			f.Println(data.RId)
			f.Println(data.LogDb)

			CurrentTx := txs.Find_Current_tx(data.RId, data.LogDb) //최신 트랜잭션 조회

			if CurrentTx != nil {
				f.Printf("UserId : %s\n", CurrentTx.UserId)
				f.Printf("Content : %s\n", CurrentTx.Content)
				f.Printf("Rid : %d\n", CurrentTx.RId)
			} else {
				f.Println("조회 결과 없음")
			}

			//구조체 마샬
			if CurrentTx != nil {
				bytes, _ := json.Marshal(CurrentTx)
				res.Write(bytes)
			} else {
				res.Write([]byte("null"))
			}

			//트랜잭션 조회 결과를 Restful Api로 응답

		}

		//res.Write([]byte(b)) //웹 브라우저에 응답
	})

	//http://192.168.10.24:81/Rid_txs 접속시 받는값은 Rid 의 0 1 2 3 값들의 최신 트랜잭션
	http.HandleFunc("/rid_txs", func(res http.ResponseWriter, req *http.Request) {

		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {

			data := &RidTxs{}

			err := json.Unmarshal(respBody, data)

			if err != nil {
				f.Println("트랜잭션 CTxs 언마샬 실패")
			}

			for _, v := range data.Rts {
				f.Printf("LogDb %d : ", v.LogDb)
				f.Printf("Rid %d : ", v.RId)

			}
			//탐색 시작
			respTXs := txs.GetRTxs(data)

			if respTXs.Txs == nil {
				res.Write([]byte("null"))
			} else {
				bytes, _ := json.Marshal(respTXs)
				res.Write(bytes)
			}

		}

	})

	http.ListenAndServe(":81", nil) //80번 포트에서 웹 서버 실행
}
