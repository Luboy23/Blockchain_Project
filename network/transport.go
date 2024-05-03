package network

// NetAddr 是一个类型别名，它表示网络地址。
// 在这个上下文中，它是一个字符串类型，用于表示网络地址。
type NetAddr string

// RPC 结构体定义了一个远程过程调用（RPC）的数据结构。
// 它包含一个From字段，表示消息的发送者；
// 一个Payload字段，表示消息的内容；
// 这个结构体可能用于在网络中传输消息。
type RPC struct {
	From string
	Payload []byte
}

// Transport 接口定义了一个网络传输的抽象。
// 它包含了一些方法，用于处理网络通信的基本操作，如接收消息、连接到其他传输、发送消息和获取传输的地址。
type Transport interface {

	// Consume方法返回一个接收RPC的通道。
	Consume() <- chan RPC

	// Connect方法用于连接到另一个Transport。
	Connect(Transport) error

	// SendMessage方法用于向指定地址发送消息。
	SendMessage(NetAddr, []byte) error

	// Addr方法返回当前Transport的地址。
	Addr() NetAddr
}