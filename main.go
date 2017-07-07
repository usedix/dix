package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/go-vgo/robotgo"
	"github.com/usedix/dix/dict"
)

var (
	Player = make(chan string)
)

func main() {
	playagain := make(chan bool)
	go func(playagain chan<- bool) {
		defer close(playagain)
		for {
			mleft := robotgo.AddEvent("`")
			if mleft == 0 {
				playagain <- true
			}
		}
	}(playagain)

	db, err := bolt.Open(".cache.dix.boltdb", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cache := &dict.Cache{
		DB:         db,
		Dictionary: &dict.Bing{},
	}

	go func(player <-chan string, cache *dict.Cache) {
		var (
			ctx    = context.Background()
			cancel = func() {}
		)
		for url := range player {
			cancel()

			ctx, cancel = context.WithCancel(context.Background())

			data, err := cache.Media(ctx, url)
			if err != nil {
				continue
			}

			mediapath := "/tmp/dix.mp3"
			err = ioutil.WriteFile(mediapath, data, 0600)
			if err != nil {
				continue
			}

			go exec.CommandContext(ctx, "ffplay", "-nodisp", mediapath).Run()
		}
	}(Player, cache)

	words := Words(cache)
	for {
		select {
		case word, ok := <-words:
			if !ok {
				return
			}

			setCurrentWord(*word)
			// print word
			fmt.Print("\x1b[0;0H \x1b[0J")
			fmt.Printf("\x1b[32;1m%s\x1b[0m", word.Word)
			fmt.Printf("\x1b[36;3m[%s]\n\x1b[0m", word.Pronunciation.US)
			for _, def := range word.Defs {
				fmt.Printf("\x1b[34m%s\x1b[0m", def.PartOfSpeech)
				fmt.Printf("\x1b[30;1m%s\n\x1b[0m", def.Def)
			}
			//			fmt.Println("\x1b[35m_____________________________________________________________________\x1b[0m")
		case ok := <-playagain:
			if !ok {
				continue
			}
			setCurrentWord(CurrentWord())
		}
	}
}

var (
	currentWord   dict.Word
	currentWordMu = new(sync.Mutex)
)

func CurrentWord() dict.Word {
	currentWordMu.Lock()
	defer currentWordMu.Unlock()

	return currentWord
}

func setCurrentWord(word dict.Word) {
	currentWordMu.Lock()
	defer currentWordMu.Unlock()

	currentWord = word
	go func() { Player <- word.Pronunciation.US_MP3URL }()
}
