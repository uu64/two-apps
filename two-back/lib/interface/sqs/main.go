package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// ReceiveMessage receives a message from the queue
func ReceiveMessage(svc *sqs.SQS, queueName string) (*sqs.ReceiveMessageOutput, error) {
	var item *sqs.ReceiveMessageOutput

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return item, err
	}

	item, err = svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            urlResult.QueueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(60),
	})

	return item, err
}

// SendMessage sends a message to the queue
func SendMessage(svc *sqs.SQS, queueName string, message string) error {
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  &message,
		QueueUrl:     urlResult.QueueUrl,
	})

	return err
}

// DeleteMessage deletes the message from the queue
func DeleteMessage(svc *sqs.SQS, queueName string, handle string) error {
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return err
	}

	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      urlResult.QueueUrl,
		ReceiptHandle: &handle,
	})

	return err
}
