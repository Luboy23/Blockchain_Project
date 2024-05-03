package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Luboy23/Blockchain_Project/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type PublicKey struct {
	Key *ecdsa.PublicKey

}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (k PrivateKey) Sign (data []byte) (*Signature, error ) {
	r, s, err :=ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		R:r, 
		S:s,
	},nil
}

func GeneratePrivatekey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
	if err != nil {
	panic(err)

	}
	return PrivateKey{
		key: key,
	}
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.key.PublicKey,
	}
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return  types.AddressFromBytes(h[len(h)-20:])
}

func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.Key, data, sig.R,sig.S)
}