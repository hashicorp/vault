package physical

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestDynamoDBBackend(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.SkipNow()
	}

	creds, err := credentials.NewEnvCredentials().Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// If the variable is empty or doesn't exist, the default
	// AWS endpoints will be used
	endpoint := os.Getenv("AWS_DYNAMODB_ENDPOINT")

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	conn := dynamodb.New(session.New(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}))

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("dynamodb", logger, map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)
}

func TestDynamoDBHABackend(t *testing.T) {
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.SkipNow()
	}

	creds, err := credentials.NewEnvCredentials().Get()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// If the variable is empty or doesn't exist, the default
	// AWS endpoints will be used
	endpoint := os.Getenv("AWS_DYNAMODB_ENDPOINT")

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	conn := dynamodb.New(session.New(&aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}))

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	table := fmt.Sprintf("vault-dynamodb-testacc-%d", randInt)

	defer func() {
		conn.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(table),
		})
	}()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("dynamodb", logger, map[string]string{
		"access_key":    creds.AccessKeyID,
		"secret_key":    creds.SecretAccessKey,
		"session_token": creds.SessionToken,
		"table":         table,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("dynamodb does not implement HABackend")
	}
	testHABackend(t, ha, ha)
}
