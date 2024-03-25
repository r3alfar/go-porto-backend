package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const envLocalURL = "http://localhost:8000"

func (app *application) connectToDynamoDB() (*dynamodb.Client, error) {
	if envLocalURL == "" {
		//Load AWS Configuration (~/.aws/config)
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
		if err != nil {
			log.Fatalf("Unable to load SDK config, %v", err)
		}

		//prepare amazon DynamoDB
		client := dynamodb.NewFromConfig(cfg)
		log.Println("Successfully connect to: DynamoDB API")

		return client, nil
	} else {
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
		log.Println("Successfully connect to: DynamoDB LOCAL")

		return client, nil
	}

}
