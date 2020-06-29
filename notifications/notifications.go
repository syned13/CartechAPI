package notifications

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

// SendNotificationToMechanics sends a notification to a the mechanics topic
func SendNotificationToMechanics(title string, body string) error {
	ctx := context.Background()

	fmt.Println(os.Getenv("SERVICE_ACCOUNT_ID"))
	opt := option.WithCredentialsFile("../credentials/cartech-12e63-ffeab1e71964.json")
	conf := &firebase.Config{
		ServiceAccountID: os.Getenv("SERVICE_ACCOUNT_ID"),
		ProjectID:        os.Getenv("FIREBASE_PROJECT_ID"),
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}
	message := messaging.Message{
		Topic: "mechanic",
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(ctx, &message)
	if err != nil {
		return err
	}

	fmt.Println("HERE")

	fmt.Println("response_from_fmc: " + response)

	return nil
}

// SendNotificationToSingleUser sends a notification to a single user with the specified token
func SendNotificationToSingleUser(token string, title string, body string) error {
	ctx := context.Background()

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	log.Println(path)

	opt := option.WithCredentialsFile("credentials/cartech-12e63-ffeab1e71964.json")
	conf := &firebase.Config{
		ServiceAccountID: os.Getenv("SERVICE_ACCOUNT_ID"),
		ProjectID:        os.Getenv("FIREBASE_PROJECT_ID"),
	}

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return err
	}
	message := messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := client.Send(ctx, &message)
	if err != nil {
		return err
	}

	fmt.Println("HERE")

	fmt.Println("response_from_fmc: " + response)

	return nil
}
