package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// check if last hash is found if none found make the new blockchain
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()

			// creating the genesis block
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
			// to be handles by the function outside

		} else {

			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})
	Handle(err)

	blockchain := Blockchain{
		LastHash: lastHash,
		Database: db,
	}
	return &blockchain

}

func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(
		// view the database in view only mode
		func(txn *badger.Txn) error {

			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err
			// copy the last Hash from the datbase into the lastHash var defined
		},
	)
	Handle(err)

	newBlock := CreateBlock(data, lastHash)
	// making the new block with the data input as well as the lastHash that we just got.

	err = chain.Database.Update(

		func(txn *badger.Txn) error {
			err := txn.Set(newBlock.Hash, newBlock.Serialize())
			Handle(err)
			// sets the stuff inside the new block
			err = txn.Set([]byte("lh"), newBlock.Hash)
			// adds the current hash as the "lh" for the new block

			chain.LastHash = newBlock.Hash

			return err

		})
	Handle(err)

}
func (chain *Blockchain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
	// custom iterator for the chain

}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get(iter.CurrentHash)
			// goes in a reverse order and gets the deserialized block.
			Handle(err)
			encodedBlock, err := item.ValueCopy(nil)
			block = Deserialize(encodedBlock)

			return err
		},
	)
	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}
