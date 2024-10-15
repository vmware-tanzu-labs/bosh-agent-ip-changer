package command

import log "github.com/sirupsen/logrus"

func initLogging(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
