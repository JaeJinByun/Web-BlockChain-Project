package Project

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	f "fmt"
	"math/big"
	"strconv"
	"time"
)

type Tx struct {
	TxID      []byte `json:"TxID"`      //sha256(Data + TimeStamp + Nonce)
	TimeStamp int64  `json:"TimeStamp"` //Tx 생성시간
	UserId    string `json:"UserId"`    //Tx 발생 시킨 유저 ID
	LogDb     int    `json:"LogDb"`     //DB 의 정보
	Content   string `json:"Content"`   //
	RId       int    `json:"RId"`       //레스토랑 아이디
}

type Txs struct {
	Txs []*Tx `json: Txs`
}

func New_Transaction_Struct() *Txs {
	return &Txs{}
}

func (txs *Txs) AddTx(data []byte) {
	newTx := NewTranscation(data)
	txs.Txs = append(txs.Txs, newTx)
}

//트랜잭션 생성
func NewTranscation(data []byte) *Tx {
	Tx := &Tx{}

	//언마샬하여  Tx 구조체형식으로 저장
	err := json.Unmarshal(data, Tx)
	if err != nil {
		panic(err)
	}

	Tx.TimeStamp = time.Now().Unix()

	timestamp := strconv.FormatInt(Tx.TimeStamp, 10)
	timeBytes := []byte(timestamp)

	//var blockBytes []byte
	blockBytes := append(timeBytes, Tx.UserId...)
	blockBytes = append(blockBytes, []byte(string(Tx.LogDb))...)
	blockBytes = append(blockBytes, []byte(Tx.Content)...)
	blockBytes = append(blockBytes, []byte(string(Tx.RId))...)
	// 		↳--------------↴
	hash := sha256.Sum256(blockBytes)
	Tx.TxID = hash[:]
	return Tx
}

//넘어온 데이터 검증
func Verify(pubKey []byte, Sign []byte, hashid []byte) bool {
	curve := elliptic.P256()

	//서명 데이터 분할
	r := big.Int{}
	s := big.Int{}
	siglen := len(Sign)
	r.SetBytes(Sign[:(siglen / 2)])
	s.SetBytes(Sign[(siglen / 2):])

	//공개키 분할
	x := big.Int{}
	y := big.Int{}
	keylen := len(pubKey)
	x.SetBytes(pubKey[:(keylen / 2)])
	y.SetBytes(pubKey[(keylen / 2):])

	//공개키 찾기
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}

	//찾은 공개키로 서명 검증
	return ecdsa.Verify(&rawPubKey, hashid, &r, &s)
}

func (tx *Tx) TPrint(i int) {
	f.Println("-------------------------------", i, "번째 트랜잭션 ---------------------------------")
	// f.Printf("TxID    	 : %x\n", tx.TxID)
	// f.Println("TimeStamp 	 :", time.Unix(tx.TimeStamp, 0))
	// f.Printf("UserID 		 : %s\n", tx.UserID)
	// f.Printf("LogDB		 : %s\n", tx.LogDB)
	// f.Printf("Content 	 : %s\n", tx.Content)
	// f.Printf("RId 	 	 : %s\n", tx.RId)
	// f.Println("-------------------------------------------------------------------------------------\n")
}

// 해당 Id를 가지고있는 트랜잭션들
func (txs *Txs) Find_tx(UserID string) *Txs {
	UserTxs := &Txs{}
	for _, v := range txs.Txs {
		if v.UserId == UserID {
			UserTxs.Txs = append(UserTxs.Txs, v)
		}
	}
	return UserTxs
}

//해당 레스토랑이름 , db의 정보를 가진 트랜잭션중 제일 최근의 트랜잭션
func (txs *Txs) Find_Current_tx(Rid int, LogDb int) *Tx {
	CurrentTxs := &Txs{}
	for _, v := range txs.Txs {
		if v.RId == Rid && v.LogDb == LogDb {

			CurrentTxs.Txs = append(CurrentTxs.Txs, v)
		}
	}

	if CurrentTxs.Txs != nil {
		return CurrentTxs.Txs[len(CurrentTxs.Txs)-1]
	} else {
		return nil
	}

}

//레스토랑 id로 하여금 특정 조건에 따른 최신 트랜잭션들
func (txs *Txs) GetRTxs(RidTxs *RidTxs) *Txs {
	resultTxs := &Txs{}

	for _, v := range RidTxs.Rts {
		respTxs := &Txs{}
		for _, v2 := range txs.Txs {
			if v.LogDb == v2.LogDb && v.RId == v2.RId {
				respTxs.Txs = append(respTxs.Txs, v2)
			}
		}
		if len(respTxs.Txs) > 0 {
			resultTxs.Txs = append(resultTxs.Txs, respTxs.Txs[len(respTxs.Txs)-1])
		}
	}

	return resultTxs
}
