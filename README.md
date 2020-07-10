# CaaSP-e2e
CaaSP e2e tests with godog.

## Requirement
* [godog](https://github.com/cucumber/godog)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [govc](https://github.com/vmware/govmomi/tree/master/govc) (vSphere)

## Run test
You need a running cluster before execution:
```bash
godog
```

## Write Tests
### Cluster Access
Use `cluster access from` in each test *Feature* to specify the kubernetes configuration file path. For example
```
Background:
    cluster access from "../cluster/cluster_1/admin.conf"
```

### JSON file
Tests use `.` notation in tests to get JSON data from file. This is particularly useful when you have infrastructure provisioned automatically with dynamic values.

For example:
```json
{
  "cluster_1": {
    "platform": "vmware",
    "master": {
      "caasp-cluster-abc-1-master-1": {
        "ip": "10.84.73.87",
        "disable": "False",
        "skuba_name": "010084073087"
      }
    }
  }
}
```
```
  Background:
    Given cluster info from "cluster/logs/cluster_state.json"

  Scenario: Cluster info
    When I search "cluster_1.platform" in cluster info
    Then it prints:
    """
    vmware
    """

    When I search "cluster_1.master.caasp-cluster-abc-1-master-1.ip" in cluster info
    Then it prints:
    """
    10.84.73.87
    """
```

