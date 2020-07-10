package main

import (
	"fmt"

	"github.com/cucumber/godog"
)

func createVSphereDatastoreDiskIfNotExist(diskSize, disk, datacenter, datastore string) error {
	err := iCanRunCmdStr("", fmt.Sprintf("govc datastore.disk.info -dc %s -ds %s %s", datacenter, datastore, disk))
	if err != nil {
		return iRunCmdStr("", fmt.Sprintf("govc datastore.disk.create -dc %s -ds %s -size %s %s", datacenter, datastore, diskSize, disk))
	}
	return nil
}

func VSphereContext(s *godog.Suite) {
	s.Step(`^I create "([^"]*)" "([^"]*)" in vSphere "([^"]*)" and "([^"]*)"$`, createVSphereDatastoreDiskIfNotExist)
}
