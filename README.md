# caasp-e2e
CaaSP e2e tests with godog

# Requirement
* [Godog](https://github.com/cucumber/godog)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [govc](https://github.com/vmware/govmomi/tree/master/govc) (vSphere)

# Run test
You need a running cluster before execution:
```bash
godog
```

# Write Tests
## Cluster Access
Use `cluster access from` in each test *Feature* to specify the kubernetes configuration file. For example
```
cluster access from "../cluster/cluster_1/admin.conf"
```

## JSON file
Tests can take JSON files and with uses of dot notation in tests to retrieve the variables. This is perticulary useful when you have pre-existing cluster.

For example:
```
  Background:
    Given cluster info from "cluster/logs/cluster_state.json"

  Scenario: Cluster info
    When I search "cluster_1.platform" in cluster info
    Then it prints:
    """
    vmware
    """
```


