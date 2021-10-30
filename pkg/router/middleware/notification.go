package middleware

import (
	"fmt"
	"github.com/6156-DonaldDuck/articles/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
)



func publishSNS(subject string, message string) (*request.Request, *sns.PublishOutput, error){

	// Initial credentials loaded from SDK's default credential chain. Such as
	sess := session.Must(session.NewSession())

	// Create a SNS client with additional configuration
	creds := credentials.NewStaticCredentials(
		config.Configuration.AWS.PublicKey,
		config.Configuration.AWS.SecretKey,
		"",
	)
	svc := sns.New(sess, &aws.Config{Credentials: creds, Region:  aws.String(config.Configuration.SNS.Region)})

	// Send sns request
	params := sns.PublishInput{
		Message: &message,
		Subject: &subject,
		TargetArn: &config.Configuration.SNS.TopicArn,
	}
	req, resp := svc.PublishRequest(&params)
	err := req.Send()
	return req, resp, err
}




// New article notification middleware
func Notification() gin.HandlerFunc {
	return func (c *gin.Context) {
		requestMap := map[string] string {"/api/v1/articles": "POST"}
		method, ok := requestMap[c.Request.RequestURI]
		if ok && method == c.Request.Method {
			subject := "New Article Created!"
			message := "New Article Created!"
			req, resp, err := publishSNS(subject, message)
			if err == nil {
				fmt.Println(resp)
			} else {
				fmt.Print(req)
			}
		}
	}
}