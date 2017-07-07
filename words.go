package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/usedix/dix/clip"
	"github.com/usedix/dix/dict"
)

type Output struct {
	Word dict.Word
}

func Words(dictionary dict.Dictionary) <-chan *dict.Word {
	var (
		clipCheckInterval = time.Second / 10
		words             = make(chan *dict.Word)
	)

	go func(words chan<- *dict.Word) {
		defer close(words)

		lastWord := ""
		for {
			selectedWord := clip.Primary()
			if lastWord == selectedWord {
				time.Sleep(clipCheckInterval)
				continue
			}

			lastWord = selectedWord

			word, err := dictionary.Search(context.Background(), selectedWord)
			if err != nil {
				Errln(err)
			}

			words <- word
		}
	}(words)
	return words
}

func Errln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
