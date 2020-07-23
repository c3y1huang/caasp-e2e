[![CircleCI](https://circleci.com/gh/c3y1huang/caasp-e2e/tree/master.svg?style=svg)](https://circleci.com/gh/c3y1huang/caasp-e2e/tree/master)

# CaaSP-e2e
CaaSP e2e tests with godog.

## Requirement
* [godog](https://github.com/cucumber/godog)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [govc](https://github.com/vmware/govmomi/tree/master/govc) (vSphere)

## Run test
You need a running cluster:
```bash
cd regression/smoke
godog\
 --tags "@enableCloudProviderVSphere,@vSphereCloudProviderStorageClass"\
 ../../features/
```

## Write Tests

### Implement Features Tests and Steps
Create:
```
/feature/<feature_name>.feature // <1>
/feature/<feature_name>/<feature_name>.go <2>
/feature/<feature_name>/<feature_name>_test.go <3>
/feature/<feature_name>/steps.go <4>
```
<1> Test scenarios.
<2> The steps maps to the test scenarios.
<3> Entry point to the steps.
<4> Step implementations.

Then import to the `/regression/smoke/smoke_test.go`.

Original skeleton reference: https://github.com/innobead/caasp-e2e

### Get Value from JSON file
Tests use `.` notation to get JSON data from file. This is particularly useful when you have infrastructure provisioned automatically with dynamic values.

```
  Background:
    Given cluster info from "cluster/logs/cluster_state.json"

  Scenario: Cluster info
    When I search "caasp-cluster.platform" in cluster info
    Then it prints:
    """
    vmware
    """

    When I search "caasp-cluster.master.caasp-master-0.ip" in cluster info
    Then it prints:
    """
    10.84.73.87
    """
```

