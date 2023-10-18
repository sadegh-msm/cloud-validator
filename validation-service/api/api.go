package api

import (
	log "github.com/sirupsen/logrus"
	"hw1/validation-service/configs"
)

func Validate() bool {
	userNationalId, err := ReadMQ()
	if err != nil {
		log.Warnln("cant read from queue")
		return false
	}

	user := Find(userNationalId)
	image1 := DownloadS3(Res.S3Sess, configs.Conf.S3Bucket, user.Image1)
	image2 := DownloadS3(Res.S3Sess, configs.Conf.S3Bucket, user.Image2)

	id1, err := faceDetection(image1)
	id2, err := faceDetection(image2)

	similarityScore := FaceSimilarity(id1, id2)

	if similarityScore >= 80 {
		Update(userNationalId, "accepted")
		user := Find(userNationalId)
		SendMail(user.State, user.Email, configs.Conf.MailGunDomain, configs.Conf.MailGunApiKey)
		log.Infof("person with %s ID has been accepted and informed by email.", userNationalId)
		return true
	} else {
		Update(userNationalId, "rejected")
		user := Find(userNationalId)
		SendMail(user.State, user.Email, configs.Conf.MailGunDomain, configs.Conf.MailGunApiKey)
		log.Infof("person with %s ID has been rejected and informed by email.", userNationalId)
		return false
	}
}
