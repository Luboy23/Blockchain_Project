package network

import (
	"bytes"
	"fmt"
	"sync"
)

// 代表一个本地传输实现。
// 包含一个地址（addr）、一个用于接收RPC的通道（consumeCh）、一个读写锁（lock）和一个用于存储对等传输的映射（peers）。
type LocalTransport struct {

	// 本地传输的地址。
	addr NetAddr

	// 用于接收RPC的通道。
	consumeCh chan RPC

	// 读写锁，用于保护对peers映射的并发访问。
	lock sync.RWMutex

	// 映射，用于存储对等传输。
	peers map[NetAddr]*LocalTransport
}

// 创建并返回一个新的LocalTransport实例。
// 受一个NetAddr类型的地址作为参数，并初始化LocalTransport的consumeCh和peers字段。
func NewLocalTransport(addr NetAddr) *LocalTransport {
	return &LocalTransport{
		addr:	addr,

		// 创建一个容量为1024的RPC通道。
		consumeCh: make(chan RPC, 1024),

		// 初始化peers映射。
		peers: make(map[NetAddr]*LocalTransport),
	}
}

// 返回LocalTransport的consumeCh通道，用于接收RPC。
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// 用于连接到另一个LocalTransport实例。
// 接受一个*LocalTransport类型的参数，并将其添加到当前实例的peers映射中。
func (t *LocalTransport) Connect(tr Transport) error {
	 // 获取写锁，保护对peers映射的并发访问。
	t.lock.Lock()

	// 确保锁在函数返回时被释放。
	defer t.lock.Unlock()

	// 将传入的LocalTransport实例添加到peers映射中。
	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

// 用于向指定地址发送消息。
// 接受一个NetAddr类型的地址和一个字节切片作为负载，然后将消息发送到对应的peer。
func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {

	// 获取读锁，保护对peers映射的并发访问。
	t.lock.RLock()
	// 确保锁在函数返回时被释放。
	defer t.lock.RUnlock()

	// 尝试从peers映射中获取对应的peer。
	peer, ok := t.peers[to]

	// 如果找不到对应的peer，返回错误。
	if !ok {
		return fmt.Errorf("%s: 无法发送消息至： %s", t.addr,to)
	}

	// 将消息发送到对应的peer。
	peer.consumeCh <- RPC {

		// 消息的发送者。
		From:	string(t.addr),

		// 消息的负载。
		Payload: bytes.NewReader(payload),
	}
	return nil
}

// 用于向所有连接的peer广播消息。
func(t *LocalTransport) Broadcast (payload []byte) error {
	for _, peer := range t.peers{
		if err := t.SendMessage(peer.Addr(), payload)
		err != nil {
			return err
		}
	}
	return nil
}

// 返回LocalTransport的地址。
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
