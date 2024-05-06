package network

import (
	"sort"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/Luboy23/Blockchain_Project/types"
)

// 用于对交易进行排序。
// 它包含一个交易切片，用于存储需要排序的交易。
type TxMapSorter struct {
	transactions []*core.Transaction
}

// 创建一个新的TxMapSorter实例。
// 它接收一个交易映射，将其转换为切片，并对切片进行排序。
func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap)) // 创建一个新的切片，长度等于交易映射的长度

	i := 0
	for _, val := range txMap { // 遍历交易映射，将每个交易添加到切片中
		txx[i] = val
		i++
	}

	s := &TxMapSorter{txx} // 创建一个新的TxMapSorter实例

	sort.Sort(s) // 对切片进行排序

	return s // 返回排序后的TxMapSorter实例
}

// 返回交易切片的长度。
// 它实现了sort.Interface接口的一部分。
func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

// 方法交换切片中两个交易的位置。
// 它实现了sort.Interface接口的一部分。
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

// 方法比较两个交易的首次见到的时间戳，返回true表示第一个交易的时间戳小于第二个交易的时间戳。
// 它实现了sort.Interface接口的一部分。
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

// 结构体代表一个交易池，用于存储待处理的交易。
// 它包含一个交易映射，用于存储交易的哈希值和交易对象。
type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

// 函数创建一个新的交易池实例。
// 它初始化交易映射，并返回新创建的交易池实例。
func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// 方法返回交易池中的所有交易，按照首次见到的时间戳排序。
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions) // 创建一个新的TxMapSorter实例
	return s.transactions // 返回排序后的交易切片
}

// 方法将一个新的交易添加到交易池中。
// 它接收一个交易对象，计算其哈希值，并将交易添加到交易映射中。
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{}) // 计算交易的哈希值

	p.transactions[hash] = tx // 将交易添加到交易映射中

	return nil // 返回nil，表示添加成功
}

// 检查交易池中是否存在指定的交易。
// 它接收一个交易的哈希值，并返回true表示交易存在。
func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash] // 检查交易映射中是否存在指定的哈希值
	return ok // 返回检查结果
}

// 返回交易池中的交易数量。
func (p *TxPool) Len() int {
	return len(p.transactions) // 返回交易映射的长度
}

//	清空交易池中的所有交易。
func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction) // 创建一个新的空交易映射
}
