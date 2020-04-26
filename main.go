package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/order"
	"github.com/CartechAPI/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/subosito/gotenv"
)

var port string

func init() {
	gotenv.Load()
	port = os.Getenv("PORT")
}

func main() {

	connectionString := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("could_not_open_db: ", err)
	}
	defer db.Close()

	queueConnection, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	if err != nil {
		log.Fatal("could_not_connect_to_queue: ", err)
	}
	defer queueConnection.Close()

	channel, err := queueConnection.Channel()
	if err != nil {
		log.Fatal("could_not_open_channel_to_queue: ", err)
	}
	defer channel.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", auth.Index()).Methods(http.MethodGet)
	router.HandleFunc("/login", auth.Login(db)).Methods(http.MethodPost)
	router.HandleFunc("/signup", auth.SignUp(db)).Methods(http.MethodPost)

	router.HandleFunc("/mechanic/signup", auth.MechanichSignUp(db)).Methods(http.MethodPost)
	router.HandleFunc("/mechanic/login", auth.MechanicLogin(db)).Methods(http.MethodPost)

	router.HandleFunc("/service", service.GetAllServices(db)).Methods(http.MethodGet)
	router.HandleFunc("/service/category", service.GetAllServiceCategories(db)).Methods(http.MethodGet)
	router.HandleFunc("/service/category/{category_id}", service.GetServicesByCategoryID(db)).Methods(http.MethodGet)

	router.HandleFunc("/order", order.CreateServiceOrder(db, channel)).Methods(http.MethodPost)
	router.HandleFunc("/order", order.GetAllServiceOrders(db)).Methods(http.MethodGet)
	router.HandleFunc("/order/{order_id}", order.UpdateServiceOrder(db)).Methods(http.MethodPatch)
	router.HandleFunc("/order/{order_id}", order.GetServiceOrder(db)).Methods(http.MethodGet)

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
