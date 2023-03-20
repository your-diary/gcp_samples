package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"compute_engine/cloud_sql"
	"compute_engine/cloud_storage"
	"compute_engine/config"
	"compute_engine/firestore"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

const configFile = "./config.json"

func main() {

	config, err := config.New(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, config)
	})

	fmt.Printf("Started listening the port %v...\n", config.Port)
	server.ListenAndServe()

}

// Same as `json.MarshalIndent()` but does NOT implicitly escape HTML entities such as `<`, `>` and `&`.
// ref: |https://stackoverflow.com/a/28596225/8776746|
func toJSON(t any, numIndent int) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	if numIndent != 0 {
		encoder.SetIndent("", strings.Repeat(" ", numIndent))
	}
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func writeErrorResponse(w http.ResponseWriter, status int, reason string) {
	w.WriteHeader(500)
	response := Response{
		Status:  "error",
		Content: reason,
	}
	res, _ := toJSON(response, 4)
	w.Write(res)

}

func handler(w http.ResponseWriter, r *http.Request, config *config.Config) {
	w.Header().Add("Content-Type", "application/json")
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, 400, "failed to parse json")
		return
	}
	if req.Content == "" {
		writeErrorResponse(w, 400, "`content` field shall not be empty")
		return
	}

	content :=
		req.Content

	cloudStorage, err := cloud_storage.New(config.CloudStorage)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	filename := fmt.Sprintf("%v.txt", time.Now().UnixMicro())
	err = cloudStorage.UploadFile(filename, content)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	url, err := cloudStorage.CreateSignedURL(filename)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	fmt.Printf("url: %v\n", url)

	dbSQL, err := cloud_sql.New(config.Postgres)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	err = dbSQL.Insert(content)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}

	dbFirestore, err := firestore.New(config.Firestore)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}
	err = dbFirestore.Insert(content)
	if err != nil {
		writeErrorResponse(w, 500, err.Error())
		return
	}

	w.WriteHeader(200)
	response := Response{
		Status:  "success",
		Content: url,
	}
	res, _ := toJSON(response, 4)
	fmt.Printf("res: %v\n", res)
	fmt.Printf("res: %v\n", string(res))
	w.Write(res)
	return
}
