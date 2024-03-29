package service

import (
	"github.com/6156-DonaldDuck/articles/pkg/db"
	"github.com/6156-DonaldDuck/articles/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"fmt"
	"strconv"
)

var tableName string = "Articles"

// TODO: pagination
func ListAllArticlesDynamo() ([]model.DArticle, error) {
	out, err := db.DynamoDBConn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		panic(err)
	}
	recs := []model.DArticle{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &recs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
	}
	return recs, err
}

// Filter Expression
// Also can filter by title
func GetArticleByAuthorIdDynamo(authorId uint) ([]string, error) {
	// Construct the filter builder with a name and value.
	filt := expression.Name("AuthorId").Equal(expression.Value(authorId))

	// Create the names list projection of names to project.
	proj := expression.NamesList(
		expression.Name("Title"),
	)

	// Using the filter and projections create a DynamoDB expression.
	expr, err := expression.NewBuilder().
		WithFilter(filt).
		WithProjection(proj).
		Build()
	if err != nil {
		fmt.Println(err)
	}

	// Use the built expression to populate the DynamoDB Scan API input parameters.
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	out, err := db.DynamoDBConn.Scan(input)
	fmt.Println(out)
	type titles struct {
		Title string
	}
	recs := []titles{}
	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &recs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
	}
	var res []string
	for _, entry := range recs {
		res = append(res, entry.Title)
	}
	return res, err
}

// Conditional Expression
func CreateArticleDynamo(article model.DArticle) error {
	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("Content"),
			expression.Value(article.Content),
		).Set(
			expression.Name("SectionId"),
			expression.Value(article.SectionId),
		),
	).WithCondition(
		expression.And(
			expression.AttributeNotExists(
				expression.Name("AuthorId"),
			),
			expression.AttributeNotExists(
				expression.Name("Title"),
			),
		),
	).Build()
	if err != nil {
		return err
	}
	_, err = db.DynamoDBConn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"AuthorId": {
				N: aws.String(strconv.Itoa(int(article.AuthorId))),
			},
			"Title": {
				S: aws.String(article.Title),
			},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	})
	return err
}

func UpdateArticleDynamo(updateInfo model.DArticle) error {
	update := expression.Set(
		expression.Name("Content"),
		expression.Value(updateInfo.Content),
	).Set(
		expression.Name("SectionId"),
		expression.Value(updateInfo.SectionId),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		return err
	}
	result, err := db.DynamoDBConn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"AuthorId": {
				N: aws.String(strconv.Itoa(int(updateInfo.AuthorId))),
			},
			"Title": {
				S: aws.String(updateInfo.Title),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: expr.Update(),
	})
	fmt.Println(result)
	return err
}

func DeleteArticleDynamo(article model.DArticle) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"AuthorId": {
				N: aws.String(strconv.Itoa(int(article.AuthorId))),
			},
			"Title": {
				S: aws.String(article.Title),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := db.DynamoDBConn.DeleteItem(input)
	return err
}

func TestTransactionConflict() error {
    var (
        transactItems1 []*dynamodb.TransactWriteItem
        transactItems2 []*dynamodb.TransactWriteItem
    )

    items1 := []model.DArticle{
        {AuthorId: 1, Title: "transaction1-1"},
        {AuthorId: 2, Title: "transaction1-2"},
        {AuthorId: 3, Title: "transaction1-3"},
    }

    for _, item := range items1 {
        itemav, _ := dynamodbattribute.MarshalMap(item)
        transactItems1 = append(transactItems1, &dynamodb.TransactWriteItem{
            Put: &dynamodb.Put{
                TableName: aws.String(tableName),
                Item:      itemav,
            },
        })
    }

    _, err := db.DynamoDBConn.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
        TransactItems: transactItems1,
    }) 
	if err != nil {
		return err
	}

    items2 := []model.DArticle{
        {AuthorId: 1, Title: "transaction2-1"},
        {AuthorId: 2, Title: "transaction1-2"},
        {AuthorId: 3, Title: "transaction1-3"},
    }

    for _, item := range items2 {
        itemav, _ := dynamodbattribute.MarshalMap(item)
        transactItems2 = append(transactItems2, &dynamodb.TransactWriteItem{
            Put: &dynamodb.Put{
                TableName:           aws.String(tableName),
                Item:                itemav,
                ConditionExpression: aws.String("attribute_not_exists(AuthorId) and attribute_not_exists(Title)"),
            },
        })
    }

    _, err = db.DynamoDBConn.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
        TransactItems: transactItems2,
    })
	return err
}