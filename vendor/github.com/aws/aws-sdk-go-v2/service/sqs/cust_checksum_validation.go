package sqs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
)

// addValidateSendMessageChecksum adds the ValidateMessageChecksum middleware
// to the stack configured for the SendMessage Operation.
func addValidateSendMessageChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateSendMessageChecksum)
}

// validateSendMessageChecksum validates the SendMessage operation's input
// message payload MD5 checksum matches that returned by the API.
//
// The input and output types must match the SendMessage operation.
func validateSendMessageChecksum(input, output interface{}) error {
	in, ok := input.(*SendMessageInput)
	if !ok {
		return fmt.Errorf("wrong input type, expect %T, got %T", in, input)
	}
	out, ok := output.(*SendMessageOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	// Nothing to validate if the members aren't populated.
	if in.MessageBody == nil || out.MD5OfMessageBody == nil {
		return nil
	}

	if err := validateMessageChecksum(*in.MessageBody, *out.MD5OfMessageBody); err != nil {
		return messageChecksumError{
			MessageID: aws.ToString(out.MessageId),
			Err:       err,
		}
	}
	return nil
}

// addValidateSendMessageBatchChecksum adds the ValidateMessagechecksum
// middleware to the stack configured for the SendMessageBatch operation.
func addValidateSendMessageBatchChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateSendMessageBatchChecksum)
}

// validateSendMessageBatchChecksum validates the SendMessageBatch operation's
// input messages body MD5 checksum matches those returned by the API.
//
// The input and output types must match the SendMessageBatch operation.
func validateSendMessageBatchChecksum(input, output interface{}) error {
	in, ok := input.(*SendMessageBatchInput)
	if !ok {
		return fmt.Errorf("wrong input type, expect %T, got %T", in, input)
	}
	out, ok := output.(*SendMessageBatchOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	outEntries := map[string]sqstypes.SendMessageBatchResultEntry{}
	for _, e := range out.Successful {
		outEntries[*e.Id] = e
	}

	var failedMessageErrs []messageChecksumError
	for _, inEntry := range in.Entries {
		outEntry, ok := outEntries[*inEntry.Id]
		// Nothing to validate if the members aren't populated.
		if !ok || inEntry.MessageBody == nil || outEntry.MD5OfMessageBody == nil {
			continue
		}

		if err := validateMessageChecksum(*inEntry.MessageBody, *outEntry.MD5OfMessageBody); err != nil {
			failedMessageErrs = append(failedMessageErrs, messageChecksumError{
				MessageID: aws.ToString(outEntry.MessageId),
				Err:       err,
			})
		}
	}

	if len(failedMessageErrs) != 0 {
		return batchMessageChecksumError{
			Errs: failedMessageErrs,
		}
	}

	return nil
}

// addValidateReceiveMessageChecksum adds the ValidateMessagechecksum
// middleware to the stack configured for the ReceiveMessage operation.
func addValidateReceiveMessageChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateReceiveMessageChecksum)
}

// validateReceiveMessageChecksum validates the ReceiveMessage operation's
// input messages body MD5 checksum matches those returned by the API.
//
// The input and output types must match the ReceiveMessage operation.
func validateReceiveMessageChecksum(_, output interface{}) error {
	out, ok := output.(*ReceiveMessageOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	var failedMessageErrs []messageChecksumError
	for _, msg := range out.Messages {
		// Nothing to validate if the members aren't populated.
		if msg.Body == nil || msg.MD5OfBody == nil {
			continue
		}

		if err := validateMessageChecksum(*msg.Body, *msg.MD5OfBody); err != nil {
			failedMessageErrs = append(failedMessageErrs, messageChecksumError{
				MessageID: aws.ToString(msg.MessageId),
				Err:       err,
			})
		}
	}

	if len(failedMessageErrs) != 0 {
		return batchMessageChecksumError{
			Errs: failedMessageErrs,
		}
	}

	return nil
}

// messageChecksumValidator provides the function signature for the operation's
// validator.
type messageChecksumValidator func(input, output interface{}) error

// addValidateMessageChecksum adds the ValidateMessageChecksum middleware to
// the stack with the passed in validator specified.
func addValidateMessageChecksum(stack *middleware.Stack, o Options, validate messageChecksumValidator) error {
	if o.DisableMessageChecksumValidation {
		return nil
	}

	m := validateMessageChecksumMiddleware{
		validate: validate,
	}
	err := stack.Initialize.Add(m, middleware.Before)
	if err != nil {
		return fmt.Errorf("failed to add %s middleware, %w", m.ID(), err)
	}

	return nil
}

// validateMessageChecksumMiddleware provides the Initialize middleware for
// validating an operation's message checksum is validate. Needs to b
// configured with the operation's validator.
type validateMessageChecksumMiddleware struct {
	validate messageChecksumValidator
}

// ID returns the Middleware ID.
func (validateMessageChecksumMiddleware) ID() string { return "SQSValidateMessageChecksum" }

// HandleInitialize implements the InitializeMiddleware interface providing a
// middleware that will validate an operation's message checksum based on
// calling the validate member.
func (m validateMessageChecksumMiddleware) HandleInitialize(
	ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, meta middleware.Metadata, err error,
) {
	out, meta, err = next.HandleInitialize(ctx, input)
	if err != nil {
		return out, meta, err
	}

	err = m.validate(input.Parameters, out.Result)
	if err != nil {
		return out, meta, fmt.Errorf("message checksum validation failed, %w", err)
	}

	return out, meta, nil
}

// validateMessageChecksum compares the MD5 checksums of value parameter with
// the expected MD5 value. Returns an error if the computed checksum does not
// match the expected value.
func validateMessageChecksum(value, expect string) error {
	msum := md5.Sum([]byte(value))
	sum := hex.EncodeToString(msum[:])
	if sum != expect {
		return fmt.Errorf("expected MD5 checksum %s, got %s", expect, sum)
	}

	return nil
}

// messageChecksumError provides an error type for invalid message checksums.
type messageChecksumError struct {
	MessageID string
	Err       error
}

func (e messageChecksumError) Error() string {
	prefix := "message"
	if e.MessageID != "" {
		prefix += " " + e.MessageID
	}
	return fmt.Sprintf("%s has invalid checksum, %v", prefix, e.Err.Error())
}

// batchMessageChecksumError provides an error type for a collection of invalid
// message checksum errors.
type batchMessageChecksumError struct {
	Errs []messageChecksumError
}

func (e batchMessageChecksumError) Error() string {
	var w strings.Builder
	fmt.Fprintf(&w, "message checksum errors")

	for _, err := range e.Errs {
		fmt.Fprintf(&w, "\n\t%s", err.Error())
	}

	return w.String()
}
