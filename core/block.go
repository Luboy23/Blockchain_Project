package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
)

type Header struct {
	Version uint32
	DataHash types.Hash
	PrevBlockHash types.Hash
	Timestamp int64
	Height uint32 
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}
	

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}

// func (b *Block) Sign()

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(b)
}  

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header: h,
		Transactions: txx,
	}
}

func (b *Block) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(b.Header.Bytes())
	if err  != nil {
		return err
	}

	b.Validator = priKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("区块没有签名！")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("区块签名不匹配！")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify()
		err != nil {
			return err
		}
	}

	return nil 

}

