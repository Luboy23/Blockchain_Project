package core

import (
	"crypto/sha256"

	"github.com/Luboy23/Blockchain_Project/types"
)

//	定义了一个哈希计算器接口，用于计算任意类型的数据的哈希值。
type Hasher[T any] interface {
	Hash(T) types.Hash
}

//	 实现Hasher接口，用于计算区块头（Header）的哈希值。
type BlockHasher struct {
}

//	使用BlockHasher计算区块头的哈希值。
func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes()) // 使用SHA-256算法计算区块头的字节表示的哈希值

	return types.Hash(h) // 将哈希值转换为types.Hash类型并返回
}

//	实现Hasher接口，用于计算交易（Transaction）的哈希值。
type TxHasher struct {
}

//	使用TxHasher计算交易的哈希值。
func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data)) // 使用SHA-256算法计算交易数据的哈希值
}
