package main

import (
	"log"

	"github.com/tasks/Microservice-Mail/config"
	"github.com/tasks/Microservice-Mail/controllers"
)

func init() {
	err := config.FunInitEnvironment()
	if err != nil {
		log.Fatal("*ERROR init: couldn't initialize configuration -> ", err.Error())
	}
}

func main() {
	controllers.MailRouter()
}
