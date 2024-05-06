package network 

import (
	"bytes"
	"encoding/gob" 
	"fmt" 
	"io" 

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/sirupsen/logrus" 
)

// 定义了一个表示消息类型的类型，基于byte，用于区分不同类型的消息
type MessageType byte

// 定义了消息类型的常量
const (
	MessageTypeTx MessageType = 0x1 // 定义了一个表示交易消息的常量
	MessageTypeBlock // 定义了一个表示区块消息的常量
)

// 定义了一个RPC结构体，包含发送者和消息负载，用于表示一个远程过程调用
type RPC struct {
	From string // 发送者的地址
	Payload io.Reader // 消息的负载，通过io.Reader接口读取
}

// 定义了一个消息结构体，包含消息头和数据，用于表示一个网络消息
type Message struct {
	Header MessageType // 消息头，表示消息的类型
	Data []byte // 消息数据，存储实际的消息内容
}

// 创建一个新的消息，接受消息类型和数据作为参数
func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t, // 设置消息头
		Data:	data, // 设置消息数据
	}
}

// 将消息编码为字节切片，用于网络传输
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{} // 创建一个缓冲区
	gob.NewEncoder(buf).Encode(msg) // 使用gob编码器将消息编码到缓冲区
	return buf.Bytes() // 返回编码后的字节切片
}

// 定义了一个解码后的消息结构体，包含发送者和数据，用于表示解码后的网络消息
type DecodeMessage struct {
	From NetAddr // 发送者的地址
	Data any // 消息数据，可以是任何类型
}

// 定义了一个RPC解码函数类型，接受一个RPC结构体，返回一个解码后的消息和错误信息
type RPCDecodeFunc func(RPC) (*DecodeMessage, error)

// 默认的RPC解码函数，用于解码RPC消息
func DefaultRPCDecodeFunc (rpc RPC) (*DecodeMessage, error) {
	msg	:= Message{} // 创建一个空的消息结构体
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil { // 使用gob解码器解码消息
		return nil, fmt.Errorf("从 %s 解码信息失败： %s",rpc.From,err) // 如果解码失败，返回错误
	}

	logrus.WithFields(logrus.Fields{ // 记录日志，包含消息类型和发送者
		"type": msg.Header,
		"from":rpc.From,
	}).Debug("传入一条新的消息")

	switch msg.Header { // 根据消息类型处理
	case MessageTypeTx: // 如果是交易消息
		tx := new(core.Transaction) // 创建一个新的交易结构体
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil { // 使用gob解码器解码交易数据
			return nil, err // 如果解码失败，返回错误
		}

		return &DecodeMessage{ // 返回解码后的消息
			From :NetAddr(rpc.From),
			Data: tx,
		}, nil
	default: // 如果是其他类型的消息
		return nil, fmt.Errorf("不正确的消息类型 % x", msg.Header) // 返回错误
	}
}

// 创建一个默认的RPC处理器，接受一个RPCProcessor接口作为参数
func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{ // 返回一个新的DefaultRPCHandler实例
		p: p, // 设置RPCProcessor接口
	}
}

// 定义了一个默认的RPC处理器结构体，用于处理RPC请求
type DefaultRPCHandler	struct {
	p RPCProcessor // RPCProcessor接口，用于处理消息
}

// 定义了一个RPC处理器接口，包含一个处理消息的方法
type RPCProcessor interface {
	ProcessMessage(*DecodeMessage) error // 处理消息的方法，接受一个解码后的消息作为参数，返回错误信息
}
