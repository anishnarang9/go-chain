package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockChain(address string) *Blockchain {
	// Creates new blocks for alr existing chain

	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}
	var lastHash []byte

	// configure and opens the existing DB
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	// reads the saved latest block hash
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	// cretes a new chain struct in memory with updated lh and db connection
	chain := Blockchain{lastHash, db}
	return &chain

}

func InitBlockChain(address string) *Blockchain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("DB already Exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// create the coinbase tx with the address to receive the rewards
		cbtx := CoinbaseTx(address, genesisData)
		genesis := Genesis(cbtx)

		// creating the genesis block
		fmt.Println("Genesis proved")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
		// to be handles by the function outside

	})

	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain

}

func (chain *Blockchain) AddBlock(txs []*Transaction) {
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

	newBlock := CreateBlock(txs, lastHash)
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
func (chain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	// final result

	spentTXOs := make(map[string][]int)
	// temporary lookup map to disregard the TXOs that have been spent

	iter := chain.Iterator()

	for {
		block := iter.Next()
		// one block at a time

		for _, tx := range block.Transactions {
			// for each transaction within that block
			txID := hex.EncodeToString(tx.ID)
			// gets the string encoding for all the byte txIDs

		Outputs:
			for outIdx, out := range tx.Outputs {
				// for each output within that transaction
				if spentTXOs[txID] != nil {
					// check if output is inside the map
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
							// if its already in the map we continue with the outputs for loop
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
					// all the outputs that can be unlocked are now UnspentTxs
				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
						// for all other transactions that are not coinbase we check to see
						// if they are referenced as inputs they are added to SpentTXOs
						// disregarded above
					}
				}
			}

		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs

}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {

	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs

}

func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	// makes a list of the outputs that were using/accumulating
	unspentTxs := chain.FindUnspentTransactions(address)
	// uses the function we defined earlier to  find all unspent transactions
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		// for each unspent transaction

		for outIdx, out := range tx.Outputs {
			// if the outputs within that transaction can be unlocked
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				// adds the accumulated value and adds it to our map of transactions were using
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
					// once we have enough it breaks
				}
			}

		}
	}
	return accumulated, unspentOuts
}
