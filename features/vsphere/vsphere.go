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

package vsphere

import (
	"fmt"
	"os"
	"text/template"

	local "github.com/c3y1huang/caasp-e2e/features/local"
)

const (
	//CloudConf file name
	CloudConf = "vsphere.conf"
	//CloudConfTemplate use to create CloudConf
	CloudConfTemplate = `[Global]
insecure-flag = "1"
[VirtualCenter "{{.VCenter}}"]
datacenters = "{{.DataCenter}}"
user = "{{.User}}"
password = "{{.Password}}"
[Workspace]
server = "{{.VCenter}}"
datacenter = "{{.DataCenter}}"
default-datastore = "{{.DataStore}}"
resourcepool-path = "{{.ResourcePool}}"
folder = "{{.ClusterFolder}}"
[Disk]
scsicontrollertype = pvscsi
[Labels]
region = "{{.Region}}"
zone = "{{.Zone}}"
`
    // EnvClusterVCenter URL
	EnvClusterVCenter = "GODOG_VSPHERE_CLUSTER_VCENTER"
	//EnvClusterPrefix used to create cluster uniqueness
	EnvClusterPrefix = "GODOG_VSPHERE_CLUSTER_PREFIX"
	//EnvClusterFolder name in vCenter where virtual machine should be in
	EnvClusterFolder = "GODOG_VSPHERE_CLUSTER_FOLDER"
	//EnvClusterDataCenter name in vCenter where virtual machine should be in
	EnvClusterDataCenter = "GODOG_VSPHERE_CLUSTER_DATACENTER"
	// EnvClusterDataStore name in vCenter to use by cluster
	EnvClusterDataStore = "GODOG_VSPHERE_CLUSTER_DATASTORE"
	// EnvClusterResourcePool path in vCenter
	EnvClusterResourcePool = "GODOG_VSPHERE_CLUSTER_RESOURCEPOOL"
	// EnvClusterRegion name in vCenter
	EnvClusterRegion = "GODOG_VSPHERE_CLUSTER_REGION"
	// EnvClusterZone name in vCenter
	EnvClusterZone = "GODOG_VSPHERE_CLUSTER_ZONE"
	// EnvUser name use to login vCenter
	EnvUser = "GODOG_VSPHERE_USER"
	// EnvPassword use to login vCenter
	EnvPassword = "GODOG_VSPHERE_PASSWORD"
)

//VSphere object
type VSphere struct {
	Config Config
	Local  local.Local // local object
}

//Config is the configuration object
type Config struct {
	VCenter       string // vCenter URL
	ClusterFolder string // vCenter virtual machines cluster folder name
	User          string // vCenter login username
	Password      string // vCenter login password
	DataCenter    string // vCenter datacenter name
	DataStore     string // vCenter datastore name
	ResourcePool  string // vCenter resource pool path
	Region        string // vCenter region
	Zone          string // vCenter zone
}

//NewVSphere returns VSphere object
func NewVSphere() (*VSphere, error) {
	var v VSphere
	return &v, nil
}

func (v *VSphere) createCloudConf() error {
	var ok bool
	if v.Config.VCenter, ok = os.LookupEnv(EnvClusterVCenter); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterVCenter)
	}
	if v.Config.ClusterFolder, ok = os.LookupEnv(EnvClusterFolder); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterFolder)
	}
	if v.Config.User, ok = os.LookupEnv(EnvUser); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvUser)
	}
	if v.Config.Password, ok = os.LookupEnv(EnvPassword); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvPassword)
	}
	if v.Config.DataCenter, ok = os.LookupEnv(EnvClusterDataCenter); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterDataCenter)
	}
	if v.Config.DataStore, ok = os.LookupEnv(EnvClusterDataStore); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterDataStore)
	}
	if v.Config.ResourcePool, ok = os.LookupEnv(EnvClusterResourcePool); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterResourcePool)
	}
	if v.Config.Region, ok = os.LookupEnv(EnvClusterRegion); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterRegion)
	}
	if v.Config.Zone, ok = os.LookupEnv(EnvClusterZone); !ok {
		return fmt.Errorf(local.MissingEnvVariable, EnvClusterZone)
	}

	//create file from template
	var t *template.Template
	fWriter, err := os.Create(CloudConf)
	if err != nil {
		return err
	}
	t = template.New(CloudConf)
	t, _ = t.Parse(string(CloudConfTemplate))
	err = t.Execute(fWriter, v.Config)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}
	return nil
}

func (v *VSphere) createDatastoreDiskIfNotExist(diskSize, disk, datacenter, datastore string) error {
	err := v.Local.RunCmdSimple("", fmt.Sprintf("govc datastore.disk.info -dc %s -ds %s %s", datacenter, datastore, disk))
	if err != nil {
		return v.Local.RunCmdSimple("", fmt.Sprintf("govc datastore.disk.create -dc %s -ds %s -size %s %s", datacenter, datastore, diskSize, disk))
	}
	return nil
}
