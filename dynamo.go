package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type BlogPost struct {
	PostTitle string `json:"PostTitle"`
	PageTitle string `json:"PageTitle"`
	Date      string `json:"Date"`
	Body      []struct {
		Heading     string   `json:"Heading"`
		HeadingText []string `json:"HeadingText"`
	} `json:"Body"`
	ID string `json:"ID"`
}

func get_post(post_id string) (post BlogPost) {
	aws_session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println("holy crap the aws session failed")
	}
	svc := dynamodb.New(aws_session)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("BlogPosts"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(post_id),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	blog_post := BlogPost{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &blog_post)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal record, %v", err))
	}

	if post.ID == "" {
		fmt.Printf("Could not find: blogpost: %v", post_id)
	}
	return blog_post
}
