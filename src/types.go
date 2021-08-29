package main

import (
	"github.com/gorilla/mux"
	"github.com/ztdevelops/go-project/src/helpers/database"
)

type Router struct {
	*mux.Router
}

type App struct {
	Router
	database.Database
}