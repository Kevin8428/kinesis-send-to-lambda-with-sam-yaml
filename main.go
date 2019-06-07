package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type configProvider struct {
}

func (c *configProvider) ClientConfig(serviceName string, cfgs ...*aws.Config) client.Config {
	return client.Config{
		Config:        &aws.Config{},
		Handlers:      request.Handlers{},
		Endpoint:      "",
		SigningRegion: "us-east-1",
		SigningName:   "",
	}
}

func main() {
	b := []byte(`{"name":"Kevin","age":31}`)
	// err := PutRecord(b)
	// fmt.Println("err emitting to firehose: ", err)
	err := PutRecordKinesis(b)
	fmt.Println("err emitting to kinesis stream: ", err)
}

// PutRecord - Put a record in the kinesis stream by a certain partition key
func PutRecord(data []byte) error {
	streamName := "kevin-kinesis-to-lambda-stack-KinesisStream-53T50NQ5NTW3"
	client := firehose.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	recordInput := firehose.PutRecordInput{
		DeliveryStreamName: &streamName,
		Record:             &firehose.Record{Data: data},
	}
	time.Sleep(5 * time.Second)
	_, err := client.PutRecord(&recordInput)
	return err
}

func PutRecordKinesis(b []byte) error {
	streamName := "kevin-kinesis-to-lambda-stack-KinesisStream-53T50NQ5NTW3"
	mySession := configProvider{}
	cfgs := aws.Config{}
	client := kinesis.New(&mySession, &cfgs)
	pk := "1"

	pri := kinesis.PutRecordInput{
		Data:         b,
		PartitionKey: &pk,
		StreamName:   &streamName,
	}
	client.PutRecord(&pri)
	return nil
}
