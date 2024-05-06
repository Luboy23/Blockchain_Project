package network

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/stretchr/testify/assert"
)

// 测试创建新的交易池时，交易池的长度是否为0。
func TestTxPool(t *testing.T) {
	p := NewTxPool() // 创建一个新的交易池
	assert.Equal(t, p.Len(), 0) // 断言交易池的长度为0
}

// 测试向交易池添加交易的功能。
// 它检查添加交易后，交易池的长度是否正确增加，
// 以及添加重复交易时，交易池的长度是否不变。
func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool() // 创建一个新的交易池
	tx := core.NewTransaction([]byte("lllu")) // 创建一个新的交易
	assert.Nil(t, p.Add(tx)) // 断言添加交易后，返回值为nil
	assert.Equal(t, p.Len(), 1) // 断言交易池的长度为1

	_ = core.NewTransaction([]byte("lllu")) // 创建另一个相同的交易
	assert.Equal(t, p.Len(), 1) // 断言交易池的长度仍为1，因为添加了重复的交易

	p.Flush() // 清空交易池
	assert.Equal(t, p.Len(), 0) // 断言清空后，交易池的长度为0
}

// 测试交易池中交易的排序功能。
// 它创建一定数量的交易，并为每个交易设置不同的首次见到的时间戳，
// 然后检查交易是否按照时间戳正确排序。
func TestSortTransactions(t *testing.T) {
	p := NewTxPool() // 创建一个新的交易池
	txLen := 1000 // 定义要添加的交易数量

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10))) // 创建一个新的交易
		tx.SetFirstSeen(int64(i * rand.Intn(10000))) // 为交易设置随机的首次见到的时间戳
		assert.Nil(t, p.Add(tx)) // 断言添加交易后，返回值为nil
	}

	assert.Equal(t, txLen, p.Len()) // 断言交易池的长度等于添加的交易数量

	txx := p.Transactions() // 获取排序后的交易切片
	for i := 0; i < len(txx) - 1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i + 1].FirstSeen()) // 断言交易按照首次见到的时间戳正确排序
	}
}
