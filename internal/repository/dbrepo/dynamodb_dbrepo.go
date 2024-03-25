package dbrepo

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const envLocalURL = "http://localhost:8000"

type DynamoDBRepo struct {
	DB *dynamodb.Client
}

// GetMovie implements repository.DatabaseRepository.
func (d DynamoDBRepo) GetMovie() (*models.Movie, error) {
	result, err := d.DB.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("movies"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Fatalln("failed to call GetMovies")
	}

	if result.Item == nil {
		fmt.Println("Item Not Found")
	}

	log.Println("success get movie:")
	log.Print(result)
	// for k, item := range result.Item {
	// 	fmt.Printf("%s: %v\n", k, item)
	// }

	item := models.Movie{}
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatalf("failed to unmarshal record: %v", err)
	}

	log.Println("success unmarshal result:")
	log.Print(item)

	return &item, nil
}

// GetMovies implements repository.DatabaseRepository.
func (d DynamoDBRepo) GetAllMovies() ([]*models.Movie, error) {
	tabName := "movies"
	result, err := d.DB.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &tabName,
		Limit:     aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("Failed to fetch AllMovies\n")
	}

	var items []*models.Movie
	for _, item := range result.Items {
		var data models.Movie
		if err := attributevalue.UnmarshalMap(item, &data); err != nil {
			fmt.Println("failed to unmarshal item: ", err)
			continue
		}
		items = append(items, &data)
	}

	fmt.Println("Items: ", items)

	return items, nil
}

// PutMovie implements repository.DatabaseRepository.
func (d DynamoDBRepo) PutMovie() (*models.Movie, error) {
	panic("unimplemented")
}

// PutMovies implements repository.DatabaseRepository.
func (d DynamoDBRepo) PutMovies() ([]*models.Movie, error) {
	panic("unimplemented")
}

// GetMovies implements repository.DatabaseRepository.
func (d DynamoDBRepo) GetMovies() ([]*models.Movie, error) {
	panic("unimplemented")
}

func ConnectDynamoDB() *dynamodb.Client {
	if envLocalURL != "" {
		//Load aws config
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
			log.Fatalf("Unable to load SDK config: %v", err)
		}

		client := dynamodb.NewFromConfig(cfg)

		return client
	} else {
		//Load AWS Configuration (~/.aws/config)
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
		if err != nil {
			log.Fatalf("Unable to load SDK config, %v", err)
		}

		//prepare amazon DynamoDB
		client := dynamodb.NewFromConfig(cfg)

		return client
	}
}
