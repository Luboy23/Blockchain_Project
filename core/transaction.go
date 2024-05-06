package core

import (
	"fmt"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
)

// 定义交易的结构体，包括交易数据、发送者的公钥、签名、哈希值和首次出现的时间戳。
type Transaction struct {
	Data []byte // 交易数据

	From      crypto.PublicKey // 发送者的公钥
	Signature *crypto.Signature // 交易签名

	hash      types.Hash // 交易的哈希值
	firstSeen int64 // 交易首次出现的时间戳
}

//	创建一个新的交易。
func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

//	计算交易的哈希值。如果交易的哈希值已经计算过，则直接返回；否则，使用提供的哈希计算器计算哈希值。
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() { // 如果交易的哈希值未计算或为零
		tx.hash = hasher.Hash(tx) // 使用哈希计算器计算哈希值
	}

	return tx.hash // 返回交易的哈希值
}

//	使用私钥对交易进行签名。
func (tx *Transaction) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(tx.Data) // 使用私钥对交易数据进行签名
	if err != nil {
		return err // 如果签名失败，返回错误
	}

	tx.From = priKey.PublicKey() // 设置发送者的公钥
	tx.Signature = sig // 设置签名

	return nil // 返回nil表示签名成功
}

//	验证交易的签名。
func (tx *Transaction) Verify() error {
	if tx.Signature == nil { // 如果交易没有签名
		return fmt.Errorf("交易没有签名！") // 返回错误
	}

	if !tx.Signature.Verify(tx.From, tx.Data) { // 如果签名验证失败
		return fmt.Errorf("不是交易的签名者") // 返回错误
	}

	return nil // 返回nil表示验证成功
}
//	使用提供的编码器对交易进行编码。
func (tx *Transaction) Encode(dec Encoder[*Transaction]) error {
	return dec.Encode(tx) // 使用编码器编码交易
}

//	使用提供的解码器对交易进行解码。
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx) // 使用解码器解码交易
}

// 设置交易首次出现的时间戳。
func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t // 设置首次出现的时间戳
}

// 获取交易首次出现的时间戳。
func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen // 返回首次出现的时间戳
}
