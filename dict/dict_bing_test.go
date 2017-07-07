package dict

import (
	"context"
	"fmt"
	"testing"
)

func TestBing_Search(t *testing.T) {
	bing := Bing{}

	word, err := bing.Search(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(word)
}
