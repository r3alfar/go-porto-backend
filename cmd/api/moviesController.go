package main

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	rd, _ := time.Parse("2006-01-02", "1986-03-07")

	highlander := models.Movie{
		ID:          "1",
		Title:       "Highlander",
		ReleaseDate: int(rd.Unix()),
		MPAARating:  "R",
		RunTime:     116,
		Description: "A very nice movie",
		CreatedAt:   int(time.Now().Unix()),
		UpdatedAt:   int(time.Now().Unix()),
	}

	movies = append(movies, highlander)

	rd, _ = time.Parse("2006-01-02", "1981-06-12")

	rotla := models.Movie{
		ID:          "2",
		Title:       "Raiders of the Lost Ark",
		ReleaseDate: int(rd.Unix()),
		MPAARating:  "PG-13",
		RunTime:     115,
		Description: "Another very nice movie",
		CreatedAt:   int(time.Now().Unix()),
		UpdatedAt:   int(time.Now().Unix()),
	}

	movies = append(movies, rotla)

	out, err := json.Marshal(movies)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) LocalDBCreateItem(w http.ResponseWriter, r *http.Request) {
	envLocalURL := "http://localhost:8000"
	fmt.Printf("Begin create at %v \n", envLocalURL)

	//Load AWS Configuration (~/.aws/config)
	// cfg, err := config.LoadDefaultConfig(context.Background(),
	// 	config.WithEndpointResolverWithOptions(
	// 		aws.EndpointResolverWithOptionsFunc(
	// 			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	// 				return aws.Endpoint{
	// 						URL: envLocalURL,
	// 					},
	// 					nil
	// 			},
	// 		),
	// 	))
	// if err != nil {
	// 	log.Fatalf("Unable to load SDK config, %v", err)
	// }
	// //integrate dynamoidb
	// client := dynamodb.NewFromConfig(cfg)

	// load aws config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load sdk config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	//construct movies
	movies := constructDummyMovies()

	writeRequest := make([]types.WriteRequest, len(movies))
	for i, movie := range movies {
		av, err := attributevalue.MarshalMap(movie)
		if err != nil {
			log.Fatalf("failed to marshal item, %v", err)
		}
		putRequest := types.PutRequest{
			Item: av,
		}
		writeRequest[i] = types.WriteRequest{
			PutRequest: &putRequest,
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"movies": writeRequest,
		},
	}

	//perform batch write opertaion
	_, err = client.BatchWriteItem(context.Background(), input)
	if err != nil {
		log.Fatalf("failed to batch write items, %v", err)
	}

	//json output after remapping to ddb json
	// out, err := json.Marshal(writeRequest)
	// if err != nil {
	// 	fmt.Println("Error marshalling items to json: ", err)
	// }

	prettyJson, err := json.MarshalIndent(movies, "", "\t")
	if err != nil {
		fmt.Println("JSON Parse Error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(prettyJson)
	// log.Println("creating:", movies)
	log.Println("Finished Create Item DynamoDB ")
}

func constructDummyMovies() []models.Movie {
	var dum []models.Movie
	rd, _ := time.Parse("2006-01-02", "1986-03-07")

	highlander := models.Movie{
		ID:          "1",
		Title:       "Highlander",
		ReleaseDate: int(rd.Unix()),
		MPAARating:  "R",
		RunTime:     116,
		Description: "Indeed a very nice movie",
		CreatedAt:   int(time.Now().Unix()),
		UpdatedAt:   int(time.Now().Unix()),
	}

	dum = append(dum, highlander)

	rd, _ = time.Parse("2006-01-02", "1981-06-12")

	rotla := models.Movie{
		ID:          "2",
		Title:       "Raiders of the Lost Ark",
		ReleaseDate: int(rd.Unix()),
		MPAARating:  "PG-13",
		RunTime:     115,
		Description: "Some Another very nice movie",
		CreatedAt:   int(time.Now().Unix()),
		UpdatedAt:   int(time.Now().Unix()),
	}

	dum = append(dum, rotla)

	return dum
}

func (app *application) DynamoDbCreateItemDummy(w http.ResponseWriter, r *http.Request) {
	log.Println("Initialize Create Item DynamoDB ")
	movies := constructDummyMovies()

	//Load AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	//prepare amazon DynamoDB
	client := dynamodb.NewFromConfig(cfg)

	prettyJson, err := json.MarshalIndent(movies, "", "\t")
	if err != nil {
		fmt.Println("JSON Parse Error", err)
		return
	}

	fmt.Println(string(prettyJson))

	writeRequest := make([]types.WriteRequest, len(movies))
	for i, movie := range movies {
		av, err := attributevalue.MarshalMap(movie)
		if err != nil {
			log.Fatalf("failed to marshal item, %v", err)
		}
		putRequest := types.PutRequest{
			Item: av,
		}
		writeRequest[i] = types.WriteRequest{
			PutRequest: &putRequest,
		}
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"movies": writeRequest,
		},
	}

	//perform batch write opertaion
	_, err = client.BatchWriteItem(context.Background(), input)
	if err != nil {
		log.Fatalf("failed to batch write items, %v", err)
	}

	//json output after remapping to ddb json
	out, err := json.Marshal(writeRequest)
	if err != nil {
		fmt.Println("Error marshalling items to json: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	log.Println("Finished Create Item DynamoDB ")
}

func (app *application) LocalGetAllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Begin LocalGetAllMovies")

	movies, err := app.DB.GetAllMovies()
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
	log.Println("Finished Get all local movies DynamoDB ")
}
