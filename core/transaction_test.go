package core

import (
	"bytes"
	"testing"

	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	priKey := crypto.GeneratePrivatekey()
	tx := &Transaction{
		Data:	[]byte("foo"),
	}

	assert.Nil(t, tx.Sign(priKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	priKey := crypto.GeneratePrivatekey()
	tx := &Transaction{
		Data:	[]byte("foo"),
	}

	assert.Nil(t, tx.Sign(priKey))
	assert.Nil(t, tx.Verify())

	priKey1 := crypto.GeneratePrivatekey()
	tx.From = priKey1.PublicKey()

	assert.NotNil(t, tx.Verify())

}
  
func TestTxEncodeAndDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t,tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))

	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *Transaction {

	priKey := crypto.GeneratePrivatekey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(priKey))

	return tx
}