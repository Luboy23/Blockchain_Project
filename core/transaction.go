package core

import (
	"fmt"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/types"
)

type Transaction struct {
	Data []byte 
	
	From crypto.PublicKey
	Signature *crypto.Signature

	hash types.Hash
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	

	return  tx.hash
}

func (tx *Transaction) Sign(priKey crypto.PrivateKey) error {
	sig, err := priKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = priKey.PublicKey()
	tx.Signature = sig

	return nil
} 

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("交易没有签名！")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("不是交易的签名者")
	}
	
	return nil 
}

func (tx *Transaction) Encode(dec Encoder[*Transaction]) error {
	return dec.Encode(tx)
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t 
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}