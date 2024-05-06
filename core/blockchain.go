package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

//	定义区块链的结构，包括存储、锁、区块头列表和验证器。
type Blockchain struct {
	store     Storage // 区块链的存储
	lock      sync.RWMutex // 用于同步访问区块链的锁
	headers   []*Header // 区块链中的所有区块头
	validator Validator // 用于验证区块的验证器
}

//	创建一个新的区块链，初始化存储和验证器，并添加创世区块。
func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{}, // 初始化区块头列表
		store:   NewMemstore(), // 初始化内存存储
	}

	bc.validator = NewBlockValidator(bc) // 初始化验证器

	err := bc.addBlockWithoutValidation(genesis) // 添加创世区块

	return bc, err // 返回新创建的区块链
}

//	设置区块链的验证器。
func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v // 设置新的验证器
}

//	获取区块链的高度。
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock() // 获取读锁
	defer bc.lock.RUnlock() // 确保在函数返回时释放读锁
	return uint32(len(bc.headers) - 1) // 返回区块链的高度
}

//	检查指定高度的区块是否存在于区块链中。
func (bc *Blockchain) HashBlock(height uint32) bool {
	return height <= bc.Height()// 如果指定高度小于或等于区块链高度，则返回true
}

//	向区块链中添加一个新的区块。
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil { // 验证区块
		return err // 如果验证失败，返回错误
	}

	return bc.addBlockWithoutValidation(b) // 添加区块
}

//	获取指定高度的区块头。
func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() { // 如果指定高度大于区块链高度
		return nil, fmt.Errorf("区块高度： %d 过高", height) // 返回错误
	}
	bc.lock.Lock() // 获取写锁
	defer bc.lock.Unlock() // 确保在函数返回时释放写锁

	return bc.headers[height], nil // 返回指定高度的区块头
}

//	 添加一个新的区块到区块链中，不进行验证。
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock() // 获取写锁
	bc.headers = append(bc.headers, b.Header) // 将新区块的头添加到区块头列表
	bc.lock.Unlock() // 释放写锁

	logrus.WithFields(logrus.Fields{ // 记录日志
		"区块高度": b.Height,
		"区块哈希": b.Hash(BlockHasher{}),
	}).Info("添加了一个新的区块")

	return bc.store.Put(b) // 将新区块存储到存储中
}
