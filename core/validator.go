package core

import "fmt"

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc:	bc,
	}
}

func (v *BlockValidator) ValidateBlock (b *Block) error {
	if v.bc.HashBlock(b.Height) {
		return fmt.Errorf("区块内已经包含了区块（%d）以及哈希（%s）", b.Height,b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height() + 1 {
		return fmt.Errorf("(%s)区块过高",b.Hash(BlockHasher{})	)
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	} 

	hash := BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash{
		return	fmt.Errorf("前区块哈希(%s)不正确！", b.PrevBlockHash)
	}
 

	if err := b.Verify()
	err != nil {
		return err
	}



	return nil
}