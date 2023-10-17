package api

import (
	log "github.com/sirupsen/logrus"
)

const bucket = "hw1-pic.s3.ir"

func Validate() bool {
	userNationalId, err := ReadMQ()
	if err != nil {
		log.Warnln("cant read from queue")
		return false
	}

	user := Find(userNationalId)
	image1 := DownloadS3(Res.S3Sess, bucket, user.Image1)
	image2 := DownloadS3(Res.S3Sess, bucket, user.Image2)

	log.Infoln(image1)
	log.Infoln(image2)

	log.Infoln(user)

	return true
}
