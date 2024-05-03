package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivatekey(t *testing.T) {
	priKey := GeneratePrivatekey()
	pubKey := priKey.PublicKey()
	address := pubKey.Address()


	fmt.Println(address.String())

}

func TestKeypairSignVerifySuccess(t *testing.T) {
	priKey := GeneratePrivatekey()
	pubKey := priKey.PublicKey()

	msg := []byte("hello,world")

	sig, err := priKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t,sig.Verify(pubKey, msg))

}

func TestKeypairSignVerifyFail(t *testing.T) {
	priKey1 := GeneratePrivatekey()
	pubKey1 := priKey1.PublicKey()
	msg1 := []byte("hello,world")
	sig1, err1 := priKey1.Sign(msg1)

	priKey2 := GeneratePrivatekey()
	pubKey2 := priKey2.PublicKey()
	msg2 := []byte("hello,world")
	sig2, err2 := priKey2.Sign(msg2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	assert.True(t,sig1.Verify(pubKey1, msg1))
	assert.True(t,sig2.Verify(pubKey2, msg2))

	assert.False(t,sig1.Verify(pubKey2, msg2))
	assert.False(t,sig2.Verify(pubKey1, msg1))

	assert.False(t, sig1.Verify(pubKey1, []byte("ashdjkadhjsahdakjs")))



}