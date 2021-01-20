package main

import (
	"gwnoo/logger"
	"gwnoo/message-service"
)

func main() {
	defer logger.Logger.Sync()
 	message_service.Start()
}