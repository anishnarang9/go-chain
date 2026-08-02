package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

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

func Genesis() *Block {
	return CreateBlock("genesis", []byte{})
}

func (b *Block) Serialize() []byte {

	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)

	// creates a new Go binary (gob) from res

	Handle(err)

	return res.Bytes()

}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	// converts the binary within the DB into the Block struct

	err := decoder.Decode(&block)

	Handle(err)

	return &block

}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
