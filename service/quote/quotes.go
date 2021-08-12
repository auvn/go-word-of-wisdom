package quote

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
)

func MustLoadEntries() []Entry {
	entries, err := LoadEntries()
	if err != nil {
		panic(err)
	}

	return entries
}

func LoadEntries() ([]Entry, error) {
	var entries []Entry
	var entry Entry
	r := bufio.NewReader(bytes.NewBuffer(quotesFile))
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("ReadLine: %w", err)
		}

		entry = append(entry, line...)

		if isPrefix {
			continue
		}

		entries = append(entries, entry)
		entry = nil
	}

	return entries, nil
}

//go:embed fixtures/quotes.txt
var quotesFile []byte
