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

package cluster

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/tidwall/gjson"

	"github.com/c3y1huang/caasp-e2e/features/local"
)

const (
	//EnvJSONFile path
	EnvJSONFile = "GODOG_CLUSTER_JSON_FILE"
)

//Cluster object
type Cluster struct {
	Clusters []byte
	Found    string
}

//NewCluster returns Cluster object
func NewCluster() (*Cluster, error) {
	var c Cluster
	var err error
	c.Clusters, err = local.ReadFileFromEnv(EnvJSONFile)
	if err != nil {
		return &c, err
	}
	return &c, nil
}

//search JSON data for testing purpose. Remove if not required
func (c *Cluster) search(find string) error {
	value := gjson.Get(string(c.Clusters), find)
	c.Found = value.String()
	return nil
}

//find JSON data for testing purpose. Remove if not required
func (c *Cluster) find(expectedDocString *godog.DocString) error {
	var expected = expectedDocString.Content
	if expected != c.Found {
		return fmt.Errorf("\nwant:\n%s\ngot:\n%s", expected, c.Found)
	}
	return nil
}
