package tests

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoMock struct {
	dynamodbiface.DynamoDBAPI
	GetItemCount     *int
	UpdateItemCount  *int
	Error            error
	ReturnGetItem    *dynamodb.GetItemOutput
	ReturnUpdateItem *dynamodb.UpdateItemOutput
}

func NewDynamoMock(err error, returnGetItem *dynamodb.GetItemOutput,
	returnUpdateItem *dynamodb.UpdateItemOutput) DynamoMock {
	return DynamoMock{
		GetItemCount:     new(int),
		UpdateItemCount:  new(int),
		Error:            err,
		ReturnGetItem:    returnGetItem,
		ReturnUpdateItem: returnUpdateItem,
	}
}

func (m DynamoMock) GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	*m.GetItemCount++

	if m.Error != nil {
		return nil, m.Error
	}

	return m.ReturnGetItem, nil
}

func (m DynamoMock) UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	*m.UpdateItemCount++

	if m.Error != nil {
		return nil, m.Error
	}

	return m.ReturnUpdateItem, nil
}
