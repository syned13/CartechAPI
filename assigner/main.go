package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/CartechAPI/notifications"
	"github.com/CartechAPI/order"
	"github.com/CartechAPI/user"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/subosito/gotenv"
)

const ASSIGNER_QUEUE = "assign-order"

func init() {
	gotenv.Load()
}

func main() {
	connectionString := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("could_not_open_db: ", err)
	}

	defer db.Close()

	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	failOnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		ASSIGNER_QUEUE, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Println("Received a message:" + string(d.Body))

			err := assignOrder(db, d.Body)
			if err != nil {
				log.Println("error_handling_message: " + err.Error())
			}

			// time.Sleep(2 * time.Second)
			log.Println("Done")
		}
	}()

	// Receive from channel and assign to value
	value := <-forever
	fmt.Println(value)
}

func assignOrder(db *sql.DB, messageBody []byte) error {
	serviceOrder := order.ServiceOrder{}
	err := json.Unmarshal(messageBody, &serviceOrder)
	if err != nil {
		return err
	}

	usr, err := user.GetUserByID(db, serviceOrder.UserID)
	if err != nil {
		return err
	}

	log.Print("user: ")
	log.Println(usr)

	log.Print("service_order: ")
	log.Println(serviceOrder)

	return notifications.SendNotificationToMechanics("¡Nueva orden!", "Hay una nueva orden disponible")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
