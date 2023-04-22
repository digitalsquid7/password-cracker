package wordhasher

import (
	"fmt"
	"hash"
)

type Hasher struct {
	hasher hash.Hash
}

func (h *Hasher) Hash(word string) string {
	h.hasher.Reset()
	h.hasher.Write([]byte(word))
	return fmt.Sprintf("%x", h.hasher.Sum(nil))
}

func New(hash hash.Hash) *Hasher {
	return &Hasher{hash}
}
