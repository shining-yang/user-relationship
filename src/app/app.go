package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "pgv3"
    "router"
)


func main() {
    fmt.Println("create database schemas ....")
    db := pgv3.ConnectDatabase()
    err := pgv3.CreateSchema(db)
    if err != nil {
        panic(err)
    }

    fmt.Println("setup routes ...")
    r := mux.NewRouter()
    r.HandleFunc("/users", router.RetrieveUsers).Methods("GET")
    r.HandleFunc("/users", router.CreateUser).Methods("POST")
    r.HandleFunc("/users/{user_id:[0-9]+}/relationships", router.RetrieveRelationship).Methods("GET")
    r.HandleFunc("/users/{user_id:[0-9]+}/relationships/{other_user_id:[0-9]+}", router.UpdateRelationship).Methods("PUT")
	r.NotFoundHandler = http.HandlerFunc(router.HandleBlackHoleRoute)

    fmt.Println("begin to serve ...")
    http.ListenAndServe(":8000", r)
}
