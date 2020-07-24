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

package remote

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/cucumber/godog"
	"github.com/tidwall/gjson"
	"github.com/tmc/scp"

	"github.com/c3y1huang/caasp-e2e/features/cluster"
	"github.com/c3y1huang/caasp-e2e/features/local"
)

const (
	//DefaultPort to use for SSH session
	DefaultPort = 22
	//DefaultUser to use for SSH session
	DefaultUser = "sles"
)

//Client is SSH client object
type Client struct {
	Host    string
	Port    int64
	Config  *ssh.ClientConfig
	Cluster []byte
	Results []Result
}

//Result of SSH stdout object
type Result struct {
	IP     string
	Stdout string
}

//NewClient returns ClientConfig object
func NewClient() (*Client, error) {
	var c Client

	//host port
	c.Port = DefaultPort

	//SSH client config
	privateBytes, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		return &c, fmt.Errorf("failed to load private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return &c, fmt.Errorf("failed to parse private key: %v", err)
	}
	c.Config = &ssh.ClientConfig{
		User: DefaultUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	//reads cluster JSON data from file
	c.Cluster, err = local.ReadFileFromEnv(cluster.EnvJSONFile)
	if err != nil {
		return &c, err
	}

	return &c, nil
}

//getClusterIPs returns list of IPs from matching hostname in cluster state JSON string
func (c *Client) getClusterIPs(hosts []string) []string {
	var ips []string
	for _, host := range hosts {
		ip := gjson.Get(string(c.Cluster), host+".ip")
		ips = append(ips, ip.String())
	}
	return ips
}

//RunCmdBlockOnHosts run commands on remote hosts and save outputs in Result object
func (c *Client) RunCmdBlockOnHosts(hosts, description string, input *godog.DocString) error {
	c.Results = []Result{}
	hosts = strings.ReplaceAll(hosts, " ", "")

	var hostList []string
	hostList = strings.Split(hosts, ",")

	var wg sync.WaitGroup
	for _, ip := range c.getClusterIPs(hostList) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			c.Results = append(
				c.Results,
				Result{
					IP:     ip,
					Stdout: c.RunCmd(input.Content, ip),
				},
			)
		}(ip)
	}
	wg.Wait()

	return nil
}

//RunCmd create ssh session and run command on single host and return stdout
func (c *Client) RunCmd(cmd, host string) string {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, c.Port), c.Config)
	if err != nil {
		fmt.Printf("failed to create connection with %s: %v\n", host, err)
		return ""
	}
	defer conn.Close()

	var stdoutBuf bytes.Buffer
	session, _ := conn.NewSession()
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return stdoutBuf.String()
}

//UploadFileToHosts copy file from local machine to hosts
func (c *Client) UploadFileToHosts(file, hosts, dir string) error {
	c.Results = []Result{}
	hosts = strings.ReplaceAll(hosts, " ", "")

	var hostList []string
	hostList = strings.Split(hosts, ",")

	var wg sync.WaitGroup
	for _, ip := range c.getClusterIPs(hostList) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			c.UploadFile(file, ip, dir)
		}(ip)
	}
	wg.Wait()

	return nil
}

//UploadFile copy file from local machine to single host
func (c *Client) UploadFile(src, host, dst string) error {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, c.Port), c.Config)
	if err != nil {
		return fmt.Errorf("failed to create connection with %s: %v", host, err)
	}
	defer conn.Close()

	session, _ := conn.NewSession()
	err = scp.CopyPath(src, dst, session)
	if err != nil {
		return fmt.Errorf("copy file to remote destination: %v", err)
	}

	return nil
}
