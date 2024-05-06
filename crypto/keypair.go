package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Luboy23/Blockchain_Project/types"
)

// 定义了一个私钥结构体，包含一个ecdsa.PrivateKey。
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// 定义了一个公钥结构体，包含一个ecdsa.PublicKey。
type PublicKey struct {
	Key *ecdsa.PublicKey
}

// 定义了一个签名结构体，包含两个大整数R和S。
type Signature struct {
	R *big.Int
	S *big.Int
}

// 使用私钥对数据进行签名，返回一个签名结构体。
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data) // 使用私钥对数据进行签名
	if err != nil {
		return nil, err // 如果签名失败，返回错误
	}

	return &Signature{ // 返回签名结构体
		R: r,
		S: s,
	}, nil
}

// 生成一个新的私钥。
func GeneratePrivatekey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // 生成一个新的ecdsa私钥
	if err != nil {
		panic(err) // 如果生成失败，抛出异常
	}
	return PrivateKey{ // 返回新生成的私钥
		key: key,
	}
}

// 返回与私钥对应的公钥。
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{ // 返回公钥结构体
		Key: &k.key.PublicKey,
	}
}

// 将公钥转换为字节切片。
func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y) // 使用椭圆曲线压缩格式将公钥转换为字节切片
}

// 计算公钥对应的地址。
func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice()) // 计算公钥的哈希值

	return types.AddressFromBytes(h[len(h)-20:]) // 返回地址
}

// 验证签名是否有效。
func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R, sig.S) // 使用公钥和数据验证签名
}
