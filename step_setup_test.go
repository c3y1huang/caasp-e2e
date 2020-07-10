/*
Copyright (c) 2020 SUSE LLC.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/tidwall/gjson"
)

var Clusters []byte

func getClusterInfoFrom(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}
	Clusters = f
	return nil
}

func getClusterAccessFrom(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	return os.Setenv("KUBECONFIG", absPath)
}

func iSearchInClusterInfo(find string) error {
	value := gjson.Get(string(Clusters), find)
	CheckString = value.String()
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^cluster info from "([^"]*)"$`, getClusterInfoFrom)
	s.Step(`^cluster access from "([^"]*)"$`, getClusterAccessFrom)
	s.Step(`^I search "([^"]*)" in cluster info$`, iSearchInClusterInfo)
}
