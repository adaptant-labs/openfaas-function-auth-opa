// Copyright (c) Adaptant Solutions AG 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package function

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/adaptant-labs/go-function-sdk"
)

var (
	opaServer string
)

type opaResponse struct {
	Result map[string]bool `json:"result"`
}

func checkAuthz(values map[string]map[string]interface{}) (opaResponse, error) {
	jsonValue, err := json.Marshal(values)
	if err != nil {
		log.Println("Failed marshalling OPA payload", err)
		return opaResponse{}, err
	}

	opaUri := opaServer + "/v1/data/openfaas/authz"

	resp, err := http.Post(opaUri, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("Failed delivering OPA payload", err)
		return opaResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed reading response body", err)
		return opaResponse{}, err
	}

	var response opaResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Failed unmarshalling response body", err)
		return opaResponse{}, err
	}

	return response, nil
}

func isAuthorized(req handler.Request) bool {
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		return false
	}

	authzRequestPayload := map[string]map[string]interface{}{
		"input": {
			"function": host,
			"port":     port,
			"method":   req.Method,
			"query":    req.QueryString,
			"user":     req.Header.Get("Authorization"),
		},
	}

	// Check for a decision from OPA - if none is available, default-deny
	response, err := checkAuthz(authzRequestPayload)
	if err != nil {
		return false
	}

	// Otherwise default-allow
	allow := true

	// Iterate over decision results looking for any deny cases
	for k, v := range response.Result {
		if k == "deny" && v == true {
			allow = false
		} else if k == "allow" && v == false {
			allow = false
		}
	}

	return allow
}

func init() {
	opaServer = os.Getenv("OPA_URL")
	if opaServer == "" {
		opaServer = "http://localhost:8181"
	}
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error
	res := handler.Response{}

	if !isAuthorized(req) {
		message := "Unauthorized."
		res.Body = []byte(message)
		res.StatusCode = http.StatusUnauthorized
		return res, err
	}

	res.StatusCode = http.StatusOK
	res.Body = []byte("Authorization OK.")

	return res, err
}
