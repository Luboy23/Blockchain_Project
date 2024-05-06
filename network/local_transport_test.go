package network

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试连接功能
func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A") // 创建一个名为"A"的本地传输实例
	trb := NewLocalTransport("B") // 创建一个名为"B"的本地传输实例

	tra.Connect(trb) // "A"连接到"B"
	trb.Connect(tra) // "B"连接到"A"
	assert.Equal(t, tra.peers[trb.Addr()], trb) // 断言"A"的peers中包含"B"
	assert.Equal(t, trb.peers[tra.Addr()], tra) // 断言"B"的peers中包含"A"
}

// 测试发送消息功能
func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A") // 创建一个名为"A"的本地传输实例
	trb := NewLocalTransport("B") // 创建一个名为"B"的本地传输实例

	tra.Connect(trb) // "A"连接到"B"
	trb.Connect(tra) // "B"连接到"A"

	msg := []byte("hello world") // 定义一个消息"hello world"
	assert.Nil(t, tra.SendMessage(trb.addr, msg)) // "A"向"B"发送消息，断言没有错误

	rpc := <-trb.Consume() // 从"B"的消费队列中获取一个消息
	b, err := io.ReadAll(rpc.Payload) // 读取消息的负载
	assert.Nil(t, err) // 断言没有错误
	assert.Equal(t, b, msg) // 断言读取到的消息与发送的消息相同
	assert.Equal(t, rpc.From, string(tra.addr)) // 断言消息来源是"A"
}

// 测试广播消息功能
func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A") // 创建一个名为"A"的本地传输实例
	trb := NewLocalTransport("B") // 创建一个名为"B"的本地传输实例
	trc := NewLocalTransport("C") // 创建一个名为"C"的本地传输实例

	tra.Connect(trb) // "A"连接到"B"
	tra.Connect(trc) // "A"连接到"C"

	msg := []byte("foo") // 定义一个消息"foo"
	assert.Nil(t, tra.Broadcast(msg)) // "A"广播消息，断言没有错误

	rpcb := <-trb.Consume() // 从"B"的消费队列中获取一个消息
	b, err := io.ReadAll(rpcb.Payload) // 读取消息的负载
	assert.Nil(t, err) // 断言没有错误
	assert.Equal(t, b, msg) // 断言读取到的消息与广播的消息相同

	rpcC := <-trc.Consume() // 从"C"的消费队列中获取一个消息
	b, err = io.ReadAll(rpcC.Payload) // 读取消息的负载
	assert.Nil(t, err) // 断言没有错误
	assert.Equal(t, b, msg) // 断言读取到的消息与广播的消息相同
}
