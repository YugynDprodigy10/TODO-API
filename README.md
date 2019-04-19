# TODO-API

This is a REST Api built using golang and MongoDB. The user can Perform basic CRUD operations.

## Requirements

* [Go](https://golang.org/doc/install)
* [Mongodb](https://docs.mongodb.com/manual/installation/)
## How to Run:

* [Download](https://github.com/sagarchoudhary96/TODO-API/archive/master.zip) or [Clone](https://github.com/sagarchoudhary96/TODO-API.git) the repository.
* Install required dependencies:
```
go get -u github.com/gorilla/mux`

go get gopkg.in/mgo.v2
```
* Open `terminal` and start mongodb using command `mongod`.
* cd into the cloned repository and run `go run main.go` to start api server.
* Api Endpoints can be accessed on `http://localhost:3000`

## API Routes {Method} :

1. **/todo  {GET}** :
  Fetch all Todos List.
2. **/todo  {POST/PUT}** : Add todo to the list. 
  __params:__ `description -> String`
3. **/todo/{id}  {GET}** : Fetch todo by ID.
4. **/todo/{id}  {PATCH}** : Mark specific todo as done.
5. **/todo/{id}  {DELETE}** : Deelete todo by ID.
