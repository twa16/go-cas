/*
 * Copyright 2017 Manuel Gauto (github.com/twa16)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
*/

package gocas

import (
	"crypto/tls"
	"net/http"
	"fmt"
	"goji.io"
	"goji.io/pat"
	"io/ioutil"
	"strings"
	"io"
)

type CASServerConfig struct {
	ServerHostname   string
	IgnoreSSLErrors  bool
	httpServerCloser io.Closer
	Response         *CASValidationResponse
}

type CASValidationResponse struct {
	IsValid bool
	Username string
}

func (server *CASServerConfig) StartLocalAuthenticationProcess() {
	userAuthURL := server.ServerHostname + "/login?service=http%3A%2F%2Flocalhost%3A8883%2Fcaslogin"
	fmt.Println("Go to: "+userAuthURL)
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/caslogin"), server.handleCASCallback)
	go http.ListenAndServe(":8883", mux)
}

func (server *CASServerConfig)handleCASCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	server.Response, _ = server.ValidateTicket(r.FormValue("ticket"))
	fmt.Fprint(w, "You may now close the window.")
}

func (server *CASServerConfig) ValidateTicket(ticket string) (*CASValidationResponse, error){
	url := server.ServerHostname+"/validate?service=http%3A%2F%2Flocalhost%3A8883%2Fcaslogin&ticket="+ticket
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)
	response := CASValidationResponse{}
	if strings.HasPrefix(body, "no") {
		response.IsValid = false
		return &response, nil
	} else {
		username := strings.Split(body, "\n")[1]
		response.IsValid = true
		response.Username = username
		return &response, nil
	}
}

func (server *CASServerConfig) getHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: server.IgnoreSSLErrors},
	}
	return &http.Client{Transport: tr}
}