package Project

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	f "fmt"
	"math/rand"
	"strconv"
	"time"
)

type Block struct {
	Hash          []byte `json:"Hash"`          //Current Hash
	PrevBlockHash []byte `json:"PrevBlockHash"` //Previous Hash
	Timestamp     int64  `json:"Timestamp"`     //Block was create To Time
	Nonce         int    `json:"Nonce"`         //Random Num that have Ordering
	TxID          []byte `json:"TxID"`          //sha256(Data + TimeStamp + Nonce)
	Bits          int    `json:"Bits"`          //Targetbits
	Pow           []byte `json:"Pow"`           //Hash from Pow
	Height        int    `json:"Height"`        //Block Height
}

func (bc *Block) BPrint() {
	f.Println("-------------------------------보낸 블록 데이터 ---------------------------------")
	f.Printf("Prev Hash    	 : %x\n", bc.PrevBlockHash)
	f.Println("TimeStamp 	 :", time.Unix(bc.Timestamp, 0))
	f.Printf("Nonce 		 : %d\n", bc.Nonce)
	f.Printf("bits		 : %d\n", bc.Bits)
	f.Printf("hash 		 : %x\n", bc.Hash)
	f.Printf("Height 	 	 : %d\n", bc.Height)
	f.Printf("TxID		:  %x\n", bc.TxID)
	f.Println("-------------------------------------------------------------------------------------\n")
}

func (block *Block) setHash(TxID []byte) {

	timestamp := strconv.FormatInt(block.Timestamp, 10)
	timeBytes := []byte(timestamp)

	//현재 블록의 해시값 = timeBytes + PrevBlockHash + Data 를 합친 값
	var blockBytes []byte
	blockBytes = append(timeBytes, block.PrevBlockHash...)
	blockBytes = append(blockBytes, block.Pow...)
	blockBytes = append(blockBytes, TxID...)
	// 		↳--------------↴
	hash := sha256.Sum256(blockBytes)
	block.Hash = hash[:]

}

var Height = 0

func NewBlock(TxID []byte, prevBlockHash []byte) *Block {
	block := &Block{}

	//만들어진 시각 , 이전 블록의 해시값 , 블록의 데이터
	block.Timestamp = time.Now().Unix()
	block.PrevBlockHash = prevBlockHash

	//Block 정보 세팅
	block.setHash(TxID)

	block.Bits = rand.Intn(10)
	pow := newProofOfWork(block)
	nonce, hash := pow.Run()
	block.Pow = hash[:]
	block.Nonce = nonce
	block.TxID = TxID
	block.Height = Height
	Height++
	return block
}

func NewGenesisBlock() *Block {
	Tx := &Data{"Genesis Block", 0, "Genesis Block", 0, []byte("Genesis Block"), []byte("Genesis Block")}
	bytes, _ := json.Marshal(Tx)
	return NewBlock(bytes, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

//해시값의 체이닝 여부
func (b *Block) EqualHash(e []byte) bool {
	return bytes.Equal(b.Hash, e)
}

//Data의 정합성 여부
func (b *Block) EqualData(e []byte) bool {
	return bytes.Equal(b.TxID, e)
}

//블록찾기
func (bc *Blockchain) FindBlock(id []byte) *Block {
	for _, v := range bc.Blocks {
		if v.EqualHash(id) {
			return v
		}
	}
	return nil
}

func (b *Block) IsGenBlock() bool {
	return bytes.Equal(b.PrevBlockHash, make([]byte, len(b.PrevBlockHash)))

}
