package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"hw1/api-service/model"
	"net/http"
)

var bucket string

func SetupRouter() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/register", newUser)
	e.GET("/status", getUserStatus)

	return e
}

func newUser(c echo.Context) error {
	sess := ConnectS3()

	name := c.FormValue("name")
	email := c.FormValue("email")
	nationalId := c.FormValue("nationalId")
	ip := c.RealIP()

	image1, err := c.FormFile("image1")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to open file")
	}
	path1 := UploadS3(sess, image1, bucket, nationalId)

	image2, err := c.FormFile("image2")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to open file")
	}
	path2 := UploadS3(sess, image2, bucket, nationalId)

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
	//conn := model.ConnectMongo
	return nil
}
