package stop_button

import (
	"os/exec"
)

func STOP() error{
	url := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	js := "javascript:var video=document.querySelector('video');if(video){video.play();}else{window.location.href='" + url + "';}"
	cmd := exec.Command("xdg-open", url, js)
	return cmd.Run()
}