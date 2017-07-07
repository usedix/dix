package clip

import (
	"bytes"
	"os/exec"
)

func Primary() string {
	return Xclip("primary")
}

func Xclip(selection string) string {
	buf := new(bytes.Buffer)

	cmd := exec.Command("xclip", "-o", "-selection", selection)
	cmd.Stdout = buf

	cmd.Run()

	return buf.String()
}
