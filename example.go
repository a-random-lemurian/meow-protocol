package main

import (
	"encoding/binary"
	"os"

	"github.com/a-random-lemurian/meow-protocol/meowproto"
)

func main() {
	meow := meowproto.MeowProtocolMessage{
		Version: 1,
		MessageType: meowproto.MtMeow,
		AnimalType: meowproto.AtCat,
		Breed: meowproto.BrCalico,
		Cuteness: 2,
		Name: "Ming",
	}

	bytes, err := meow.ToBytes()
	if err != nil {
		panic(err)
	}

	binary.Write(os.Stdout, binary.BigEndian, bytes)
}
