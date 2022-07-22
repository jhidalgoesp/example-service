package awsProvider

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Provider struct {
	session *session.Session
}

func NewProvider() Provider {
	return Provider{
		session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))}
}

func (aws Provider) NewDynamoClient() *dynamodb.DynamoDB {
	return dynamodb.New(aws.session)
}

func (aws Provider) NewSSMClient() *ssm.SSM {
	return ssm.New(aws.session)
}
