package core

import (
	"fmt"
	"testing"

	"github.com/Luboy23/Blockchain_Project/types"
	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain{
	bc, err := NewBlockchain(randomBlock(0,types.Hash{}))
	assert.Nil(t, err)

return bc
}
	
func TestAddBlock(t *testing.T){
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i <1000; i++ {
		block := randomBlockWithSignature(t,uint32(i + 1), getPrevBlockHash(t, bc, uint32(i + 1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks) )
	assert.Equal(t, len(bc.headers), lenBlocks + 1 )

	assert.NotNil(t, bc.AddBlock(randomBlock(48,types.Hash{})))
}

func TestNewBlockchain (t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.NotNil(t,bc.validator)

	assert.Equal(t, bc.Height(),uint32(0))

	fmt.Println(bc.Height())
}

func TestTheHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.True(t, bc.HashBlock(0))
	assert.False(t, bc.HashBlock(1))
	assert.False(t, bc.HashBlock(100))

}

func TestAddBlockToHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, 1, getPrevBlockHash(t, bc, uint32(1)))))

	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 3,types.Hash{})))
} 
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	// lenBlocks := 1000
	for i := 0; i <1000; i++ {
		block := randomBlockWithSignature(t,uint32(i + 1),getPrevBlockHash(t, bc, uint32(i +1)))
		assert.Nil(t, bc.AddBlock(block))

		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)

	}
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height -1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}
