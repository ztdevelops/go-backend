package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ztdevelops/go-project/src/helpers/custom_types"
)

// HandleRoutes initialises the connections to all the explicitly coded routes.
func (a *App) HandleRoutes() {
	app, err := InitFirebase()
	if err != nil {
		log.Fatal("failed to init firebase:", err)
	}

	a.App = *app
	a.Router.HandleFunc("/", DefaultHandler).Methods((custom_types.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom_types.RGET))

	// APIs (No need for auth)
	a.Router.HandleFunc("/api/signup", NotImplemented).Methods((custom_types.RPOST))
	a.Router.HandleFunc("/api/signin", LogInWithFirebase).Methods((custom_types.RPOST))

	// APIs (Require auth)
	a.Router.HandleFunc("/api/test", a.TestVerifyToken).Methods(custom_types.RGET)
	a.Router.HandleFunc("/api/upload", UploadHandler).Methods(custom_types.RPOST)
	http.Handle("/", a.Router)
}

func transform(w http.ResponseWriter, r *http.Request) (custom_types.CustomWriter, custom_types.CustomRequest) {
	writer := custom_types.CustomWriter{w}
	request := custom_types.CustomRequest{r}
	
	return writer, request
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	uploadType := request.GetURIParam("type")
	writer.SetContentType(custom_types.ContentTypeJSON)

	if uploadType == "file" {
		request.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("myFile")
		if err != nil {
			log.Println("error retrieving file:", err)
			return
		}
		defer file.Close()
		
		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		tempFilename := "upload-*" + filepath.Ext(handler.Filename)
		tempFile, err := ioutil.TempFile("../temp", tempFilename)
		if err != nil {
			log.Println("failed to create tempFile...", err)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("failed to read file:", err)
			return
		}
		tempFile.Write(fileBytes)

		log.Println("file successfully uploaded.")
	}
}

func (a *App) TestVerifyToken(w http.ResponseWriter, r *http.Request) {
	if err := VerifyToken(&a.App, r); err != nil {
		log.Println("token failed:", err)
		return
	}
	log.Println("token ok.")
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func TestAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom_types.ENDPOINT_HIT, "test")
	users := []custom_types.User{
	}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}
