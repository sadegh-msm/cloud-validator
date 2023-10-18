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
	//image2 := DownloadS3(Res.S3Sess, bucket, user.Image2)

	faceDetection(image1)

	//MGapiKey := "880e398409b13a654b0e5f564017f933-3750a53b-2bbe965e"
	//domain := "sandbox537fc23d9dfc4ff085da5c7b23074837.mailgun.org"

	//SendMail(user.State, user.Email, domain, MGapiKey)

	log.Infoln(user)

	return true
}
