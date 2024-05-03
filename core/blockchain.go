package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

//	State => State
//
//

type Blockchain struct {
	store Storage
	lock 	sync.RWMutex
	headers []*Header
	validator Validator
}

func NewBlockchain(genesis *Block)( *Blockchain , error){
	bc := &Blockchain {
		headers: []*Header{},
		store: 	NewMemstore(),
	}

	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) HashBlock (height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b) 
	err  != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil,fmt.Errorf("区块高度： %d 过高", height)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"区块高度": b.Height,
		"区块哈希":b.Hash(BlockHasher{}),
	}).Info("添加了一个新的区块")

	return bc.store.Put(b)
}

