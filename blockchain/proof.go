package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s.
// This is essentially the difficulty and ensures that the amount of bitcoin mined is limited and doesnt go crazy.

const Difficulty = 22

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
	// Creates a new Big int and sets a lower limit target using left shift
	// The lower the number the harder it is to find
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
		// updates the data to take in the nonce and difficulty
	)
	return data
}
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	// starting nonce as 0

	for nonce < math.MaxInt64 {
		// arbitrarily high number 9 quintillion
		// does this to run through every single possibility

		data := pow.InitData(nonce)
		// writes the new data using out InitData Algorithm
		hash = sha256.Sum256(data)
		//converts it Into a sha256 Hash

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])
		// converts the hash into Int

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
		// compares it to our target and keeps going till one lower than target is found

	}
	fmt.Println()
	return nonce, hash[:]
}
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)

	// use data, calculate hash convert it to int and compare them and if its less than target were good.

	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1

}

func ToHex(num int64) []byte {

	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
	// Helper function that converts integer to Binary in Big endian format
}
