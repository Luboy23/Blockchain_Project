package network

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// //	测试两个LocalTransport实例之间的连接功能。
// //	创建两个LocalTransport实例，分别表示两个节点，然后测试它们之间的连接是否成功。
// func TestConnect(t *testing.T) {
// 	// 创建一个LocalTransport实例，地址为"A"。
// 	tra := NewLocalTransport("A")

// 	// 创建一个LocalTransport实例，地址为"B"。
// 	trb := NewLocalTransport("B")

// 	//	尝试将"A"连接到"B"。
// 	tra.Connect(trb)
// 	// 尝试将"B"连接到"A"。
// 	trb.Connect(tra)

// 	// 断言"A"的peers映射中包含"B"，并且"B"的地址与"A"的实例相同。
// 	assert.Equal(t,tra.peers[trb.addr], trb)

// 	// 断言"B"的peers映射中包含"A"，并且"A"的地址与"B"的实例相同。
// 	assert.Equal(t,trb.peers[tra.addr], tra)
// }

// // 测试在两个LocalTransport实例之间发送消息的功能。
// // 创建两个LocalTransport实例，分别表示两个节点，然后测试从一个节点发送消息到另一个节点的过程。
// func TestSendMessage(t *testing.T) {
// 	// 创建一个LocalTransport实例，地址为"A"。
// 	tra := NewLocalTransport("A")

// 	// 创建一个LocalTransport实例，地址为"B"。
// 	trb := NewLocalTransport("B")

// 	//	尝试将"A"连接到"B"。
// 	tra.Connect(trb)
// 	// 尝试将"B"连接到"A"。
// 	trb.Connect(tra)

// 	// 定义要发送的消息。
// 	msg := []byte("hello world")
// 	// 断言从"A"发送消息到"B"不返回错误。
// 	assert.Nil(t,tra.SendMessage(trb.addr,msg))

// 	// 从"B"的consumeCh通道中接收一个RPC消息。
// 	rpc := <- trb.Consume()

// 	// 断言接收到的消息的负载与发送的消息相同。
// 	assert.Equal(t, rpc.Payload, msg)

// 	// 断言接收到的消息的发送者地址与"A"的地址相同。
// 	assert.Equal(t, rpc.From, tra.addr)
// }