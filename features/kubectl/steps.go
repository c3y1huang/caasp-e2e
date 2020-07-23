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

package kubectl

import (
	"fmt"

	"github.com/cucumber/godog"

	"github.com/c3y1huang/caasp-e2e/features/local"
)

//Kubectl object
type Kubectl struct {
	Local local.Local
}

//NewKubectl returns Kubectl object
func NewKubectl() (*Kubectl, error) {
	var k Kubectl
	return &k, nil
}

func (k *Kubectl) createK8sResourceWithManifest(resource string, manifest *godog.DocString) error {
	return k.Local.RunCmdSimple("", fmt.Sprintf("cat <<EOF | kubectl apply -f -\n%s\nEOF", manifest.Content))
}
