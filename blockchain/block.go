package blockchain

type Blockchain struct {
	Blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// Archived Derive Hash since POW implemented
/*
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.prevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]

}
*/

func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), PrevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	//block.DeriveHash()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}
func Genesis() *Block {
	return CreateBlock("genesis", []byte{})
}

func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{(Genesis())}}
}
