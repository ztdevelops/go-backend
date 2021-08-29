package routes

import "github.com/gorilla/mux"

// HTTP Request Methods
const RGET string 		= "GET"
const RPOST string 		= "POST"
const RPUT string		= "PUT"
const RDELETE string 	= "DELETE"

const ENDPOINT_HIT 		= "ENDPOINT HIT:"

type Router struct { 
	*mux.Router 
}
