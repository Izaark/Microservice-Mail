package controllers

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/tasks/Microservice-Mail/models"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const ctimeLayout = "2006-01-02T15:04:05.000Z"

var strFromEmail string
var strPasswordEmail string
var envLogType string
var envLogFile string

//checkError: func for generate error and logs
func checkError(strMessage string, err error) {
	if err != nil {
		err = errors.New(strMessage + "->" + err.Error())
		log.Println(err)
	}
}

func MailRouter() {
	var osLogFile *os.File
	var err error
	//type deploy Mode
	if os.Getenv("ENV_DEPLOY_MODE") == "DEBUG" {
		gin.SetMode(gin.DebugMode)

	} else if os.Getenv("ENV_DEPLOY_MODE") == "RELEASE" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("ERROR MailRouter: invalid environment mode")
		return
	}
	//logFile
	envLogType = os.Getenv("ENV_LOG_TYPE")
	envLogFile = os.Getenv("ENV_LOG_FILE")
	if envLogType == "FILE" || envLogType == "MIXED" {
		osLogFile, err = os.Create(envLogFile + "." + time.Now().UTC().Format(ctimeLayout))
		defer osLogFile.Close()
		checkError("ERROR MailRouter: couldn't create log file", err)
	}
	switch envLogType {
	case "FILE":
		gin.DefaultWriter = io.MultiWriter(osLogFile)
	case "MIXED":
		gin.DefaultWriter = io.MultiWriter(osLogFile, os.Stdout)
	case "CONSOLE":
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	}
	log.SetOutput(gin.DefaultWriter)

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
		mail.POST("/users/email/file", handlerSendEmail)
	}

	router.Run(":" + os.Getenv("API_PORT"))
}

//handlerSendEmail: send the mail from frontend whit a specific struct
func handlerSendEmail(c *gin.Context) {
	var vEmail models.ObjUserInfo
	var vginResponse gin.H
	var err error

	err = c.BindJSON(&vEmail)
	if err != nil {
		err = errors.New("ERROR handlerSendEmail: couldn't bind payload provided ObjUserInfo struct -> " + err.Error())
		vginResponse = gin.H{"message": "error reading payload provided", "response": nil, "error": "Response Error", "status": http.StatusBadRequest}
		c.JSON(http.StatusBadRequest, vginResponse)
		return
	}

	//send the email depending by strKey
	strURL := strings.TrimPrefix(c.Request.RequestURI, "/api/users/email/")
	strKey := strings.Split(strURL, "/")[0]

	switch strKey {
	case "text":
		err = funSendEmail(vEmail.Email, vEmail.ObjTemplate.Body)
		if err != nil {
			err = errors.New("ERROR handlerSendEmail: couldn't init funSendEmail -> " + err.Error())
			vginResponse = gin.H{"message": "error funSendEmail", "response": nil, "error": "Response Error", "status": http.StatusBadRequest}
			c.JSON(http.StatusBadRequest, vginResponse)
			return
		}
	case "file":
		err = funSendEmailAttach(vEmail.Email, vEmail.ObjTemplate.Body)
		if err != nil {
			err = errors.New("ERROR handlerSendEmail: couldn't init funSendEmail -> " + err.Error())
			vginResponse = gin.H{"message": "error funSendEmailAttach", "response": nil, "error": "Response Error", "status": http.StatusBadRequest}
			c.JSON(http.StatusBadRequest, vginResponse)
			return
		}
	}

	vginResponse = gin.H{"message": "Message send", "response": vEmail.Email, "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}

//funSendEmail: generate struct email whit text
func funSendEmail(strToEmail, strBody string) error {
	var buff bytes.Buffer
	var vTemplate models.ObjTemplate

	fromAddress := mail.Address{"Administración:", os.Getenv("FROM_EMAIL")}
	strFromEmail = os.Getenv("FROM_EMAIL")
	strPasswordEmail = os.Getenv("FROM_EMAIL_PASSWORD")

	//parse template, for matching vTemplate struct
	templ, err := template.ParseFiles("./templates/html/email_template.html")
	checkError("Error funSendEmail: error parsing html template", err)

	vTemplate.Date = time.Now()
	vTemplate.Body = strBody
	err = templ.Execute(&buff, vTemplate)
	checkError("funSendEmail: couldn't Execute template", err)

	message := "From: " + fromAddress.String() + "\n" + "To: " + strToEmail + "\n" + "Subject: Bienvenid@\n" + "Content-type: text/html" + buff.String() + "\n"

	err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", strFromEmail, strPasswordEmail, "smtp.gmail.com"), strFromEmail, []string{strToEmail}, []byte(message))
	checkError("Error funSendEmail: error smtp to get credentials", err)
	return nil
}

//funSendEmailAttach: Send email whit files, in this case is a image !
func funSendEmailAttach(strToEmail, strBody string) error {
	fromAddress := mail.Address{"Soporte Técnico", os.Getenv("FROM_EMAIL")}
	strFromEmail = os.Getenv("FROM_EMAIL")
	strPasswordEmail = os.Getenv("FROM_EMAIL_PASSWORD")

	message := gomail.NewMessage()
	message.SetHeader("From", fromAddress.String())
	message.SetHeader("To", strToEmail)
	message.SetHeader("Subject", "llego la ayuda")
	message.SetBody("text/plain", strBody)
	message.Attach("./templates/images/RandM.png")

	dial := gomail.NewDialer("smtp.gmail.com", 587, strFromEmail, strPasswordEmail)

	err := dial.DialAndSend(message)
	checkError("ERROR funSendEmailAttach: couldn't send email", err)
	return nil
}
