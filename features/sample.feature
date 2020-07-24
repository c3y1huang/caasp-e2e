@sample
Feature: Sample

  Background:
    Given "kubectl" installed on local machine
    And environment variables exported:
      """
      GODOG_CLUSTER_JSON_FILE
      """

  Scenario: Cluster info
    When I search "cluster-0.platform" in cluster info
    Then I found:
    """
    vmware
    """
    When I search "cluster-0.lb.ip" in cluster info
    Then I found:
    """
    10.84.72.44
    """
