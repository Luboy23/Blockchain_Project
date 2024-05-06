package core

import (
	"fmt"
	"testing"

	"github.com/Luboy23/Blockchain_Project/types"
	"github.com/stretchr/testify/assert"
)

//	创建一个包含创世区块的新区块链。
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0, types.Hash{})) // 创建新的区块链，包含创世区块
	assert.Nil(t, err) // 断言没有错误

	return bc // 返回新创建的区块链
}

//	 测试向区块链中添加新的区块。
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t) // 创建一个新的区块链

	lenBlocks := 1000 // 定义要添加的区块数量
	for i := 0; i < lenBlocks; i++ { // 循环添加区块
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1))) // 创建并签名新的区块
		assert.Nil(t, bc.AddBlock(block)) // 断言添加区块不返回错误
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks)) // 断言区块链高度正确
	assert.Equal(t, len(bc.headers), lenBlocks+1) // 断言区块头列表长度正确

	assert.NotNil(t, bc.AddBlock(randomBlock(48, types.Hash{}))) // 断言添加无效区块返回错误
}

//	测试创建新的区块链。
func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t) // 创建一个新的区块链

	assert.NotNil(t, bc.validator) // 断言验证器不为空

	assert.Equal(t, bc.Height(), uint32(0)) // 断言区块链高度为0

	fmt.Println(bc.Height()) // 打印区块链高度
}

//	测试检查区块是否存在于区块链中。
func TestTheHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t) // 创建一个新的区块链

	assert.True(t, bc.HashBlock(0)) // 断言第0个区块存在
	assert.False(t, bc.HashBlock(1)) // 断言第1个区块不存在
	assert.False(t, bc.HashBlock(100)) // 断言第100个区块不存在
}

//	测试向特定高度添加区块。
func TestAddBlockToHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t) // 创建一个新的区块链

	assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, 1, getPrevBlockHash(t, bc, uint32(1))))) // 断言添加第1个区块不返回错误

	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 3, types.Hash{}))) // 断言添加第3个区块返回错误
}

//	测试获取指定高度的区块头。
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t) // 创建一个新的区块链

	for i := 0; i < 1000; i++ { // 循环添加区块
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1))) // 创建并签名新的区块
		assert.Nil(t, bc.AddBlock(block)) // 断言添加区块不返回错误

		header, err := bc.GetHeader(uint32(i + 1)) // 获取指定高度的区块头
		assert.Nil(t, err) // 断言获取区块头不返回错误
		assert.Equal(t, header, block.Header) // 断言获取的区块头与预期相同
	}
}

//	获取指定高度的前一个区块的哈希值。
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1) // 获取前一个区块的头
	assert.Nil(t, err) // 断言获取区块头不返回错误

	return BlockHasher{}.Hash(prevHeader) // 返回前一个区块的哈希值
}
