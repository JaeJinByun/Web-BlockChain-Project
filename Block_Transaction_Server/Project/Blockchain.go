package Project

type Blockchain struct {
	Blocks []*Block `json:"Blocks`
}

func (blockchain *Blockchain) AddBlock(TxID []byte) {

	//이전 블록 take
	prevBlock := blockchain.Blocks[len(blockchain.Blocks)-1]
	//새로운 블록생성 넘길값 data , prevBlockHash
	newBlock := NewBlock(TxID, prevBlock.Hash)
	blockchain.Blocks = append(blockchain.Blocks, newBlock)

}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}
