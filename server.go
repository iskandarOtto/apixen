package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maszuari/apixen/db"
	handlerorg "github.com/maszuari/apixen/handlers"
	modelorg "github.com/maszuari/apixen/models"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	d := db.DBConnect()
	h := handlerorg.NewHandler(modelorg.NewOrgModel(d))

	r := mux.NewRouter()
	r.HandleFunc("/orgs/{orgname}/comments/", h.GetCommentsByOrgName).Methods("GET")
	r.HandleFunc("/orgs/{orgname}/comments/", h.SaveComment).Methods("POST")
	r.HandleFunc("/orgs/{orgname}/comments", h.DeleteCommentsByOrgName).Methods("DELETE")
	r.HandleFunc("/orgs/{orgname}/members/", h.GetMembersByOrgName).Methods("GET")
	r.HandleFunc("/hello", h.Hello).Methods("GET")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}

}
