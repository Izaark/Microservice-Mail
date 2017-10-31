package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
)

func FunInitEnvironment() error {
	err := godotenv.Load("environment.env")
	if err != nil {
		err = errors.New("FunInitConfig: couldn't initialize environment " + "->" + err.Error())
		log.Println(err)
	}
	return nil
}
