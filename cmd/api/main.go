package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"fmt"
	"log"
	"net/http"
)

const port = 8800

type application struct {
	DataSource string
	Domain     string
	DB         repository.DatabaseRepository
}

func main() {
	//set application config
	var app application

	//connect to the database
	client, err := app.connectToDynamoDB()
	if err != nil {
		log.Fatalf("Unable to connect to Dynamodb Initialization")
	}
	app.DB = dbrepo.DynamoDBRepo{DB: client}

	log.Println("Starting Application on port: ", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
