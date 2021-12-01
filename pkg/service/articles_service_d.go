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

// Update Expression
// func CreateArticleDynamo(article model.DArticle) (model.DArticle, error){
// 	// Create an update to set two fields in the table.
// 	update := expression.Set(
// 		expression.Name("Content"),
// 		expression.Value(article.Content),
// 	).Set(
// 		expression.Name("SectionId"),
// 		expression.Value(article.SectionId),
// 	)

// 	// Create the DynamoDB expression from the Update.
// 	expr, err := expression.NewBuilder().
// 		WithUpdate(update).
// 		Build()

// 	// Use the built expression to populate the DynamoDB UpdateItem API
// 	// input parameters.
// 	input := &dynamodb.UpdateItemInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"AuthorId": {
// 				// N: aws.Int(strconv.Itoa(int(article.AuthorId))),
// 				N: aws.Int(1),
// 			},
// 			"Title": {
// 				S: aws.String(article.Title),
// 			},
// 		},
// 		ReturnValues:     aws.String("ALL_NEW"),
// 		TableName:        aws.String(tableName),
// 		UpdateExpression: expr.Update(),
// 	}

// 	result, err := db.DynamoDBConn.UpdateItem(input)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	res := model.DArticle{}
// 	err = dynamodbattribute.UnmarshalMap(result.Attributes , &res)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	return res, err
// }

func UpdateArticleDynamo(updateInfo model.Article) error {

}

// func DeleteArticleByIdDynamo(articleId uint) error {

// }