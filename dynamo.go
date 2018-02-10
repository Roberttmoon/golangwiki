package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type WikiPage struct {
	Title string
	Post  string
}

func dynamodb_svc() (*dynamodb.DynamoDB) {
	aws_session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		fmt.Printf("ERROR: aws err: %v", err)
		os.Exit(1)
	}
	svc := dynamodb.New(aws_session)
	return svc
}

func get_page (title string) (WikiPage, error) {
	svc := dynamodb_svc()
	
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("WikiPages"),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {S: aws.String(title),
			}}})

	page := WikiPage{
		Title: title,
	}

	if err != nil {
		fmt.Printf("WARN could not find entry: %v\n", err)
		return page, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &page)
	if err != nil {
		panic(fmt.Sprintf("ERROR: failed to unmarshal record: %v\n", err))
	}

	if page.Title == "" {
		fmt.Printf("WARN: could not find wiki page: %v\n", title)
		return page, err
	}

	return page, err
}

func put_page (title string, post string) (error) {
	svc := dynamodb_svc()

	page := WikiPage{
		Title: title,
		Post:  post,
	}
	wpage, err := dynamodbattribute.MarshalMap(page)

	if err != nil {
		fmt.Printf("WARN could not post wiki page: %v\n", title)
		return err
	}
	
	input := &dynamodb.PutItemInput{
		Item:      wpage,
		TableName: aws.String("WikiPages"),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Printf("Got error calling PutItem: %v\n", err)
		return err
	}
	return err
}
