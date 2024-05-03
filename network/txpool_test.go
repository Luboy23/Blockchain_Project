package network

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t,p.Len(),0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("lllu"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	_ = core.NewTransaction([]byte("lllu"))
	assert.Equal(t, p.Len(),1)


	p.Flush()
	assert.Equal(t,p.Len(),0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen,p.Len())

	txx := p.Transactions()
	for i := 0; i < len(txx) - 1 ;i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i + 1].FirstSeen())
	}
}