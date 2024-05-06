package core

import (
	"bytes"
	"testing"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/stretchr/testify/assert"
)

// 测试交易的签名功能。
func TestSignTransaction(t *testing.T) {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	tx := &Transaction{
		Data: []byte("foo"), // 创建一个交易，数据为"foo"
	}

	assert.Nil(t, tx.Sign(priKey)) // 断言签名操作不返回错误
	assert.NotNil(t, tx.Signature) // 断言交易签名不为空
}

// 测试交易的验证功能。
func TestVerifyTransaction(t *testing.T) {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	tx := &Transaction{
		Data: []byte("foo"), // 创建一个交易，数据为"foo"
	}

	assert.Nil(t, tx.Sign(priKey)) // 断言签名操作不返回错误
	assert.Nil(t, tx.Verify()) // 断言验证操作不返回错误

	priKey1 := crypto.GeneratePrivatekey() // 生成另一个私钥
	tx.From = priKey1.PublicKey() // 更改交易的发送者公钥

	assert.NotNil(t, tx.Verify()) // 断言验证操作返回错误，因为签名不匹配
}

// 测试交易的编码和解码功能。
func TestTxEncodeAndDecode(t *testing.T) {
	tx := randomTxWithSignature(t) // 创建一个已签名的随机交易
	buf := &bytes.Buffer{} // 创建一个缓冲区用于编码
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf))) // 断言编码操作不返回错误

	txDecoded := new(Transaction) // 创建一个新的交易用于解码
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf))) // 断言解码操作不返回错误

	assert.Equal(t, tx, txDecoded) // 断言编码后解码的交易与原交易相等
}

// 创建一个已签名的随机交易。
func randomTxWithSignature(t *testing.T) *Transaction {
	priKey := crypto.GeneratePrivatekey() // 生成一个私钥
	tx := &Transaction{
		Data: []byte("foo"), // 创建一个交易，数据为"foo"
	}
	assert.Nil(t, tx.Sign(priKey)) // 断言签名操作不返回错误

	return tx // 返回已签名的交易
}
