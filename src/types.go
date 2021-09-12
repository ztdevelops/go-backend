package main

import (
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/ztdevelops/go-project/src/helpers/database"
)

type Router struct {
	*mux.Router
}

type App struct {
	Router
	database.Database
	firebase.App
}

type Writer struct {
	http.ResponseWriter
}

type Response struct {
	Status 	int
	Message interface{}
}

const PORT string = "8000"

var SharedApp App