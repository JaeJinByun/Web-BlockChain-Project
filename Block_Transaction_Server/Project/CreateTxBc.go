package Project

import (
	"encoding/json"
	f "fmt"
	"io/ioutil"
	"os"
)

var cnt = 1

func (txs *Txs) CreateTx(data []byte) {
	//트랜잭션 생성
	txs.AddTx(data)

	tx := txs.Txs[len(txs.Txs)-1]
	tx.TPrint(cnt)
	cnt++

	Tx_doc, _ := json.MarshalIndent(txs.Txs, "", " ") //BlockChain Blocks
	err := ioutil.WriteFile("TxFile.json", Tx_doc, os.FileMode(0644))

	if err != nil {
		f.Println(err)
		return
	}
}

func (bc *Blockchain) CreateBc(TxID []byte) {

	//블록 생성 , 블록 추가
	bc.AddBlock(TxID)

	//Copy Origin Block struct then Paste Format of Json
	block_doc, _ := json.MarshalIndent(bc.Blocks, "", " ") //BlockChain Blocks
	err := ioutil.WriteFile("BlockFile.json", block_doc, os.FileMode(0644))

	if err != nil {
		f.Println(err)
		return
	}

}
