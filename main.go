package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/order"
	"github.com/CartechAPI/service"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/subosito/gotenv"
)

var (
	// create a 1 request/second limiter and
	// every token bucket in it will expire 1 hour after it was initially set.
	defaultLimiter = tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	loginLimiter = tollbooth.NewLimiter(0.4, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
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

	configureLimiters()

	router := mux.NewRouter()
	defineRoutes(db, channel, router)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("Listening on port", port)
	// log.Println()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), loggedRouter))
}

func configureLimiters() {
	defaultLimiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	loginLimiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
}

func defineRoutes(db *sql.DB, channel *amqp.Channel, router *mux.Router) {
	router.HandleFunc("/", auth.Index()).Methods(http.MethodGet)

	router.Handle("/login", tollbooth.LimitHandler(loginLimiter, auth.Login(db))).Methods(http.MethodPost)
	router.HandleFunc("/signup", auth.SignUp(db)).Methods(http.MethodPost)
	router.Handle("/session", auth.StoreSession(db)).Methods(http.MethodPost)

	router.Handle("/mechanic/signup", tollbooth.LimitHandler(defaultLimiter, auth.MechanichSignUp(db))).Methods(http.MethodPost)
	router.HandleFunc("/mechanic/login", auth.MechanicLogin(db)).Methods(http.MethodPost)

	router.HandleFunc("/service", service.GetAllServices(db)).Methods(http.MethodGet)
	router.HandleFunc("/service/category", service.GetAllServiceCategories(db)).Methods(http.MethodGet)
	router.HandleFunc("/service/category/{category_id}", service.GetServicesByCategoryID(db)).Methods(http.MethodGet)

	router.HandleFunc("/order", order.CreateServiceOrder(db, channel)).Methods(http.MethodPost)
	router.HandleFunc("/order", order.GetAllServiceOrders(db)).Methods(http.MethodGet)
	router.Handle("/order/past", order.GetAllPastServiceOrders(db)).Methods(http.MethodGet)
	router.Handle("/order/current", order.GetAllCurrentOrders(db)).Methods(http.MethodGet)
	router.HandleFunc("/order/{order_id}", order.UpdateServiceOrder(db)).Methods(http.MethodPatch)
	router.HandleFunc("/order/{order_id}", order.GetServiceOrder(db)).Methods(http.MethodGet)
	router.HandleFunc("/order/{order_id}/mechanic", order.AssignMechanicToOrder(db)).Methods(http.MethodPut)
}
