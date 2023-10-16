package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"hw1/api-service/model"

	"net/http"
)

const bucket = "hw1-pic.s3.ir"

func SetupRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/register", newUser)
	e.GET("/status", getUserStatus)

	return e
}

func newUser(c echo.Context) error {
	//listMyBuckets(sess)
	//model.PingDB(model.DB)

	name := c.FormValue("name")
	email := c.FormValue("email")
	nationalId := c.FormValue("nationalId")
	ip := c.RealIP()

	image1, err := c.FormFile("image1")
	if err != nil {
		log.Warnln("image1 is broken")
		return c.JSON(http.StatusBadRequest, "Unable to open file")
	}
	path1 := UploadS3(model.Res.S3Sess, image1, bucket, nationalId)

	image2, err := c.FormFile("image2")
	if err != nil {
		log.Warnln("image2 is broken")
		return c.JSON(http.StatusBadRequest, "Unable to open file")
	}
	path2 := UploadS3(model.Res.S3Sess, image2, bucket, nationalId)

	err = model.Insert(name, email, nationalId, ip, path1, path2)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "user already exist")
	}

	err = WriteMQ(model.Res.RabbitConnection, nationalId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "cant write to MQ")
	}
	//log.Println("users:", model.GetAll())

	user := model.User{
		Name:       name,
		Email:      email,
		NationalId: nationalId,
		IP:         ip,
		Image1:     path1,
		Image2:     path2,
		State:      "checking",
	}

	return c.JSON(http.StatusCreated, user)
}

func getUserStatus(c echo.Context) error {
	nationalId := c.QueryParam("id")

	user := model.Find(nationalId)
	if user.IP != c.RealIP() {
		return c.String(http.StatusUnauthorized, "unauthorized")
	}

	return c.String(http.StatusOK, "You are in "+user.State+" state")
}
