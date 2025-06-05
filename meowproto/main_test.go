package meowproto_test

import (
	"crypto/rand"
	"testing"
	"unsafe"

	"github.com/a-random-lemurian/meow-protocol/meowproto"
	"github.com/google/go-cmp/cmp"
)

var meow = &meowproto.MeowProtocolMessage{
	Version:     1,
	MessageType: meowproto.MtMeow,
	AnimalType:  meowproto.AtCat,
	Breed:       meowproto.BrCalico,
	Cuteness:    2,
	Name:        "Ming",
}

var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generate(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return *(*string)(unsafe.Pointer(&b))
}

func TestGoodMessage(t *testing.T) {
	bytes, err := meow.ToBytes()
	if err != nil {
		t.Error(err)
	}

	readback, err := meowproto.ReadMessage(bytes)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(meow, readback); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestSerializeLongName(t *testing.T) {
	var long = &meowproto.MeowProtocolMessage{
		Version:     1,
		MessageType: meowproto.MtMeow,
		AnimalType:  meowproto.AtCat,
		Breed:       meowproto.BrCalico,
		Cuteness:    2,
	}
	long.Name = generate(257)

	bytes, err := long.ToBytes()
	if err != meowproto.ErrNameTooLong {
		t.Errorf("We did not get an error, bad message bytes %+v", bytes)
	}
}

func BenchmarkSerialize(b *testing.B) {
	for b.Loop() {
		meow.ToBytes()
	}
}

func BenchmarkReadMessage(b *testing.B) {
	bytes, err := meow.ToBytes()
	if err != nil {
		b.Error(err)
	}

	for b.Loop() {
		_, err := meowproto.ReadMessage(bytes)
		if err != nil {
			b.Error(err)
		}
	}
}
