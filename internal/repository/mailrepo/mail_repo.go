package mailrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type MailRepo struct {
	SESClient   *ses.Client
	FromAddress string
}

func (r *MailRepo) SendWelcomeEmail(to string, firstname string, lastname string) error {
	subject := "Welcome to the Share2Teach platform"
	body := fmt.Sprintf("Hello %s %s,\n\nWelcome to the Share2Teach platform. We are excited to have you on board!", firstname, lastname)

	// Create the SES input
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(r.FromAddress),
	}

	// Send the email using SES
	_, err := r.SESClient.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Welcome email sent successfully!")
	return nil
}

func (r *MailRepo) SendPasswordResetRequest(to, resetKey string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf("Use the following token to reset your password: %s", resetKey)

	// Create the SES input
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(r.FromAddress),
	}

	// Send the email using SES
	_, err := r.SESClient.SendEmail(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Println("Password reset email sent successfully!")
	return nil
}
