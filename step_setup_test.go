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
