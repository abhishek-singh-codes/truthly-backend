package main

import (
	"fmt"
	"truthly/internals/app"
	"truthly/internals/util/logger"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Application started")
	log := logger.InitLogger()
	log.Info("Logger initialized")

	if err := godotenv.Load(); err != nil {
		log.Error("Error loading env file: " + err.Error())
		//panic(err)
	}

	app.Start(log)
}
