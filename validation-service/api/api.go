package api

import log "github.com/sirupsen/logrus"

func Validate() bool {
	userNationalId, err := ReadMQ()
	if err != nil {
		log.Warnln("cant read from queue")
		return false
	}

	user := Find(userNationalId)
	log.Infoln(user)

	return true
}
