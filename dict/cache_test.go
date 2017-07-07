package dict

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func TestCache_Media(t *testing.T) {
	db, err := bolt.Open("/tmp/dix.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cache := &Cache{
		DB: db,
	}

	data, err := cache.Media(context.TODO(), "https://dictionary.blob.core.chinacloudapi.cn/media/audio/tom/6a/58/6A580F1D8DA8190EC9A1E7E5E60D75CB.mp3")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(data))
}
