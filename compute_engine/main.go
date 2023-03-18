package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

const PORT = "8080"
const ENTRY_POINT = "http" //arbitrary string but shall match `Entry point` specified in the console
const BUCKET_NAME = "test-bucket-001-a"

func main() {

	var server = &http.Server{
		Addr:         fmt.Sprintf(":%v", PORT),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	http.HandleFunc("/test", handler)

	fmt.Printf("Started listening the port %v...\n", PORT)
	server.ListenAndServe()

}

//Same as `json.MarshalIndent()` but does NOT implicitly escape HTML entities such as `<`, `>` and `&`.
//ref: |https://stackoverflow.com/a/28596225/8776746|
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

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(400)
		response := Response{
			Status:  "error",
			Content: "failed to parse json",
		}
		res, _ := toJSON(response, 4)
		w.Write(res)
		return
	}
	if req.Content == "" {
		w.WriteHeader(400)
		response := Response{
			Status:  "error",
			Content: "`content` field shall not be empty",
		}
		res, _ := toJSON(response, 4)
		w.Write(res)
		return
	}

	filename := fmt.Sprintf("%v.txt", time.Now().UnixMicro())
	url, err := uploadFile(BUCKET_NAME, filename, req.Content)
	if err != nil {
		w.WriteHeader(500)
		response := Response{
			Status:  "error",
			Content: err.Error(),
		}
		res, _ := toJSON(response, 4)
		w.Write(res)
		return
	}
	fmt.Printf("url: %v\n", url)

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

func uploadFile(bucketName, filename, content string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create a client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	bucket := client.Bucket(bucketName)
	o := bucket.Object(filename)
	ww := o.NewWriter(ctx)
	if _, err = fmt.Fprintf(ww, content); err != nil {
		return "", fmt.Errorf("failed to write to a writer: %v", err)
	}
	if err := ww.Close(); err != nil {
		return "", fmt.Errorf("failed to close a writer: %v", err)
	}

	opts := storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Second * 30),
	}
	url, err := bucket.SignedURL(filename, &opts)
	if err != nil {
		return "", fmt.Errorf("failed to create a signed URL: %v", err)
	}
	return url, nil
}
