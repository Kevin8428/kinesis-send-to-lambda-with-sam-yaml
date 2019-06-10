package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func main() {
	b := []byte(`{"name":"Kevin","age":31}`)
	err := PutRecordKinesis(b)
	fmt.Println("err emitting to kinesis stream: ", err)
}

func PutRecordKinesis(b []byte) error {
	streamName := "kevin-kinesis-to-lambda-stack-KinesisStream-53T50NQ5NTW3"
	client := kinesis.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	pk := "1"

	pri := kinesis.PutRecordInput{
		Data:         b,
		PartitionKey: &pk,
		StreamName:   &streamName,
	}
	client.PutRecord(&pri)
	return nil
}
