package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
	"github.com/stretchr/testify/assert"
)

//	测试区块的哈希功能。
func TestHashBlock(t *testing.T) {
	b := randomBlock(0, types.Hash{}) // 创建一个随机区块
	fmt.Println(b.Hash(BlockHasher{})) // 打印区块的哈希值
}

//	测试区块的签名功能。
func TestSignBlock(t *testing.T) {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	b := randomBlock(0, types.Hash{}) // 创建一个随机区块

	assert.Nil(t, b.Sign(priKey)) // 断言签名操作不会返回错误
	assert.NotNil(t, b.Signature) // 断言区块有签名
}

//	测试区块的验证功能。
func TestBlockVerify(t *testing.T) {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	b := randomBlock(0, types.Hash{}) // 创建一个随机区块

	assert.Nil(t, b.Sign(priKey)) // 断言签名操作不返回错误
	assert.Nil(t, b.Verify()) // 断言验证操作不返回错误

	priKey1 := crypto.GeneratePrivatekey() // 生成另一个私钥
	b.Validator = priKey1.PublicKey() // 更改验证者公钥
	assert.NotNil(t, b.Verify()) // 断言验证操作返回错误，因为验证者公钥不匹配

	b.Height = 100 // 更改区块高度
	assert.NotNil(t, b.Verify()) // 断言验证操作返回错误，因为区块高度不匹配
}

//	创建一个随机区块。
func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1, // 设置区块版本为1
		PrevBlockHash: prevBlockHash, // 设置前一个区块的哈希值
		Height:        height, // 设置区块高度
		Timestamp:     time.Now().UnixNano(), // 设置当前时间戳
	}

	return NewBlock(header, []Transaction{}) // 创建一个新的区块，包含指定的区块头和空的交易列表
}

//	创建一个随机区块，并为其签名。
func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	b := randomBlock(height, prevBlockHash) // 创建一个随机区块

	tx := randomTxWithSignature(t) // 创建一个随机交易，并为其签名
	b.AddTransaction(tx) // 将交易添加到区块中
	b.Sign(priKey) // 为区块签名

	assert.Nil(t, b.Sign(priKey)) // 断言签名操作不返回错误
	return b
}
