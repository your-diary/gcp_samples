package abc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type Request struct {
	Content string `json:"content"`
}

type Response struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

func init() {
	entryPoint := "http" //arbitrary string but shall match `Entry point` specified in the console
	functions.HTTP(entryPoint, handler)
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
		res, _ := json.MarshalIndent(response, "", "    ")
		w.Write([]byte(res))
		return
	}
	if req.Content == "" {
		w.WriteHeader(400)
		response := Response{
			Status:  "error",
			Content: "`content` field shall not be empty",
		}
		res, _ := json.MarshalIndent(response, "", "    ")
		w.Write([]byte(res))
		return
	}
	w.WriteHeader(200)
	response := Response{
		Status:  "success",
		Content: fmt.Sprintf("Hello, %v!", req.Content),
	}
	res, _ := json.MarshalIndent(response, "", "    ")
	w.Write([]byte(res))
	return
}
