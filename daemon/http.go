package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

func httpServer() {
	loadTemplates()

	router := mux.NewRouter()
	router.HandleFunc("/", DashboardHandler)
	router.HandleFunc("/systems/create", CreateSystemHandler)

	log.Printf("Server Started...")
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func DashboardHandler(response http.ResponseWriter, request *http.Request) {
	conn := pool.Get()
	defer conn.Close()
	unixTime, err := redis.Int64(conn.Do("GET", "ap-srv1"))
	if err != nil {
		// TODO: Create an error template to handle these correctly, for now...
		log.Printf("Error retrieving value from redis: %s", err)
	}
	timestamp := time.Unix(unixTime, 0).UTC()
	details := struct {
		Time string
	}{
		timestamp.Format(time.RFC3339),
		// fmt.Sprintf("%s", timestamp),
	}
	err = templates.ExecuteTemplate(response, "dashboard.html", details)
	if err != nil {
		log.Printf("Template Error %s", err.Error())
		http.Error(response, "Template Error", 500)
	}
}

func CreateSystemHandler(response http.ResponseWriter, request *http.Request) {
	conn := pool.Get()
	defer conn.Close()
	timestamp, err := conn.Do("GET", "systemList")
	if err != nil {
		// TODO: Create an error template to handle these correctly, for now...
		log.Printf("Error retrieving value from redis: %s", err)
	}
	details := struct {
		Time string
	}{
		fmt.Sprintf("%s", timestamp),
	}
	t, _ := template.ParseFiles("templates/dashboard.html")
	t.Execute(response, details)
}

func loadTemplates() {
	var err error
	t := template.New("gonocmon").Funcs(template.FuncMap{})
	templates, err = t.ParseGlob("templates/*.html")
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
}
