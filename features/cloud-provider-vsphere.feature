@vSphereCloudProvider
Feature: vSphere Cloud Provider

  Background:
    Given "kubectl" installed on local machine
    And environment variables exported:
      """
      KUBECONFIG
      GODOG_CLUSTER_JSON_FILE
      """
    And I run command to "verify cluster running":
      """
      kubectl cluster-info
      """

  @vSphereCloudProviderNodeMeta
  Scenario: vSphere Cloud Provider Node Meta
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

  @vSphereCloudProviderPersistentVolume
  Scenario Outline: vSphere Cloud Provider Persistent Volume
    Given "govc" installed on local machine
    * I create "<diskSize>" "<disk>.vmdk" in vSphere "<datacenter>" and "<datastore>"
    And "kubectl" installed on local machine
    When I create "persistent volume" with manifest:
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
    And I create "persistent volume claim" with manifest:
      """
      apiVersion: v1
      kind: PersistentVolumeClaim
      metadata:
        name: disk-pv-pvc
        labels:
          app: sample-disk-pv
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
          app: sample-disk-pv
          tier: sample-disk-pv
      spec:
        selector:
          matchLabels:
            app: sample-disk-pv
            tier: sample-disk-pv
        strategy:
          type: Recreate
        template:
          metadata:
            labels:
              app: sample-disk-pv
              tier: sample-disk-pv
          spec:
            tolerations:
            - effect: NoSchedule
              key: node-role.kubernetes.io/master
            containers:
            - image: busybox
              name: sample-disk-pv
              volumeMounts:
              - name: pv-volume
                mountPath: /data
              command: [ "sleep", "infinity" ]
            volumes:
            - name: pv-volume
              persistentVolumeClaim:
                claimName: disk-pv-pvc
      """
    Then I run command to "get pvc":
      """
      kubectl get pvc | sed 1d 
      """
    And output contains "Bound"
    And I wait until "pod" "Running":
      """
      kubectl get pod | sed 1d
      """
    And I run command to "add data":
      """
      POD=$(kubectl get pod --selector=app=sample-disk-pv -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "echo 'I am here' > /data/hello.txt"
      """
    And I run command to "see data input":
      """
      POD=$(kubectl get pod --selector=app=sample-disk-pv -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "cat /data/hello.txt"
      """
    And output contains "I am here"

    Examples:
    |                  disk | diskSize | datacenter | datastore |
    | c3y1-storage/c3y1Test |       1G |      PROVO |      3PAR |
    
  @vSphereCloudProviderStorageClass
  Scenario Outline: vSphere Cloud Provider Storage Class
    Given "kubectl" installed on local machine
    When I create "storage class" with manifest:
      """
      kind: StorageClass
      apiVersion: storage.k8s.io/v1
      metadata:
        name: disk-sc
        annotations:
          storageclass.kubernetes.io/is-default-class: "true"
      provisioner: kubernetes.io/vsphere-volume
      parameters:
        datastore: "<datastore>"
      """
    And I create "persistent volume claim" with manifest:
      """
      apiVersion: v1
      kind: PersistentVolumeClaim
      metadata:
        name: disk-sc-pvc
        labels:
          app: sample-disk-sc
      spec:
        accessModes:
        - ReadWriteOnce
        resources:
          requests:
            storage: <diskSize>i
      """
    And I create "deployment" with manifest:
      """
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: disk-sc-deployment
        labels:
          app: sample-disk-sc
          tier: sample-disk-sc
      spec:
        selector:
          matchLabels:
            app: sample-disk-sc
            tier: sample-disk-sc
        strategy:
          type: Recreate
        template:
          metadata:
            labels:
              app: sample-disk-sc
              tier: sample-disk-sc
          spec:
            tolerations:
            - effect: NoSchedule
              key: node-role.kubernetes.io/master
            containers:
            - image: busybox
              name: sample-disk-sc
              volumeMounts:
              - name: sc-volume
                mountPath: /data
              command: [ "sleep", "infinity" ]
            volumes:
            - name: sc-volume
              persistentVolumeClaim:
                claimName: disk-sc-pvc
      """
    Then I wait until "pvc" "Bound":
      """
      kubectl get pvc | sed 1d 
      """
    And I wait until "pod" "Running":
      """
      kubectl get pod | sed 1d
      """
    And I run command to "add data":
      """
      POD=$(kubectl get pod --selector=app=sample-disk-sc -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "echo 'I am here' > /data/hello.txt"
      """
    And I run command to "see data input":
      """
      POD=$(kubectl get pod --selector=app=sample-disk-sc -o jsonpath='{.items[*].metadata.name}') &&
      kubectl exec -it ${POD} -- sh -c "cat /data/hello.txt"
      """
    And output contains "I am here"

    Examples:
    | diskSize  | datastore |
    |       1G  |      3PAR |
