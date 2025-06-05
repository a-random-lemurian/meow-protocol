package meowproto

import (
	"bytes"
	"errors"
	"strings"

	"github.com/icza/bitio"
)

const (
	MtUnknown  uint64 = iota
	MtMeow            = 1
	MtPurr            = 2
	MtScratch         = 3
	MtBite            = 4
	MtPawAtYou        = 5
	MtGrowl           = 6
	MtHiss            = 7
)

const (
	AtCat uint64 = 1 << iota
	AtHumanM
	AtHumanF
	AtHumanUnknown
)

const (
	BrUnknown uint64 = iota
	BrCalico
	BrWhite
	BrSiamese
)

type Cuteness uint64

type MeowProtocolMessage struct {
	Version     uint64
	MessageType uint64
	AnimalType  uint64
	Breed       uint64
	Cuteness    uint64
	Name        string
}

func (msg *MeowProtocolMessage) ToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	w := bitio.NewWriter(&buf)

	w.WriteBits(msg.Version, 4)
	w.WriteBits(msg.Cuteness, 4)
	w.WriteBits(msg.MessageType, 8)
	w.WriteBits(msg.AnimalType, 8)
	w.WriteBits(msg.Breed, 16)
	w.WriteBits(uint64(len(msg.Name)), 8)

	nameBytes := []byte(msg.Name)
	nameBytes = append(nameBytes, 0)
	buf.Write(nameBytes)

	err := w.Close()
	return buf.Bytes(), err
}

var (
	ErrBadMessage    = errors.New("invalid message")
	ErrBadNameLength = errors.New("length of name does not actually match up")
)

// TODO error handling
func ReadMessage(b []byte) (*MeowProtocolMessage, error) {
	buf := bytes.NewBuffer(b)
	r := bitio.NewReader(buf)

	// version 1 has 46 bits at start, rounded up to 6 bytes
	// good heuristic for detecting message corruption
	if len(b) < 6 {
		return nil, ErrBadMessage
	}

	m := &MeowProtocolMessage{}

	version, err := r.ReadBits(4)
	if err != nil {
		return nil, err
	}
	m.Version = version

	m.Cuteness, err = r.ReadBits(4)
	if err != nil {
		return nil, err
	}

	m.MessageType, err = r.ReadBits(8)
	if err != nil {
		return nil, err
	}
	m.AnimalType, err = r.ReadBits(8)
	if err != nil {
		return nil, err
	}

	m.Breed, err = r.ReadBits(16)
	if err != nil {
		return nil, err
	}

	nameLen, err := r.ReadBits(8)
	if err != nil {
		return nil, err
	}

	name, err := buf.ReadString(0x00)
	name = strings.TrimRight(name, "\x00")

	if err != nil {
		return nil, err
	}

	if nameLen != uint64(len(name)) {
		return nil, ErrBadNameLength
	}

	m.Name = name

	return m, nil
}
