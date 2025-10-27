package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ReadRequestJSON reads from the body of http request,
// decodes it into a pointer variable to struct of desired type and returns the pointer
func ReadRequestJSON[T any](w http.ResponseWriter, r *http.Request) *T {
	var data T

	// read from the body of the http request sent to server into a slice of bytes
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println("!!!")
		log.Printf("Error: %v\n", err)
		return nil
	}

	// decode slice of bytes into a pointer variable of type struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error: %v\n", err)
	}

	// return a pointer to the struct containing decoded request json data
	return &data

}
