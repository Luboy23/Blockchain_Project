package core

import "fmt"

// 定义了一个区块验证器接口，包含一个ValidateBlock方法，用于验证区块。
type Validator interface {
	ValidateBlock(*Block) error
}

// 实现了Validator接口，提供了一个区块验证的实现。
type BlockValidator struct {
	bc *Blockchain // 区块链的引用
}

// 创建一个新的BlockValidator实例。
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc, // 初始化区块链引用
	}
}

// 实现了Validator接口的ValidateBlock方法，用于验证区块。
func (v *BlockValidator) ValidateBlock(b *Block) error {
	// 检查区块是否已经存在于区块链中
	if v.bc.HashBlock(b.Height) {
		return fmt.Errorf("区块内已经包含了区块（%d）以及哈希（%s）", b.Height, b.Hash(BlockHasher{}))
	}

	// 检查区块的高度是否正确
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("(%s)区块过高", b.Hash(BlockHasher{}))
	}

	// 获取前一个区块的头部信息
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err // 如果获取失败，返回错误
	}

	// 计算前一个区块的哈希值
	hash := BlockHasher{}.Hash(prevHeader)

	// 检查前一个区块的哈希值是否正确
	if hash != b.PrevBlockHash {
		return fmt.Errorf("前区块哈希(%s)不正确！", b.PrevBlockHash)
	}

	// 验证区块的内容
	if err := b.Verify(); err != nil {
		return err // 如果验证失败，返回错误
	}

	return nil // 如果所有检查都通过，返回nil表示验证成功
}
