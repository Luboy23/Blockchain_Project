package main

import (
	"time"

	"github.com/Luboy23/Blockchain_Project/network"
)

//	服务器
//	传输层 ->tcp,udp
//	区块
//	交易
//	密钥对

func main() {

trLocal := network.NewLocalTransport("LOCAL")
trRemote := network.NewLocalTransport("REMOTE")

trLocal.Connect(trRemote)
trRemote.Connect(trLocal)

go func() {
	for{
	trRemote.SendMessage(trLocal.Addr(), []byte("Hello World"))
		time.Sleep(1 * time.Second)
	}
}()

	opts:= network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}