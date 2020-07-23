@enableCloudProvider
Feature: Enable Cloud Provider

  Background:
    Given environment variables exported:
      """
      KUBECONFIG
      GODOG_CLUSTER_JSON_FILE
      GODOG_VSPHERE_CLUSTER_VCENTER
      GODOG_VSPHERE_CLUSTER_PREFIX
      GODOG_VSPHERE_CLUSTER_FOLDER
      GODOG_VSPHERE_CLUSTER_DATACENTER
      GODOG_VSPHERE_CLUSTER_DATASTORE
      GODOG_VSPHERE_CLUSTER_RESOURCEPOOL
      GODOG_VSPHERE_CLUSTER_REGION
      GODOG_VSPHERE_CLUSTER_ZONE
      GODOG_VSPHERE_USER
      GODOG_VSPHERE_PASSWORD
      """
    And "kubectl" installed on local machine
    And I run command to "verify cluster running":
      """
      kubectl cluster-info
      """

  @enableCloudProviderVSphere
  Scenario: Enable vSphere Cloud Provider
    When I run command to "create cluster folder":
      """
      DATACENTER=$GODOG_VSPHERE_CLUSTER_DATACENTER
      FOLDER=$GODOG_VSPHERE_CLUSTER_FOLDER
      govc folder.create /$DATACENTER/vm/$FOLDER
      """
    And I run command to "move virtual machines to cluster folder":
      """
      DATACENTER=$GODOG_VSPHERE_CLUSTER_DATACENTER
      PREFIX=$GODOG_VSPHERE_CLUSTER_PREFIX
      FOLDER=$GODOG_VSPHERE_CLUSTER_FOLDER
      govc object.mv /$DATACENTER/vm/$PREFIX-\* /$DATACENTER/vm/$FOLDER
      """
    And I run command to "enable disk.UUID on all vms":
      """
      DATACENTER=$GODOG_VSPHERE_CLUSTER_DATACENTER
      PREFIX=$GODOG_VSPHERE_CLUSTER_PREFIX
      VMS=("$PREFIX-master-0" "$PREFIX-master-1" "$PREFIX-worker-0")
    
      function setup {
        NAME=$1
        echo "[$NAME]"
        govc vm.power -dc=$DATACENTER -off $NAME
     
        govc vm.change -dc=$DATACENTER -vm=$NAME -e="disk.enableUUID=1" &&\
          echo "Configured disk.enabledUUID: 1"
     
        govc vm.power -dc=$DATACENTER -on $NAME
      }
     
      for vm in ${VMS[@]}
      do
        setup $vm &
      done
      wait
      
      # Allow virtual machines to react to the bootup
      sleep 10
      """
    And I wait until "nodes" "Ready":
      """
      kubectl get nodes | sed 1d
      """
    And I run command to "update providerID":
      """
      DATACENTER=$GODOG_VSPHERE_CLUSTER_DATACENTER
      PREFIX=$GODOG_VSPHERE_CLUSTER_PREFIX
      for vm in $(govc ls "/$DATACENTER/vm/$PREFIX-cluster"); do
        VM_INFO=$(govc vm.info -json -dc=$DATACENTER -vm.ipath="/$vm" -e=true)
        VM_NAME=$(jq -r ' .VirtualMachines[] | .Name' <<< $VM_INFO)
        [[ $VM_NAME == *"-lb-"* ]] && continue
        VM_UUID=$(jq -r ' .VirtualMachines[] | .Config.Uuid' <<< $VM_INFO)
        echo "Patching $VM_NAME with UUID:$VM_UUID"
        sleep 1
        kubectl patch node $VM_NAME -p "{\"spec\":{\"providerID\":\"vsphere://$VM_UUID\"}}"
      done
      """
    And I create vsphere.conf
    And I scp "vsphere.conf" to "cluster-0.master.master-0, cluster-0.master.master-1, cluster-0.worker.worker-0" directory "/tmp"
    And I run command on "cluster-0.master.master-0, cluster-0.master.master-1, cluster-0.worker.worker-0" to "move vsphere.conf":
      """
      sudo mv /tmp/vsphere.conf /etc/kubernetes/
      """
    And I run command to "save kubeadm-config.conf":
      """
      kubectl -n kube-system get cm/kubeadm-config -o yaml > kubeadm-config.conf
      """
    And I open configmap kubeadm-config from "kubeadm-config.conf"
    * enable "vsphere" cloud provider in configmap kubeadm-config
    * save configmap kubeadm-config to "kubeadm-config.conf"
    And I run command to "apply kubeadm-config.conf":
      """
      kubectl apply -f kubeadm-config.conf
      """
    And I run command on "cluster-0.master.master-0, cluster-0.master.master-1" to "update kubelet":
      """
      sudo systemctl stop kubelet
      source /var/lib/kubelet/kubeadm-flags.env
      echo KUBELET_KUBEADM_ARGS='"'--cloud-config=/etc/kubernetes/vsphere.conf --cloud-provider=vsphere $KUBELET_KUBEADM_ARGS'"' > /tmp/kubeadm-flags.env
      sudo mv /tmp/kubeadm-flags.env /var/lib/kubelet/kubeadm-flags.env
      sudo systemctl start kubelet
      """
    And I run command on "cluster-0.master.master-0, cluster-0.master.master-1" to "update control-plane":
      """
      sudo kubeadm upgrade node phase control-plane --etcd-upgrade=false
      """
    And I run command on "cluster-0.worker.worker-0" to "update kubelet":
      """
      sudo systemctl stop kubelet
      source /var/lib/kubelet/kubeadm-flags.env
      echo KUBELET_KUBEADM_ARGS='"'--cloud-config=/etc/kubernetes/vsphere.conf --cloud-provider=vsphere $KUBELET_KUBEADM_ARGS'"' > /tmp/kubeadm-flags.env
      sudo mv /tmp/kubeadm-flags.env /var/lib/kubelet/kubeadm-flags.env
      sudo systemctl start kubelet
      """
    Then I wait until "nodes" "Ready":
      """
      kubectl get nodes | sed 1d
      """