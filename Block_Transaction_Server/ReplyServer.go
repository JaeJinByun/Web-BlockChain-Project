package main

import (
	"encoding/json"
	f "fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var flag = 0

//PBFT 응답 구조체
type Result struct {
	Hash   []byte `json:"Hash"`
	Height int64  `json:"Height"`
}

func main() {
	startTime := time.Now()

	//PBFT 합의 완료 답장
	http.HandleFunc("/reply", func(res http.ResponseWriter, req *http.Request) {
		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			data := &Result{}
			err := json.Unmarshal([]byte(respBody), data)

			if err != nil {
				f.Println("에러  : 합의 답장 실패")

			} else {

				f.Println("---------------------- 합의 완료 / 받은 블록정보 ----------------------\n")
				f.Printf("Hash : %x\n", data.Hash)
				f.Printf("Height : %d\n", data.Height)
				f.Println("-----------------------------------------------------------------------\n")
				if data.Height == 10000 {
					delta := time.Now().Sub(startTime)
					f.Printf("done in %.3fs.\n", delta.Seconds())
				}
				//cntHeight = data.Height
				flag = 1 //전송 실행 상태 On
			}
		}

		defer req.Body.Close()

	})

	http.ListenAndServe(":83", nil) //83번 포트에서 웹 서버 실행

}
