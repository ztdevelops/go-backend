package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/ztdevelops/go-project/src/lib/custom"
	"github.com/ztdevelops/go-project/src/lib/middleware"
	"google.golang.org/api/option"
)

// HandleRoutes initialises the connections to all the explicitly coded routes.
func (a *App) HandleRoutes() {
	app, err := middleware.InitFirebase()
	if err != nil {
		log.Fatal("failed to init firebase:", err)
	}

	a.App = *app
	a.Router.HandleFunc("/", DefaultHandler).Methods((custom.RGET))
	a.Router.HandleFunc("/test", TestAPIHandler).Methods((custom.RGET))

	// APIs (No need for auth)
	a.Router.HandleFunc("/api/signup", NotImplemented).Methods((custom.RPOST))
	a.Router.HandleFunc("/api/signin", SignInHandler).Methods((custom.RPOST))

	// APIs (Require auth)
	a.Router.HandleFunc("/api/test", a.TestVerifyToken).Methods(custom.RGET)
	a.Router.HandleFunc("/api/upload", UploadHandler).Methods(custom.RPOST)
	http.Handle("/", a.Router)
}

func transform(w http.ResponseWriter, r *http.Request) (custom.CustomWriter, custom.CustomRequest) {
	writer := custom.CustomWriter{w}
	request := custom.CustomRequest{r}

	return writer, request
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	writer.SetContentType(custom.ContentTypeJSON)

	received := custom.User{}
	if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
		log.Println("failed to decode user:", err)
		return
	}

	// Default to always true.
	received.ReturnSecureToken = true
	u, err := json.Marshal(received)
	if err != nil {
		log.Println("error marshalling received user:", err)
		return
	}

	resp, err := middleware.LoginWithFirebase(u)
	if err != nil {
		log.Println("error querying api:", err)
		return
	}

	log.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return
	}

	var result custom.UserReponse
	if err = json.Unmarshal(body, &result); err != nil {
		log.Println("error unmarshalling result:", err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	uploadType := request.GetURIParam("type")
	writer.SetContentType(custom.ContentTypeJSON)

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

		storageClient, err := storage.NewClient(request.Context(), option.WithCredentialsFile("storage.json"))
		if err != nil {
			log.Println("failed to init google cloud client:", err)
			return
		}

		sw := storageClient.Bucket("amp-bucket").Object(handler.Filename).NewWriter(request.Context())

		if _, err = io.Copy(sw, file); err != nil {
			log.Println("error copying file to cloud storage:", err)
			return
		}

		if err = sw.Close(); err != nil {
			log.Println("failed to close connection to cloud storage:", err)
			return
		}

		log.Println("file successfully uploaded.")
	}
}

func (a *App) TestVerifyToken(w http.ResponseWriter, r *http.Request) {
	if err := middleware.VerifyToken(&a.App, r); err != nil {
		log.Println("token failed:", err)
		return
	}
	log.Println("token ok.")
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom.ENDPOINT_HIT, "default")
	fmt.Fprintf(w, "Default landing page.")
}

func TestAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(custom.ENDPOINT_HIT, "test")
	users := []custom.User{}
	json.NewEncoder(w).Encode(users)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented.")
}
