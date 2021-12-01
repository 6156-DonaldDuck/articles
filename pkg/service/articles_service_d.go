package service

import (
	"github.com/6156-DonaldDuck/articles/pkg/db"
	"github.com/6156-DonaldDuck/articles/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"fmt"
	// "strconv"
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

	// Using the filter and projections create a DynamoDB expression from the two.
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
func CreateArticleDynamo(article model.DArticle) {
	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("AuthorId"),
			expression.Value(article.authorId),
		),
	).WithCondition(
		expression.And(
			expression.AttributeNotExists(
				expression.Name("AuthorId"),
			),
			expression.Equal(
				expression.Name("Title"),
				expression.Value(article.title),
			),
		),
	).Build()
	if err != nil {
		panic(err)
	}
	ut, err := db.DynamoDBConn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"AuthorId": &types.AttributeValueMemberS{Value: article.authorId},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
}

func UpdateArticleDynamo(updateInfo model.Article) error {
	update := expression.Set(
		expression.Name("Title"),
		expression.Value(updateInfo.Title),
	).Set(
		expression.Name("Content"),
		expression.Value(updateInfo.Content),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()

	if err != nil {
		panic(err)
	}
	result, err := db.DynamoDBConn.UpdateItem(&dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key: map[string]*dynamodb.AttributeValue{
			"AuthorId": {
				S: aws.String(string(updateInfo.AuthorId)),
			},
			"Title": {
				S: aws.String(updateInfo.Title),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: expr.Update(),
	})
	if err != nil {
		panic(err)
	} else {
		return err
	}

	fmt.Println(result)
	return nil
}

func DeleteArticleByIdDynamo(articleId uint) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ArticleId": {
				S: aws.String(string(articleId)),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := db.DynamoDBConn.DeleteItem(input)

	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	return nil
}
