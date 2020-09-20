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

package routes

import (
	"encoding/json"
	"io/ioutil"
)

type routerule struct {
  Name     string `json:"name,omitempty"`
  Backend  string `json:"backend,omitempty"`
  BasePath string `json:"basePath,omitempty"`
}

var routerules = []routerule{}
var defaultRouterule = routerule{
	Backend: "httpbin.org",
	BasePath: "/",
}

func ReadRoutesFile() error {
	routeListBytes, err := ioutil.ReadFile("routes.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(routeListBytes, &routerules); err != nil {
		return err
	}

	for _, routerule := range routerules {
		if routerule.Name == "default" {
			defaultRouterule.BasePath = routerule.BasePath
			defaultRouterule.Backend = routerule.Backend
		}
	}

	return nil
}

func ListRouteRules() []routerule {
  return routerules
}

func GetDefaultRouteRule() (string, string) {
	return defaultRouterule.Backend, defaultRouterule.BasePath
}

func GetRouteRule(name string) (string, string) {
  for _, routerule := range routerules {
		if routerule.Name == name {
			return routerule.Backend, routerule.BasePath
		}
	}
  return GetDefaultRouteRule()
}
