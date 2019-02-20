package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session, _ = mgo.Dial("127.0.0.1")
var c = session.DB("todoDb").C("todo")

// todo item struct
type ToDoItem struct {
	ID          bson.ObjectId `bson: "_id, omitempty"`
	Date        time.Time
	Description string
	Done        bool
}

func main() {
	//adding mongodb
	if session != nil {
		session.SetMode(mgo.Monotonic, true)
		fmt.Println("mongod running")
		defer session.Close()
	} else {
		log.Fatal("MongoDB not running")
	}

	// define port
	port := ":3000"
	router := mux.NewRouter()

	// routes handlers

	// add/update todo
	router.HandleFunc("/todo", AddToDo).Methods("POST", "PUT")

	// perform system healt check
	router.HandleFunc("/health", Health).Methods("GET")

	// start API Server
	fmt.Println("API Listening at port: ", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// healthcheck func
func Health(w http.ResponseWriter, r *http.Request) {
	// set status to OK
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content_Type", "application/json")

	// send the json Respone
	io.WriteString(w, `{"alive": true}`)
}

// add todo func
func AddToDo(w http.ResponseWriter, r *http.Request) {

	_ = c.Insert(ToDoItem{
		bson.NewObjectId(),
		time.Now(),
		r.FormValue("description"),
		false,
	})

	result := ToDoItem{}
	// lookup for the inserted data in the mongodb collection and store in result
	_ = c.Find(bson.M{"description": r.FormValue("description")}).One(&result)
	// return result
	json.NewEncoder(w).Encode(result)
}
