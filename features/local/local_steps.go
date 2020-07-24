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

package local

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

const (
	waitRetry = 160
	waitInterval = 5*time.Second
	//MissingEnvVariable error message
	MissingEnvVariable = "please setup environment variable before running test:\n export %s=\"\""
)

//Local object
type Local struct {
	InterceptedStdout bytes.Buffer
}

//NewLocal returns Local object
func NewLocal() (*Local, error) {
	var local Local
	return &local, nil
}

//AssertEnvsExist checks if all environment variables have value
func AssertEnvsExist(input *godog.DocString) error {
	envs := strings.Split(string(input.Content), "\n")
	for _, name := range envs {
		value, ok := os.LookupEnv(strings.TrimSpace(name))
		if !ok || len(value) == 0 {
			return fmt.Errorf(MissingEnvVariable, name)
		}
	}
	return nil
}

//IsInstalled returns error if file not found in PATH
func IsInstalled(name string) error {
	_, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s is not installed: %v", name, err)
	}

	return nil
}

//RunCmdBlock gets commands as DocString, execute and store stdout in Local object
func (l *Local) RunCmdBlock(desc string, input *godog.DocString) error {
	l.InterceptedStdout.Reset()
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", input.Content)
	cmd.Stdout = &l.InterceptedStdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %v", cmd, err)
	}
	// l.InterceptedStdout = cmd.String()
	return nil
}

//RunCmdSimple gets commands as string, execute and store stdout in Local object
func (l *Local) RunCmdSimple(desc, input string) error {
	l.InterceptedStdout.Reset()
	// cmd, err := RunCmd(input, &l.InterceptedStdout)
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", input)
	cmd.Stdout = &l.InterceptedStdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %s: %v", cmd, err)
	}
	return nil
}

//AssertOutputContains matching string in every lines of the intercepted stdout
func (l *Local) AssertOutputContains(match string) error {
	if len(l.InterceptedStdout.String()) == 0 {
		return fmt.Errorf("expect output")
	}

	lines := strings.Split(l.InterceptedStdout.String(), "\n")

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if !strings.Contains(line, match) {
			return fmt.Errorf("expect output to contain %s\nactual:\n%s", match, line)
		}
	}
	return nil
}

//WaitUntilOutputContains matching string in every lines of the intercepted stdout
func (l *Local) WaitUntilOutputContains(resource, match string, input *godog.DocString) error {
	err := retry(waitRetry, waitInterval, func() (err error) {
		l.RunCmdSimple("", input.Content)
		err = l.AssertOutputContains(match)
		return
	})
	if err != nil {
		return err
	}
	return nil
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

//TranslateAbsolutePathFromEnv checks if give environment variable 
//resolves to a valid path. Returns absolute path.
func TranslateAbsolutePathFromEnv(name string) (string, error) {
	fpath, ok := os.LookupEnv(name)
	if !ok {
		return "", fmt.Errorf(MissingEnvVariable, name)
	}
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return "", err
	}
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

//ReadFileFromEnv validate and opens file from environment variable 
func ReadFileFromEnv(name string) ([]byte, error) {
	var f []byte
	fpath, err := TranslateAbsolutePathFromEnv(name)
	if err != nil {
		return f, err
	}
	f, err = ioutil.ReadFile(fpath)
	if err != nil {
		return f, err
	}
	return f, nil
}
