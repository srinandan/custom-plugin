// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
)

const port = "0.0.0.0:7080"

//ErrorMessage hold the return value when there is an error
type ErrorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

var errorMessage = ErrorMessage{StatusCode: http.StatusInternalServerError}

func main() {

	fmt.Println("Starting server - ", port)

	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ProcessRequest(w, r)
	}))

}

func ProcessRequest(w http.ResponseWriter, r *http.Request) {
	var req *http.Request
	downstreamURL := "http://127.0.0.1:9080" + r.URL.Path
	client := &http.Client{}

	fmt.Println(">> Request received")

	req, err := http.NewRequest("GET", downstreamURL, nil)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		ErrorHandler(w, err)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorHandler(w, err)
		return
	}
	fmt.Println("%s", string(respBody))
	fmt.Println(">> Response sent")

	var jsonMap map[string]string
	_ = json.Unmarshal(respBody, &jsonMap)

	ResponseHandler(w, jsonMap, false)

}

func ErrorHandler(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	errorMessage.Message = err.Error()

	if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
		fmt.Println(err)
	}
}

//ResponseHandler returns a 200 when the response is successful
func ResponseHandler(w http.ResponseWriter, response interface{}, text bool) {
	if !text {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println(err)
		}
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		if str, ok := response.(string); ok {
			w.Write([]byte(str))
		}
	}
}
