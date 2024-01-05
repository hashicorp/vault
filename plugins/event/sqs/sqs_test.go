package sqs

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/vault/sdk/event"
)

func getTestClient(t *testing.T) *sqs.Client {
	awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		t.Fatal(err)
	}
	return sqs.NewFromConfig(awsConfig)
}

func createQueue(t *testing.T, client *sqs.Client, queueName string) string {
	resp, err := client.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	if err != nil {
		t.Fatal(err)
	}
	return *resp.QueueUrl
}

func deleteQueue(t *testing.T, client *sqs.Client, queueURL string) {
	_, err := client.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
		QueueUrl: &queueURL,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func receiveMessage(t *testing.T, client *sqs.Client, queueURL string) string {
	resp, err := client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:        &queueURL,
		WaitTimeSeconds: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, resp.Messages, 1)
	msg := resp.Messages[0]
	_, err = client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		t.Fatal(err)
	}
	return *msg.Body
}

// TestSQS_SendOneMessage tests that the plugin basic flow of subscribe/sendevent/unsubscribe will send a message to SQS.
func TestSQS_SendOneMessage(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		t.Skip("Must set AWS_REGION")
	}
	sqsClient := getTestClient(t)
	temp, err := uuid.GenerateUUID()
	assert.Nil(t, err)
	tempQueueName := "event-sqs-test-" + temp
	tempQueueURL := createQueue(t, sqsClient, tempQueueName)
	t.Cleanup(func() {
		deleteQueue(t, sqsClient, tempQueueURL)
	})

	backend := New()
	ctx := context.Background()
	err = backend.Initialize(ctx)
	assert.Nil(t, err)

	subID, err := uuid.GenerateUUID()
	assert.Nil(t, err)

	err = backend.Subscribe(ctx, &event.SubscribeRequest{
		SubscriptionID: subID,
		Config: map[string]interface{}{
			"queue_name": tempQueueName,
			"region":     os.Getenv("AWS_REGION"),
		},
		VerifyConnection: false,
	})
	assert.Nil(t, err)

	err = backend.SendSubscriptionEvent(subID, "{}")
	assert.Nil(t, err)

	msg := receiveMessage(t, sqsClient, tempQueueURL)
	assert.Equal(t, "{}", msg)

	err = backend.Unsubscribe(context.Background(), subID)
	assert.Nil(t, err)
}
