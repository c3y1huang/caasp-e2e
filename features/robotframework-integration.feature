Feature: RobotFramework Integration

  Background:
    Given "kubectl" installed on local machine
    And cluster access from "../cluster/cluster_1/admin.conf"
    And cluster info from "cluster/logs/cluster_state.json"

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
