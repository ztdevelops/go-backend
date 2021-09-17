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
	a.Router.HandleFunc("/api/signup", SignUpHandler).Methods((custom.RPOST))
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

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	writer.SetContentType(custom.ContentTypeJSON)

	// 1. decode received data into User struct
	received := custom.UserForFirebase{}
	if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
		errMsg := fmt.Sprint("failed to decode user:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	// 2. marshal received data into JSON
	received.ReturnSecureToken = true
	u, err := json.Marshal(received)
	if err != nil {
		errMsg := fmt.Sprint("failed to marshal user:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	// 3. query firebase to create new user
	resp, err := middleware.SignUpWithFirebase(u)
	if err != nil {
		errMsg := fmt.Sprint("failed to sign up with firebase:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	// 4. read response from firebase
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprint("failed to read response:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	// 5. unmarshal response into a struct
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		errMsg := fmt.Sprint("error unmarshalling results:", err)
		log.Println(errMsg)
		writer.Respond(resp.StatusCode, errMsg)
		return
	}

	// 6. check for potential errors revealed in api response from firebase
	if resp.StatusCode != http.StatusOK {
		errJSON, ok := result["error"].(map[string]interface{})
		if !ok {
			log.Println("something went wrong...")
			writer.Respond(http.StatusInternalServerError, "something went wrong... Maybe firebase changed their API response?")
			return
		}
		errMsg := errJSON["message"]
		log.Println(errMsg)
		writer.Respond(resp.StatusCode, errMsg)
		return
	}

	// 7. return response
	result["status"] = http.StatusOK
	json.NewEncoder(writer).Encode(result)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	writer.SetContentType(custom.ContentTypeJSON)

	received := custom.UserForFirebase{}
	if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
		errMsg := fmt.Sprint("failed to decode user:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	// Default to always true.
	received.ReturnSecureToken = true
	u, err := json.Marshal(received)
	if err != nil {
		errMsg := fmt.Sprint("failed to marshal user:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusBadRequest, errMsg)
		return
	}

	log.Println(string(u))
	resp, err := middleware.LoginWithFirebase(u)
	if err != nil {
		errMsg := fmt.Sprint("failed to query API:", err)
		log.Println(errMsg)
		writer.Respond(http.StatusInternalServerError, errMsg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response body:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := "invalid signin credentials"
		log.Printf("Status %v: %v", resp.StatusCode, errMsg)
		writer.Respond(http.StatusUnauthorized, errMsg)
		return
	}

	var result custom.UserSignInResponse
	if err = json.Unmarshal(body, &result); err != nil {
		log.Println("error unmarshalling result:", err)
		return
	}
	result.Status = http.StatusOK
	json.NewEncoder(w).Encode(result)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	writer, request := transform(w, r)
	uploadType := request.GetURIParam("type")
	writer.SetContentType(custom.ContentTypeJSON)

	switch uploadType {
	case "file":
		handleFileUpload(writer, request)
	}
}

func handleFileUpload(w custom.CustomWriter, r custom.CustomRequest) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		errMsg := fmt.Sprint("error retrieving file from request:", err)
		log.Println(errMsg)
		w.Respond(http.StatusBadRequest, errMsg)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	opt := custom.GetOpt("CLOUD_JSON")
	storageClient, err := storage.NewClient(r.Context(), opt)
	if err != nil {
		errMsg := fmt.Sprint("failed to init google cloud client:", err)
		log.Println(errMsg)
		w.Respond(http.StatusInternalServerError, errMsg)
		return
	}

	sw := storageClient.Bucket("amp-bucket").Object(handler.Filename).NewWriter(r.Context())
	if _, err = io.Copy(sw, file); err != nil {
		errMsg := fmt.Sprint("error copying file to cloud storage:", err)
		log.Println(errMsg)
		w.Respond(http.StatusInternalServerError, errMsg)
		return
	}
	if err = sw.Close(); err != nil {
		errMsg := fmt.Sprint("failed to close connection to cloud storage:", err)
		log.Println(errMsg)
		w.Respond(http.StatusInternalServerError, errMsg)
		return
	}

	log.Println("file successfully uploaded.")
	w.Respond(http.StatusOK, "file successfully uploaded")
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
