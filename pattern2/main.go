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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"net/http"
)

const address = "0.0.0.0:7080"

//ErrorMessage hold the return value when there is an error
type ErrorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

var errorMessage = ErrorMessage{StatusCode: http.StatusInternalServerError}

type HttpHandler struct{}

func main() {

	var wait time.Duration

	fmt.Println("Starting server - ", address)

	handler := HttpHandler{}

	//the following code is from gorilla mux samples
	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)

	fmt.Println("Shutting down")

	os.Exit(0)
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req *http.Request
	downstreamURL := "https://httpbin.org/" + r.URL.Path //"http://127.0.0.1:9080" + r.URL.Path
	client := &http.Client{}

	fmt.Println(">> Request received")

	//TODO: add request logic here.

	req, err := http.NewRequest(r.Method, downstreamURL, nil)
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

	//TODO: add response logic here.

	fmt.Println(">> Response sent")
	ResponseHandler(w, resp.Header, resp.StatusCode, respBody)
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
func ResponseHandler(w http.ResponseWriter, headers http.Header, statusCode int, response []byte) {

	for headerName, headerValue := range headers {
		w.Header().Set(headerName, strings.Join(headerValue, ","))
	}
	w.WriteHeader(statusCode)
	_, err := w.Write(response)
	if err != nil {
		fmt.Println(err)
	}
}
