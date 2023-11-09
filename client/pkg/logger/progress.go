package logger

import "log"

func StartToEnd(f func(), startMsg, endMsg string) {
	log.Println(startMsg)
	f()
	log.Println(endMsg)
}
