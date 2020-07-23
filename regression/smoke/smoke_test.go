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

package smoke_test

import (
	"github.com/cucumber/godog"

	cluster "github.com/c3y1huang/caasp-e2e/features/cluster"
	local "github.com/c3y1huang/caasp-e2e/features/local"
	remote "github.com/c3y1huang/caasp-e2e/features/remote"
	kubeadm "github.com/c3y1huang/caasp-e2e/features/kubeadm"
	kubectl "github.com/c3y1huang/caasp-e2e/features/kubectl"
	vsphere "github.com/c3y1huang/caasp-e2e/features/vsphere"
)

func InitalizeScenarioContext(ctx *godog.ScenarioContext) {
	cluster.InitializeScenario(ctx)
	local.InitializeScenario(ctx)
	remote.InitializeScenario(ctx)
	kubeadm.InitializeScenario(ctx)
	kubectl.InitializeScenario(ctx)
	vsphere.InitializeScenario(ctx)
}