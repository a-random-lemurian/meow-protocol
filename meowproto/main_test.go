package meowproto_test

import (
	"testing"

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
