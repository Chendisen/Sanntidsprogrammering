package stop_button

import (
	"os/exec"
)

func STOP() error{
	url := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	cmd := exec.Command("xdg-open", url)
	return cmd.Run()
}