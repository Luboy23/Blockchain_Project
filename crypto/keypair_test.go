package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试私钥生成功能。
func TestGeneratePrivatekey(t *testing.T) {
	priKey := GeneratePrivatekey() // 生成一个新的私钥
	pubKey := priKey.PublicKey() // 获取私钥对应的公钥
	address := pubKey.Address() // 计算公钥对应的地址

	fmt.Println(address.String()) // 打印地址
}

// 测试私钥和公钥的签名和验证成功的情况。
func TestKeypairSignVerifySuccess(t *testing.T) {
	priKey := GeneratePrivatekey() // 生成一个新的私钥
	pubKey := priKey.PublicKey() // 获取私钥对应的公钥

	msg := []byte("hello,world") // 定义一个消息

	sig, err := priKey.Sign(msg) // 使用私钥对消息进行签名
	assert.Nil(t, err) // 断言签名操作不返回错误

	assert.True(t, sig.Verify(pubKey, msg)) // 断言签名验证成功
}

// 测试私钥和公钥的签名和验证失败的情况。
func TestKeypairSignVerifyFail(t *testing.T) {
	priKey1 := GeneratePrivatekey() // 生成第一个私钥
	pubKey1 := priKey1.PublicKey() // 获取第一个私钥对应的公钥
	msg1 := []byte("hello,world") // 定义第一个消息
	sig1, err1 := priKey1.Sign(msg1) // 使用第一个私钥对消息进行签名

	priKey2 := GeneratePrivatekey() // 生成第二个私钥
	pubKey2 := priKey2.PublicKey() // 获取第二个私钥对应的公钥
	msg2 := []byte("hello,world") // 定义第二个消息
	sig2, err2 := priKey2.Sign(msg2) // 使用第二个私钥对消息进行签名

	assert.Nil(t, err1) // 断言第一个签名操作不返回错误
	assert.Nil(t, err2) // 断言第二个签名操作不返回错误

	assert.True(t, sig1.Verify(pubKey1, msg1)) // 断言第一个签名验证成功
	assert.True(t, sig2.Verify(pubKey2, msg2)) // 断言第二个签名验证成功

	assert.False(t, sig1.Verify(pubKey2, msg2)) // 断言使用第二个公钥验证第一个签名失败
	assert.False(t, sig2.Verify(pubKey1, msg1)) // 断言使用第一个公钥验证第二个签名失败

	assert.False(t, sig1.Verify(pubKey1, []byte("ashdjkadhjsahdakjs"))) // 断言使用不同的消息验证第一个签名失败
}
