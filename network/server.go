package network

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second // 定义了默认的区块生成时间间隔

// 定义了一个服务器的配置选项。
// 它包含一个名为Transports的切片，用于存储服务器可以使用的传输方式。
type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc // RPC解码函数
	RPCProcessor RPCProcessor // RPC处理器
	Transports []Transport // 传输方式
	BlockTime time.Duration // 区块生成时间间隔
	PrivateKey *crypto.PrivateKey // 私钥，用于验证交易
}

// 定义了一个服务器的抽象。
// 它可能包含一些用于处理网络请求的方法和字段。
type Server struct {
	ServerOpts // 服务器配置选项
	memPool *TxPool // 内存池，用于存储待处理的交易
	isValidator bool // 是否是验证者
	rpcCh chan RPC // RPC通道，用于接收RPC请求
	quitCh chan struct{} // 退出通道，用于接收退出信号
}

// 创建一个新的服务器实例
func NewServer(opts ServerOpts) *Server {
	// 如果没有指定区块生成时间间隔，则使用默认值
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	// 如果没有指定RPC解码函数，则使用默认的解码函数
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	// 创建一个新的服务器实例
	s := &Server{
		ServerOpts: opts,
		memPool: NewTxPool( ),
		isValidator: opts.PrivateKey != nil,
		rpcCh: make(chan RPC),
		quitCh: make(chan struct{}, 1),
	}

	// 如果没有指定RPC处理器，则使用服务器自身作为处理器
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	 return s
}

// 启动服务器
func (s *Server) Start() {
	s.initTransports() // 初始化传输方式

	ticker := time.NewTicker(s.BlockTime) // 创建一个定时器，用于定时生成新的区块

free:
	for {
		select {
		case rpc := <- s.rpcCh: // 接收RPC请求
			msg, err := s.RPCDecodeFunc(rpc) // 解码RPC请求
			if err != nil {
				logrus.Error(err) // 如果解码失败，记录错误
			}

			if err := s.RPCProcessor.ProcessMessage(msg) // 处理解码后的消息
			err != nil {
				logrus.Error(err) // 如果处理失败，记录错误
			}

		case <-s.quitCh: // 接收退出信号
			break free
		case <- ticker.C: // 接收定时器的信号
			if s.isValidator{ // 如果是验证者
				s.createNewBlock() // 创建新的区块
			}
		}
	}

	fmt.Println("服务器关闭") // 服务器关闭时打印消息
}

// 处理解码后的消息
func (s *Server) ProcessMessage(msg *DecodeMessage) error { 
	switch t := msg.Data.(type) { // 根据消息数据的类型进行处理
	case *core.Transaction : // 如果是交易
			return s.processTransaction(t) // 处理交易
	}
	return nil 
}

// 广播消息
func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports{ // 遍历所有的传输方式
		if err := tr.Broadcast(payload) // 广播消息
		err != nil {
			return nil // 如果广播失败，返回错误
		}
	}
	return nil // 广播成功，返回nil
}
// 处理交易的函数。
// 它接收一个交易对象，并执行一系列步骤来处理这个交易。
func (s *Server) processTransaction(tx *core.Transaction) error {
	// 使用core.TxHasher计算交易的哈希值。
	hash := tx.Hash(core.TxHasher{})

	// 检查内存池中是否已经存在这个交易。
	// 如果存在，则记录日志并返回nil，表示处理成功。
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"哈希为：":hash,
		}).Info("交易已经在内存池里了")
		return nil
	}

	// 验证交易的有效性。
	// 如果验证失败，返回错误。
	if err := tx.Verify(); err != nil {
		return err
	}

	// 设置交易的首次见到的时间戳。
	tx.SetFirstSeen(time.Now().UnixNano())
	// 记录日志，表示交易已经被添加到内存池。
	logrus.WithFields(logrus.Fields{
		"哈希为：":hash,
		"交易池的长度为：":s.memPool.Len(),
	}).Info("添加了一笔新的交到内存池")

	// 异步广播交易。
	go s.broadcastTx(tx)

	// 将交易添加到内存池。
	// 如果添加成功，返回nil。
	return s.memPool.Add(tx)
}

// 广播交易的函数。
// 它接收一个交易对象，将其编码为字节流，然后通过broadcast函数广播出去。
func (s *Server) broadcastTx(tx *core.Transaction) error {
	// 创建一个新的字节缓冲区，用于存储编码后的交易数据。
	buf := &bytes.Buffer{}
	// 使用core.NewGobTxEncoder将交易对象编码为字节流。
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}
	
	// 创建一个新的消息对象，包含交易类型和编码后的交易数据。
	msg := NewMessage(MessageTypeTx, buf.Bytes())

	// 调用broadcast函数，将消息广播到网络上。
	return s.broadcast(msg.Bytes())
}

// 创建新的区块的函数。
// 目前这个函数只是打印一条消息，表示正在创建新的区块。
func (s *Server) createNewBlock() error {
	// 打印一条消息，表示正在创建新的区块。
	fmt.Println("创建一个新的区块")

	// 返回nil，表示函数执行成功。
	return nil
}

// 初始化传输方式的函数。
// 它遍历服务器配置中的所有传输方式，并为每个传输方式启动一个goroutine来消费RPC请求。
func (s *Server) initTransports() {
	// 遍历服务器配置中的所有传输方式。
	for _, tr := range s.Transports{
		// 为每个传输方式启动一个goroutine。
		go func(tr Transport) {
			// 消费传输方式的RPC请求。
			for rpc := range tr.Consume() {
				// 将接收到的RPC请求发送到服务器的RPC通道。
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
