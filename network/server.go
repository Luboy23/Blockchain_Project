package network

import (
	"fmt"
	"time"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlocktime = 5 * time.Second

// ServerOpts 结构体定义了一个服务器的配置选项。
// 它包含一个名为Transports的切片，用于存储服务器可以使用的传输方式。
type ServerOpts struct {
	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey

}

// Server 结构体定义了一个服务器的抽象。
// 它可能包含一些用于处理网络请求的方法和字段。
// 目前，这个结构体是空的，可能是因为它是一个接口或者它的实现还在开发中。
type Server struct {
	ServerOpts
	blockTime time.Duration
	memPool *TxPool
	isValidator bool
	rpcCh chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlocktime
	}
	return &Server{
		ServerOpts: opts,
		blockTime: opts.BlockTime,
		memPool: NewTxPool( ),
		isValidator: opts.PrivateKey != nil,
		rpcCh: make(chan RPC),
		quitCh:  make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <- s.rpcCh:
			fmt.Printf("RPC 为 %+v \n", rpc)
		case <-s.quitCh:
			break free
		case <- ticker.C:
			if s.isValidator{
				s.createNewBlock()
			}
		}
	}

	fmt.Println("服务器关闭")
}
func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify()
	err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"哈希为：":hash,
		}).Info("交易已经子啊内存池里了")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"哈希为：":hash,
	}).Info("添加了一笔新的交到内存池")

	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("创建一个新的区块")

	return nil
}

func (s * Server) initTransports() {
	for _, tr := range s.Transports{
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				//	handle
				s.rpcCh <- rpc
			}
		}(tr)
	}
}