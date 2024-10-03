package main

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
)

func main() {
	sio := socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
	)
	for {
		select {
		case <-sio.Done:

		}
	}

}
