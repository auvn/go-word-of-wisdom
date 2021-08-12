package quote

import (
	"math/rand"
	"time"
)

var DefaultProducer = NewRandomProducer(MustLoadEntries())

type RandomProducer struct {
	r *rand.Rand
	entries    []Entry
}

func NewRandomProducer(entries []Entry) *RandomProducer {
	if len(entries) == 0 {
		panic("empty entries")
	}

	return &RandomProducer{
		r: rand.New(rand.NewSource(time.Now().Unix())),
		entries:    entries,
	}
}

func (p *RandomProducer) Produce() (Entry, error) {
	r := p.r.Intn(len(p.entries))
	return p.entries[r], nil
}

type Entry = []byte
