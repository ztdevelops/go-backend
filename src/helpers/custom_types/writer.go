package custom_types

import (
	"encoding/json"
)

// Respond responds to the API requests with a status code and message.
func (w *CustomWriter) Respond(code int, message interface{}) {
	r := Response{
		Status:  code,
		Message: message,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(r)
}

// SetContentType sets the content type of the response.
func (w *CustomWriter) SetContentType(ct string) {
	w.Header().Set("Content-Type", ct)
}
