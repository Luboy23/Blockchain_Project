package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T){
	b :=randomBlock(0, types.Hash{})
	fmt.Println(b.Hash(BlockHasher{}))
}

func TestSignBlock(t *testing.T) {
	priKey := crypto.GeneratePrivatekey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(priKey))
	assert.NotNil(t,b.Signature)
}

func TestBlockVerify(t *testing.T) {
	priKey := crypto.GeneratePrivatekey()
	b := randomBlock(0, types.Hash{})

	 assert.Nil(t, b.Sign(priKey))
	assert.Nil(t,b.Verify())

	priKey1 := crypto.GeneratePrivatekey()
	b.Validator = priKey1.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())

}
 

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version: 1,
		PrevBlockHash: prevBlockHash,
		Height:	height,
		Timestamp:  time.Now().UnixNano(),
	}

	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T,height uint32, prevBlockHash types.Hash) *Block {
	priKey := crypto.GeneratePrivatekey()
	b := randomBlock(height,prevBlockHash)

	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	b.Sign(priKey)

	assert.Nil(t, b.Sign(priKey))
	return b
}