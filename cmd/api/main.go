package main

import (
	_ "backend/cmd/env"
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = 8080

type application struct {
	DataSource string
	Domain     string
	DB         repository.DatabaseRepository
}

func main() {
	//set application config
	var app application

	//access env
	envPort := os.Getenv("PORT")
	if envPort == "" {
		log.Fatal("PORT is not set")
	}

	//connect to the database
	client, err := app.connectToDynamoDB()
	if err != nil {
		log.Fatalf("Unable to connect to Dynamodb Initialization")
	}
	app.DB = dbrepo.DynamoDBRepo{DB: client}

	log.Println("Starting Application on port: ", envPort)
	log.Println("ENV port: ", envPort)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%v", envPort), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
