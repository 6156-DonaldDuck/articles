package db

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/6156-DonaldDuck/articles/pkg/config"
    "fmt"
    "log"
)

var DynamoDBConn *dynamodb.DynamoDB

type DynamoArticle struct {
	Title string
	AuthorId uint
	Content string
	SectionId uint
}

func init() {
	tableName := "Articles"
	createTable(tableName)
	insertTestItem(tableName)
}

func createTable(tableName string) {
	sess := session.Must(session.NewSession())
	creds := credentials.NewStaticCredentials(
		config.Configuration.AWS.PublicKey,
		config.Configuration.AWS.SecretKey,
		"",
	)
	DynamoDBConn = dynamodb.New(sess, &aws.Config{Credentials: creds, Region:  aws.String(config.Configuration.SNS.Region)})

	listTableInput := &dynamodb.ListTablesInput{}
	result, err := DynamoDBConn.ListTables(listTableInput)
    if err != nil {
        fmt.Println(err.Error())
    }

    for _, n := range result.TableNames {
		if tableName == *n {
			fmt.Println("Table already exists", tableName)	
			return
		}
    }
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("AuthorId"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},		
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("AuthorId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("RANGE"),
			},		
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err = DynamoDBConn.CreateTable(createTableInput)
	if err != nil {
		log.Fatalf("Got error calling CreateTable: %s", err)
	}
	fmt.Println("Created the table", tableName)	
}

func insertTestItem(tableName string) {
	testArticle := DynamoArticle{
		AuthorId: 1,
		Title: "test article title",
		Content: "test article content",
		SectionId: 1,
	}

	item, err := dynamodbattribute.MarshalMap(testArticle)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return
	}
	fmt.Println(item)
	updateInput := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(tableName),
	}
	_, err = DynamoDBConn.PutItem(updateInput)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		return
	}
	fmt.Println("Successfully put item into table")
}