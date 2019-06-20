# WindowsOperator
## Description
Provides a simple fully integrated way to management the addition of Windows Nodes to OpenShift clusters. Using a user-provided Windows machine that has WIndows WMI enabled, it will proceded to upgrade/install and configure the windows node, and join it to cluster. The user provides administrator username and password to the machine, or user configured default value. User creates a simple file, and uses oc create -f newnode.yaml which starts the process. User can watch the install via Node Events in OpenShift.

## Installation of WinOperator
Administartor installs 2 containers via simple script or uses OperatorHub.



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
Components - Ignition Files

#### Template:
```
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
```
 
### Component
A component consists of a [Ignition](https://coreos.com/ignition/docs/latest/configuration-v2_1.html) file. The WinOperator using a subset of the ignition format, which was originally designed for RHCOS for the defenition of components. OCP 4.x uses ignition for the configuration of nodes from boot, and WInOperator extends that to Windows Nodes. As the original ignition specification is highly specific to Linux, and reimplmentation was done for Windows, called [libigniton](https://github.com/glennswest/libignition).  

A Component  - Ignition File with a embedded powershell script:
```
{ 
  "ignition": { 
    "version": "2.2.0" 
  }, 
  "storage": { 
    "files": [ 
      { 
        "path": "/bin/metadata/prewin1809_v1.0.metadata", 
        "filesystem": "", 
        "mode": 420, 
        "contents": { 
          "source": "data:text/plain;charset=utf-8;base64,ewog...KfQo=" 
        } 
      }, 
      { 
        "path": "/bin/prereq1809.ps1", 
        "filesystem": "", 
        "mode": 420, 
        "contents": { 
          "source": "data:text/plain;charset=utf-8;base64,JEVy...Cgo=" 
        } 
      } 
    ] 
  } 
} 
```

In this example, the component has a embedded powershell script, and a the mandatory metadata file. The metadata file gives us the details of the component, as well as the contents of the powershell script. Data needed for a component can be directl in the content section, base64 encoded in the content section, or remotely via url. 

A Component metadata file:
```
{
  "name":          "prewin1809",
  "version":       "v1.0",
  "description":    "Prequisite Setup for WIndows 1809",
  "install_message": "",
  "package_url":     "",
  "install": {
      "commands": ["\\bin\\prereq1809.ps1"],
      "reboot":   "no"
      }
  "uninstall": {
      "priority": 100,
      "lprecmds": [],
      "commands": [],
      "lpstcmds": [],
      "reboot":   "no"
      }
  "pre_upgrade": {
      "priority": 100,
      "uninstall": "no",
      "lprecmds": [],
      "commands": [],
      "lpstcmds": [],
      "reboot":   "no"
      }
  "post_upgrade": {
      "priority": 100,
      "lprecmds": [],
      "commands": [],
      "lpstcmds": [],
      "reboot":   "no"
      }
  "files": []
}
```
In this example you can see the 4 main sections of the metadata. Install is the commands needed for install of a new node, uninstall is for removing all the components added. pre_upgrade is run before a upgrade of the component, and post_upgrade is run after the upgrade of a component. 

#### Adding more content
To add additional content, create a new json file, with the name,version, and at least the install section. The new component will needed to be added to a template, and that template used by a new node install.



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

### WinOperator EcoSystem
[winoperator](https://github.com/glennswest/winoperator) - This project
[winmachineman](https://github.com/glennswest/winmachineman) - Manages template and lifecycle of winnodemanager 
[winnodemanager](https://github.com/glennswest/winnodemanager) - Manages life cycle of windows machine 
[winoperatordata](https://github.com/glennswest/winoperatordata) - The components in ignition format plus templates
[libignition](https://github.com/glennswest/libignition) - Ignition Lite - Used to parse and generation ignition files for Windows
[libpowershell](https://github.com/glennswest/libpowershell) - Used to communicate to windows


