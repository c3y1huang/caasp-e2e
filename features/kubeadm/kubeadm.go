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

package kubeadm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"

	vsphere "github.com/c3y1huang/caasp-e2e/features/vsphere"
)

const (
	//K8sConfigDir path where have all K8s cluster configurations
	K8sConfigDir = "/etc/kubernetes"
)

//Kubeadm object
type Kubeadm struct {
	Configmap            Configmap
	ClusterConfiguration ClusterConfiguration
}

//Configmap of kubeadm-config object
type Configmap struct {
	APIVersion string `yaml:"apiVersion"`
	Data       struct {
		ClusterConfiguration string `yaml:"ClusterConfiguration"`
		ClusterStatus        string `yaml:"ClusterStatus"`
	} `yaml:"data"`
	Kind     string                 `yaml:"kind"`
	Metadata map[string]interface{} `yaml:"metadata"`
}

//ClusterConfiguration object
type ClusterConfiguration struct {
	APIServer struct {
		CertSANs  []string `yaml:"certSANs"`
		ExtraArgs struct {
			CloudConfig                  string `yaml:"cloud-config"`
			CloudProvider                string `yaml:"cloud-provider"`
			AuthorizationMode            string `yaml:"authorization-mode"`
			EnableAdmissionPlugins       string `yaml:"enable-admission-plugins"`
			OidcCaFile                   string `yaml:"oidc-ca-file"`
			OidcClientID                 string `yaml:"oidc-client-id"`
			OidcGroupsClaim              string `yaml:"oidc-groups-claim"`
			OidcIssuerURL                string `yaml:"oidc-issuer-url"`
			OidcUsernameClaim            string `yaml:"oidc-username-claim"`
			ServiceAccountIssuer         string `yaml:"service-account-issuer"`
			ServiceAccountSigningKeyFile string `yaml:"service-account-signing-key-file"`
		} `yaml:"extraArgs"`
		ExtraVolumes           []ExtraVolume `yaml:"extraVolumes"`
		TimeoutForControlPlane string        `yaml:"timeoutForControlPlane"`
	} `yaml:"apiServer"`
	APIVersion           string `yaml:"apiVersion"`
	CertificatesDir      string `yaml:"certificatesDir"`
	ClusterName          string `yaml:"clusterName"`
	ControlPlaneEndpoint string `yaml:"controlPlaneEndpoint"`
	ControllerManager    struct {
		ExtraArgs struct {
			CloudConfig   string `yaml:"cloud-config"`
			CloudProvider string `yaml:"cloud-provider"`
		} `yaml:"extraArgs"`
		ExtraVolumes []ExtraVolume `yaml:"extraVolumes"`
	} `yaml:"controllerManager"`
	DNS struct {
		ImageRepository string `yaml:"imageRepository"`
		ImageTag        string `yaml:"imageTag"`
		Type            string `yaml:"type"`
	} `yaml:"dns"`
	Etcd struct {
		Local struct {
			DataDir         string `yaml:"dataDir"`
			ImageRepository string `yaml:"imageRepository"`
			ImageTag        string `yaml:"imageTag"`
		} `yaml:"local"`
	} `yaml:"etcd"`
	ImageRepository   string `yaml:"imageRepository"`
	Kind              string `yaml:"kind"`
	KubernetesVersion string `yaml:"kubernetesVersion"`
	Networking        struct {
		DNSDomain     string `yaml:"dnsDomain"`
		PodSubnet     string `yaml:"podSubnet"`
		ServiceSubnet string `yaml:"serviceSubnet"`
	} `yaml:"networking"`
	Scheduler struct {
	} `yaml:"scheduler"`
}

//ExtraVolume object
type ExtraVolume struct {
	HostPath  string `yaml:"hostPath"`
	MountPath string `yaml:"mountPath"`
	Name      string `yaml:"name"`
	PathType  string `yaml:"pathType"`
	ReadOnly  bool   `yaml:"readOnly"`
}

//NewKubeadm returns Kubeadm object
func NewKubeadm() (*Kubeadm, error) {
	var k Kubeadm
	return &k, nil
}

// readKubeadmConfig unmarshals YAML file to K8s ConfigMap obejct
func (k *Kubeadm) readKubeadmConfig(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	err = yaml.Unmarshal([]byte(f), &k.Configmap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	err = yaml.Unmarshal([]byte(k.Configmap.Data.ClusterConfiguration), &k.ClusterConfiguration)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	// fmt.Printf("%v", l.ClusterConfiguration.APIServer.ExtraArgs.CloudConfig)
	return nil
}

//writeKubeadmConfigFile from Configmap object to file
func (k *Kubeadm) writeKubeadmConfigFile(path string) error {
	k.ClusterConfiguration.APIServer.ExtraArgs.CloudConfig = filepath.Join(K8sConfigDir, vsphere.CloudConf)
	configByte, err := yaml.Marshal(&k.ClusterConfiguration)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %v", err)
	}

	k.Configmap.Data.ClusterConfiguration = string(configByte)
	configmapByte, err := yaml.Marshal(&k.Configmap)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %v", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	_, err = f.Write(configmapByte)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %v", path, err)
	}
	return nil
}

//enableCloudProvider update ClouderConfiguration object with cloudprovider enabled
func (k *Kubeadm) enableCloudProvider(cloudProvider string) error {
	if k.ClusterConfiguration.APIServer.ExtraArgs.CloudProvider == cloudProvider {
		return nil
	}
	k.ClusterConfiguration.APIServer.ExtraArgs.CloudProvider = cloudProvider
	cloudConfig := path.Join(K8sConfigDir, cloudProvider+".conf")
	k.ClusterConfiguration.APIServer.ExtraArgs.CloudConfig = cloudConfig
	k.ClusterConfiguration.APIServer.ExtraVolumes = append(
		k.ClusterConfiguration.APIServer.ExtraVolumes,
		ExtraVolume{
			HostPath:  cloudConfig,
			MountPath: cloudConfig,
			Name:      "cloud-config",
			PathType:  "FileOrCreate",
			ReadOnly:  true,
		},
	)
	k.ClusterConfiguration.ControllerManager.ExtraArgs.CloudConfig = cloudConfig
	k.ClusterConfiguration.ControllerManager.ExtraArgs.CloudProvider = cloudProvider
	k.ClusterConfiguration.ControllerManager.ExtraVolumes = append(
		k.ClusterConfiguration.ControllerManager.ExtraVolumes,
		ExtraVolume{
			HostPath:  cloudConfig,
			MountPath: cloudConfig,
			Name:      "cloud-config",
			PathType:  "FileOrCreate",
			ReadOnly:  true,
		},
	)
	return nil
}
