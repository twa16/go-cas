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
package main

import (
	"github.com/twa16/go-cas/client"
	"fmt"
	"os"
	"time"
)

func main() {
	config := gocas.CASServerConfig{}
	config.ServerHostname = "https://login.gmu.edu"
	config.IgnoreSSLErrors = false

	config.StartLocalAuthenticationProcess()
	for true {
		if config.Response != nil {
			fmt.Printf("Logged in as %s\n", config.Response.Username)
			break
		}
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)
}
