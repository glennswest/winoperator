# WindowsOperator
## Description
Provides a simple fully integrated way to management the addition of Windows Nodes to OpenShift clusters. Using a user-provided Windows machine that has WIndows WMI enabled, it will proceded to upgrade/install and configure the windows node, and join it to cluster. The user provides administrator username and password to the machine, or user configured default value. User creates a simple file, and uses oc create -f newnode.yaml which starts the process. User can watch the install via Node Events in OpenShift.

## Installation of WinOperator
Administartor installs 3 containers via simple script or uses OperatorHub.



## How It Works
![alt text](https://raw.githubusercontent.com/glennswest/winoperator/master/doc/overviewuml.png)
 
General Flow of Operation
![alt text](https://raw.githubusercontent.com/glennswest/winoperator/master/doc/overview.png)

## WinOperator 
The Windows Operator has two types of components, infrastructure, and content.
The infrastructure consists of 2 containers
   winoperator
   winmachineman
The content consists of a data container:
   winoperatordata

The infrastructor components are slow change items. There not expected to change much over time, as they provide the execution environment for the content.
The content is expected to change with windows versions as well as openshift changes. 

### WinOperatorData
This is a data container, containing the components for multiple versions of windows. This includes two types of data.
Templates - A json format file that lists components versions for a version of windows.

A Template:
{
  "name":          "win2019",
  "version":       "0.001",
  "description":    "Initial Version of Win2019",
  "install_message": "OpenShift Windows 4.x Windows 2019 ",
  "packages": ["prewin1809_v1.0",
               "docker_v1.0.1",
               "pause_v1.0.1",
               "nssm_2.24",
               "node_1.0.2",
               "cni_0.3.1",
               "sdn_v1.0.1",
               "ovn_16e1a3cf",
               "ovs_2.70beta",
               "kube_v1.11.3"
               ]
}



### WinOperator
Listens for creation of node from kubernetes
Builds parameters for installation from local object store for auth.
Converts annotations and labels to parameters
Sends Restful request to WinMachineManager

### WinMachineManager
1. Validates connection to windows machine. 
2. Lifecycle Management of WindowsNodeManager (Its contained in WinMachineManager)
    a. Uninstalled WinMachineMan if already installed
    b. Installs latest version of WInMachineMan
3. Selects Template and Finalize parameters
4. Sends Rest Request to WindowsNode Mangager
5. Optionally Serves content - Templates and Components


### WindowsNodeManager
1. Downloads Ignition Files
2. Deploys Ignition Files
3. Handles Metadata 
   a. Execute powershell commands that are needed
   b. Execute master commands - before and after powershell
4. Handles reboot and continuation
5. Support Install, Uninstall and Update

### WinOperator Techical Flow
![alt text](https://raw.githubusercontent.com/glennswest/winoperator/master/doc/winoperator.png)


