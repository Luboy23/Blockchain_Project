package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
)

//	定义区块头的结构，包括版本、数据哈希、前一个区块哈希、时间戳和高度。
type Header struct {
	Version       uint32     // 区块版本号
	DataHash      types.Hash // 区块中所有交易的哈希值
	PrevBlockHash types.Hash // 前一个区块的哈希值
	Timestamp     int64      // 区块创建的时间戳
	Height        uint32     // 区块在链上的高度
}

//	定义区块的结构，包括区块头、交易列表、验证者公钥、签名和哈希。
type Block struct {
	*Header                        // 指向区块头的指针
	Transactions []Transaction     // 区块中包含的交易列表
	Validator    crypto.PublicKey  // 验证区块的公钥
	Signature    *crypto.Signature // 区块的签名

	hash types.Hash // 区块的哈希值
}

//	将区块头转换为字节流，用于序列化。
func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}     // 创建一个缓冲区
	enc := gob.NewEncoder(buf) // 创建一个新的gob编码器
	enc.Encode(h)              // 将区块头编码到缓冲区

	return buf.Bytes() // 返回缓冲区中的字节流
}

//	计算区块的哈希值，如果哈希值未被计算过，则使用提供的哈希函数计算。
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() { // 检查哈希值是否为零
		b.hash = hasher.Hash(b.Header) // 使用提供的哈希函数计算哈希值
	}

	return b.hash // 返回计算得到的哈希值
}

//	将区块编码为字节流，用于序列化。
func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(b) // 使用提供的编码器将区块编码到输出流
}

//	从字节流中解码区块，用于反序列化。
func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(b) // 使用提供的解码器从输入流解码区块
}

//	向区块中添加一个交易。
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx) // 将交易添加到区块的交易列表中
}

//	创建一个新的区块，包含指定的区块头和交易列表。
func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{ // 返回一个新的区块实例
		Header:       h,   // 设置区块头
		Transactions: txx, // 设置交易列表
	}
}

//	为区块签名，使用提供的私钥对区块头进行签名，并设置验证者公钥和签名。
func (b *Block) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(b.Header.Bytes()) // 使用私钥对区块头进行签名
	if err != nil {                           // 检查签名过程中是否出错
		return err // 如果出错，返回错误
	}

	b.Validator = priKey.PublicKey() // 设置验证者公钥
	b.Signature = sig                // 设置签名

	return nil // 返回nil表示成功
}

//	验证区块的签名和交易，确保区块的完整性和有效性。
func (b *Block) Verify() error {
	if b.Signature == nil { // 检查区块是否有签名
		return fmt.Errorf("区块没有签名！") // 如果没有签名，返回错误
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) { // 验证签名是否有效
		return fmt.Errorf("区块签名不匹配！") // 如果签名无效，返回错误
	}

	for _, tx := range b.Transactions { // 遍历区块中的所有交易
		if err := tx.Verify(); err != nil { // 验证每个交易是否有效
			return err // 如果有无效的交易，返回错误
		}
	}

	return nil // 如果所有验证通过，返回nil表示成功
}
