package core

//	定义一个存储接口，包含一个Put方法，用于将区块存储到存储系统中。
type Storage interface {
	Put(*Block) error
}

//	实现Storage接口，提供了一个内存存储的实现。
type MemoryStore struct {
}

//	创建一个新的MemoryStore实例。
func NewMemstore() *MemoryStore {
	return &MemoryStore{}
}

//	实现Storage接口的Put方法，用于将区块存储到内存中。
// 注意：这里的实现仅仅是一个占位符，实际上并没有将数据存储到任何地方。
func (s *MemoryStore) Put(b *Block) error {
	return nil //	返回nil表示操作成功
}
