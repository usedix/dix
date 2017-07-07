package clip

import (
	"fmt"
	"testing"
	"time"
)

func TestXclip(t *testing.T) {
	for {
		fmt.Print("\x1b[0;0H\x1b[0J")
		fmt.Println(Primary())
		time.Sleep(time.Second / 10)
	}
}
