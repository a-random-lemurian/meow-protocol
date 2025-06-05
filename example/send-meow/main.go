package main

import (
	"net"
	"time"

	"github.com/a-random-lemurian/meow-protocol/meowproto"
)

func main() {
	udpServer, err := net.ResolveUDPAddr("udp4", "localhost:32390")
	if err != nil {
		panic(err)
	}

	sock, err := net.DialUDP("udp4", nil, udpServer)
	if err != nil {
		panic(err)
	}

	meow := meowproto.MeowProtocolMessage{
		Version:     1,
		MessageType: meowproto.MtMeow,
		AnimalType:  meowproto.AtCat,
		Breed:       meowproto.BrCalico,
		Cuteness:    4,
		Name:        "Calico-148",
	}
	bytes, err := meow.ToBytes()
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(1 * time.Second)
	for range t.C {
		sock.Write(bytes)
	}
}
