package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/Luboy23/Blockchain_Project/core"
	"github.com/Luboy23/Blockchain_Project/crypto"
	"github.com/Luboy23/Blockchain_Project/network"
	"github.com/sirupsen/logrus"
)

// 主函数，程序的入口点
func main() {
	// 创建两个本地传输实例，分别命名为"LOCAL"和"REMOTE"
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	// 将"LOCAL"和"REMOTE"传输实例连接起来
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	// 启动一个goroutine，不断地发送交易
	go func() {
		for {
			// 发送交易，如果发生错误，记录错误信息
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			// 每次发送交易后，等待1秒
			time.Sleep(1 * time.Second)
		}
	}()

	// 创建服务器选项，包含一个传输实例列表
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	// 使用选项创建一个新的服务器实例
	s := network.NewServer(opts)
	// 启动服务器
	s.Start()
}

// 用于创建并发送一个新的交易
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	// 生成一个新的私钥
	priKey := crypto.GeneratePrivatekey()
	// 创建一个包含随机数据的新交易
	data := []byte(strconv.FormatInt(int64(rand.Intn(100000)), 10))
	tx := core.NewTransaction(data)

	// 使用私钥对交易进行签名
	tx.Sign(priKey)

	// 创建一个缓冲区，用于编码交易
	buf := &bytes.Buffer{}

	// 将交易编码为字节切片，如果发生错误，返回错误
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	// 创建一个新的网络消息，包含交易数据
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	// 通过指定的传输实例将消息发送到指定的地址，如果发生错误，返回错误
	return tr.SendMessage(to, msg.Bytes())
}
