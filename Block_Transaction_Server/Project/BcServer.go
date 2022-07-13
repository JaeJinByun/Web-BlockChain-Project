package Project

import (
	b "bytes"
	"encoding/json"
	f "fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"net/http"
)

//PBFT 응답 구조체
type Result struct {
	Hash   []byte `json:"Hash"`
	Height int64  `json:"Height"`
}

//UserId - PubKey 받아오기 위한 구조체
type UserKey struct {
	UserId string `json:"UserId"`
	PubKey []byte `json:"PubKey"`
}

var cntHeight int64
var flag = 1

func StartBCServer() {
	//서버 키면 지정한 경로에 Level DB 생성
	db, err := leveldb.OpenFile("/Users/byunjaejin/Go/level_DB", nil)
	if err != nil {
		panic(err)
	}

	//서버 키면 제네시스 블럭을 PBFT 로 전송
	bc := NewBlockchain()

	Block := bc.Blocks[len(bc.Blocks)-1]
	bytes, _ := json.Marshal(Block)
	buff := b.NewBuffer(bytes)
	Block.BPrint()

	resp, err := http.Post("http://192.168.10.57:4000/pbft", "application/json", buff)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// localhsot:80/create_bc 에 접속시
	// 넘어오는 값은 트랜잭션 내용
	http.HandleFunc("/create_bc", func(res http.ResponseWriter, req *http.Request) {

		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			f.Println("트랜잭션 요청이 왔습니다. 검증을 시작합니다...")
			TxData := []byte(respBody)

			DataForSign := &Data{}

			err = json.Unmarshal(TxData, DataForSign)
			if err != nil {
				println("Json Unmarshal Fail")
			}
			f.Println("\n---------------------------------------------------------------------")
			f.Printf("UserId  : %s\n", DataForSign.UserId)
			f.Printf("LogDb   : %d\n", DataForSign.LogDb)
			f.Printf("Content : %s\n", DataForSign.Content)
			f.Printf("Rid     : %d\n", DataForSign.RId)
			f.Printf("Sign    : %x\n", DataForSign.Sign)
			f.Printf("HashId  : %x\n", DataForSign.HashId)
			f.Println("---------------------------------------------------------------------\n")

			/*
				1. DataForSign 의 UserId 와 Sign 을 뽑아온다.
				2. LevelDB에서 UserId 가 가진 공개키를 가져온다.
				3. 가져온 공개키로 Sign 을 Verify() 한다.
				4. bool 값에 따라 처리한다.
			*/

			//1. UserId, Sign 값 가져오기
			UserID := DataForSign.UserId
			Sign := DataForSign.Sign
			HashId := DataForSign.HashId
			//2.levelDB에서 ID에 맞는공개키 가져오기
			data, err := db.Get([]byte(UserID), nil)
			if err != nil {
				println("공개키 못 찾았음")
			}
			//3. Verify()
			if !Verify(data, Sign, HashId) {
				f.Println("<검증 성공> 트랜잭션을 생성합니다..")
				bytes, _ := json.Marshal(DataForSign)
				buff := b.NewBuffer([]byte(bytes))
				resp, err := http.Post("http://localhost:81/create_tx", "application/json", buff)

				if err != nil {
					panic(err)
				}

				defer resp.Body.Close()
			}

		}
		defer req.Body.Close()

	})

	/* 트랜잭션 검증이 끝난후 Tx서버에서 넘긴 TxID로 블록 생성 요청  */
	http.HandleFunc("/newblock", func(res http.ResponseWriter, req *http.Request) {
		respBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			TxID := []byte(respBody)
			bc.CreateBc(TxID) //블록 생성
		}

		/*--------- PBFT 서버에 POST 전달 --------*/
		if flag == 1 {
			flag = 1 //전송대기
			Block := bc.Blocks[len(bc.Blocks)-1]
			bytes, _ := json.Marshal(Block)
			buff := b.NewBuffer(bytes)
			f.Println("PBFT 메인노드에 블록정보를 보냅니다.")
			resp, err := http.Post("http://192.168.10.57:4000/pbft", "application/json", buff)

			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()

			respBody, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				str := string(respBody)
				println(str)
			}
		}
		/*--------------------------------------*/
		cntHeight++
	})

	//회원의 ID에 맞는 공개키를 levelDB에 저장
	http.HandleFunc("/save_key", func(res http.ResponseWriter, req *http.Request) {
		respBody, err := ioutil.ReadAll(req.Body)
		UserKey := &UserKey{}

		if err == nil {
			err := json.Unmarshal(respBody, UserKey)
			if err != nil {
				println("유저키 언마샬 실패")
			}

			//Input DB			Key		        Value
			db.Put([]byte(UserKey.UserId), UserKey.PubKey, nil)
			//Key 로 Vaule Get
			DBdata, _ := db.Get([]byte(UserKey.UserId), nil)
			f.Println("저장된 유저 아이디와 공개키 정보")
			f.Printf("UserId : %s\n", UserKey.UserId)
			f.Printf("PubKey : %x\n", DBdata)

		}

	})

	http.ListenAndServe(":80", nil) //80번 포트에서 웹 서버 실행
}
