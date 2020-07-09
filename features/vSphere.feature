Feature: vSphere Cloud Provider

  Background:
    Given "kubectl" installed on local machine
    And I can run command to "verify cluster running":
      """
      kubectl cluster-info
      """

  Scenario: vSphere Cloud Provider
    When I run command to "get node ProviderID":
      """
      kubectl describe nodes | grep 'ProviderID'
      """
    Then output contains "vsphere://"
    When I run command to "get node region":
      """
      kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\tregion: "}{.metadata.labels.failure-domain\.beta\.kubernetes\.io/region}{"\n"}{end}'
      """
    Then output contains "vcp-provo"
    When I run command to "get node zone":
      """
      kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\tzone: "}{.metadata.labels.failure-domain\.beta\.kubernetes\.io/zone}{"\n"}{end}'
      """
    Then output contains "vcp-cluster-jazz"

  Scenario Outline: vSphere Cloud Provider Persistent Volume
    Given "govc" installed on local machine
    * I create "<diskSize>" "<disk>.vmdk" in vSphere "<datacenter>" and "<datastore>"
    And "kubectl" installed on local machine
    * I create "persistent volume" with manifest:
    """
    apiVersion: v1
    kind: PersistentVolume
    metadata:
      name: disk-pv
    spec:
      capacity:
        storage: <diskSize>i
      accessModes:
        - ReadWriteOnce
      persistentVolumeReclaimPolicy: Delete
      vsphereVolume:
        volumePath: "[<datastore>] <disk>.vmdk"
        fsType: ext4
    """
    * I create "persistent volume claim" with manifest:
    """
    apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      name: disk-pv-pvc
      labels:
        app: sample
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: <diskSize>i
    """
    * I create "deployment" with manifest:
    """
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: disk-pv-deployment
      labels:
        app: sample
        tier: sample
    spec:
      selector:
        matchLabels:
          app: sample
          tier: sample
      strategy:
        type: Recreate
      template:
        metadata:
          labels:
            app: sample
            tier: sample
        spec:
          containers:
          - image: busybox
            name: sample
            volumeMounts:
            - name: sample-volume
              mountPath: /data
            command: [ "sleep", "infinity" ]
          volumes:
          - name: sample-volume
            persistentVolumeClaim:
              claimName: disk-pv-pvc
    """
    When I run command to "get pvc":
      """
      kubectl get pvc | sed 1d 
      """
    Then output contains "Bound"
    And I wait until "pod" "Running":
      """
      kubectl get pod | sed 1d
      """
    When I run command to "add data":
      """
      POD=$(kubectl get pod --selector=app=sample -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "echo 'I am here' > /data/hello.txt"
      """
    Then I run command to "see data input":
      """
      POD=$(kubectl get pod --selector=app=sample -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "cat /data/hello.txt"
      """

    Examples:
    |          disk | diskSize | datacenter | datastore |
    | c3y1/c3y1Test |       1G |      PROVO |      3PAR |
    
