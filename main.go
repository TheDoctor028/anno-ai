package main

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"time"
)

func main() {
	socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
		"/?EIO=3&transport=websocket&sid=v0-ZQfPRbL6XfwMnjFGe",
		5*time.Second)

}
