package service

import (
	"context"
	"log"

	"github.com/mailersend/mailersend-go"
)

func SendEmail(apiKey, fromEmail, toEmail, subject, message string) error {

	client := mailersend.NewMailersend(apiKey)

	from := mailersend.From{
		Name:  "Your Name",
		Email: "EMAIL",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  "Recipient Name",
			Email: toEmail,
		},
	}

	email := client.Email.NewMessage()
	email.SetFrom(from)
	email.SetRecipients(recipients)
	email.SetSubject(subject)
	email.SetHTML(message)
	email.SetText(message)

	ctx := context.TODO()
	_, err := client.Email.Send(ctx, email)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
