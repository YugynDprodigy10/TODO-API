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

	// get all todo or single todo by id
	router.HandleFunc("/todo", GetToDo).Methods("GET")
	router.HandleFunc("/todo/{id}", GetToDo).Methods("GET")

	// update todo state
	router.HandleFunc("/todo/{id}", MarkDone).Methods("PATCH")

	// delete todo
	router.HandleFunc("/todo/{id}", DeleteToDo).Methods("DELETE")

	// perform system health check
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

// get todo by id
func GetByID(id string) []ToDoItem {
	var result ToDoItem
	var res []ToDoItem
	_ = c.Find(bson.M{"id": bson.ObjectIdHex(id)}).One(&result)
	res = append(res, result)
	return res
}

// get all todo item
func GetToDo(w http.ResponseWriter, r *http.Request) {
	// result array
	var resArr []ToDoItem

	// get url params
	params := mux.Vars(r)

	// get todo id
	id := params["id"]

	// get todo/todo's
	if id != "" {
		resArr = GetByID(id)
	} else {
		_ = c.Find(nil).All(&resArr)
	}

	// return result array
	json.NewEncoder(w).Encode(resArr)
}

// update todo state
func MarkDone(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// get the todo id and convert to bson id type
	id := bson.ObjectIdHex(params["id"])

	// update the db
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"done": true}})

	// if db update fails
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": `+err.Error()+`}`)
	} else {
		// update successful
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"Updated": true}`)
	}
}

// delete specific todo
func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// get the todo id and convert to bson id type
	id := bson.ObjectIdHex(params["id"])

	// revome from db
	err := c.Remove(bson.M{"id": id})

	if err == mgo.ErrNotFound {
		json.NewEncoder(w).Encode(err.Error())
	} else {
		io.WriteString(w, "{result: 'OK'}")
	}
}
