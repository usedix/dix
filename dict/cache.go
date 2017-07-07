package dict

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
)

var (
	BucketRawwords = []byte("rawwords")
	BucketWords    = []byte("words")
	BucketMedia    = []byte("media")
)

type Cache struct {
	Dictionary
	DB *bolt.DB
}

// TODO: enable ctx
func (d Cache) Media(ctx context.Context, url string) ([]byte, error) {
	// directly return if hinting cache
	var data []byte
	err := d.DB.View(func(tx *bolt.Tx) error {
		media := tx.Bucket(BucketMedia)
		if media == nil {
			return errors.New("not found bucket")
		}

		data = media.Get([]byte(url))
		if len(data) == 0 {
			return errors.New("not found key")
		}
		return nil
	})
	if err == nil {
		//		log.Print("[media]Hint: ", url)
		return data, nil
	}

	// download media via url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// cache media data
	d.DB.Update(func(tx *bolt.Tx) error {
		media, err := tx.CreateBucketIfNotExists(BucketMedia)
		if err != nil {
			return err
		}

		media.Put([]byte(url), []byte(data))
		//		log.Print("[media]Set: ", url)
		return nil
	})

	return data, nil
}

func (d Cache) Search(ctx context.Context, rawword string) (*Word, error) {
	// directly return if cache hinting
	var word = Word{}
	err := d.DB.View(func(tx *bolt.Tx) error {
		rawwords := tx.Bucket(BucketRawwords)
		if rawwords == nil {
			return errors.New("not found bucket")
		}
		wd := rawwords.Get([]byte(rawword))
		if len(wd) == 0 {
			return errors.New("not found key")
		}

		words := tx.Bucket(BucketWords)
		if words == nil {
			return errors.New("not found bucket")
		}

		cachedWord := words.Get(wd)
		if len(cachedWord) == 0 {
			return errors.New("not found key")
		}

		return json.Unmarshal(cachedWord, &word)
	})

	if err == nil {
		//log.Print("Hint:", rawword)
		return &word, nil
	}

	wd, err := d.Dictionary.Search(ctx, rawword)
	if err != nil {
		return nil, err
	}
	word = *wd

	// cache the result
	d.DB.Update(func(tx *bolt.Tx) error {
		rawwords, err := tx.CreateBucketIfNotExists(BucketRawwords)
		if err != nil {
			return err
		}
		words, err := tx.CreateBucketIfNotExists(BucketWords)
		if err != nil {
			return err
		}

		data, err := json.Marshal(word)
		if err != nil {
			return err
		}

		rawwords.Put([]byte(rawword), []byte(word.Word))
		words.Put([]byte(word.Word), data)
		//log.Print("Set:", rawword, " = ", word.Word)
		return nil
	})

	return &word, nil
}
