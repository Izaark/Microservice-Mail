package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/tasks/Microservice-Mail/models"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

const ctimeLayout = "2006-01-02T15:04:05.000Z"

var strFromEmail string
var strPasswordEmail string

func checkError(strMessage string, err error) {

	if err != nil {
		err = errors.New(strMessage + "->" + err.Error())
		log.Println(err)
		fmt.Println(err)
	}
}

func MailRouter() {
	/*var osLogFile *os.File
	var err error

	osLogFile, err = os.Create(os.Getenv("ENV_LOG_FILE") + "." + time.Now().UTC().Format(ctimeLayout))
	defer osLogFile.Close()

	checkError("ERROR MailRouter: couldn't create log file", err)
	log.SetOutput(osLogFile)
	fmt.Println(osLogFile)*/

	router := gin.Default()
	//CORS
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE ,OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	//Endpoints Microservice
	mail := router.Group("api")
	{
		mail.POST("/users/email/text", handlerSendEmail)
	}

	router.Run(":" + os.Getenv("API_PORT"))
}

func handlerSendEmail(c *gin.Context) {
	var vEmail models.ObjUserInfo
	var vginResponse gin.H
	var err error

	idUser := vEmail.Id
	strToEmail := vEmail.Email
	strBody := vEmail.Body

	err = c.BindJSON(&vEmail)
	fmt.Println(vEmail)
	if err != nil {
		err = errors.New("ERROR handlerSendEmail: couldn't bind payload provided ObjUserInfo struct -> " + err.Error())
		vginResponse = gin.H{"message": "error reading payload provided", "response": nil, "error": "Response Error", "status": http.StatusBadRequest}
		c.JSON(http.StatusBadRequest, vginResponse)
		return
	}
	err = funSendEmail(idUser, strToEmail, strBody)
	if err != nil {
		err = errors.New("ERROR handlerSendEmail: couldn't init funSendEmail -> " + err.Error())
		vginResponse = gin.H{"message": "error reading payload provided", "response": nil, "error": "Response Error", "status": http.StatusBadRequest}
		c.JSON(http.StatusBadRequest, vginResponse)
		return
	}

	vginResponse = gin.H{"message": "Message send :D", "response": vEmail, "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}

func funSendEmail(idUser, strToEmail, strBody string) error {
	strFromEmail = os.Getenv("FROM_EMAIL")
	strPasswordEmail = os.Getenv("FROM_EMAIL_PASSWORD")

	message := "From: " + strFromEmail + "\n" + "To: " + strToEmail + "\n" + "Subject: Hello there\n\n" + strBody

	err := smtp.SendMail("smtp.gmail.com:465", smtp.PlainAuth("", strFromEmail, strPasswordEmail, "smtp.gmail.com"), strFromEmail, []string{strToEmail}, []byte(message))
	checkError("Error funSendEmail, couldn't get credentials", err)
	return nil

}
