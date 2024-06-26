package main

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Wrold! from %s", app.Domain)
}

func (app *application) dummyJson(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Farrel backend",
		Version: "1.0.0alpha",
	}

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	//settingup API RESPONSE HEADER
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

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
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
							URL: envLocalURL,
						},
						nil
				},
			),
		))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	//construct movies
	movies := constructDummyMovies()

	//integrate dynamoidb
	client := dynamodb.NewFromConfig(cfg)

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

func (app *application) LocalPutMovie(w http.ResponseWriter, r *http.Request) {

	//VALIDATOR

	//Content Type validator
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// var request map[string]interface{}
	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		fmt.Println("FAILED TO DECODE")
		return
	}
	fmt.Println("DECODED: ", movie)

	movie.UpdatedAt = int(time.Now().Unix())
	movie.CreatedAt = int(time.Now().Unix())

	// //create additional fields
	// extra := map[string]interface{}{
	// 	"created_at": time.Now().Unix(),
	// 	"updated_at": time.Now().Unix(),
	// }
	//merge request data with extra fields
	// for key, value := range extra {
	// 	movie[key] = value
	// }

	//convert to free string and to JSON
	prettyJson, err := json.Marshal(movie)
	if err != nil {
		fmt.Println("JSON Parse Error", err)
		return
	}

	fmt.Println("request to JSON: ", string(prettyJson))

	fmt.Println("Received JSON data as struct: ", movie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received POST Request"))

}
