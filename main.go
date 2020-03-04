package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CartechAPI/auth"
	"github.com/CartechAPI/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	router := mux.NewRouter()

	router.HandleFunc("/", auth.Index()).Methods(http.MethodGet)
	router.HandleFunc("/login", auth.Login(db)).Methods(http.MethodPost)
	router.HandleFunc("/signup", auth.SignUp(db)).Methods(http.MethodPost)

	router.HandleFunc("/service", service.GetAllServices(db)).Methods(http.MethodGet)

	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
