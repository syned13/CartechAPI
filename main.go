package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/CartechAPI/auth"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

const PORT = "5000"

func init() {
	gotenv.Load()
}

func main() {

	connectionString := os.Getenv("DB_CONNECTION")
	_, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("could_not_open_db: ", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", auth.Index())
	fmt.Println("Listening on port", PORT)
	http.ListenAndServe(":"+PORT, router)

}
