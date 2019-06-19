# WindowsOperator
## Description
Listens for the additional of nodes, and if the node is a windows node, it will be auto installed if needed.
User creates a new kubernets node object, and the rest of the process is fully automated.

## WinOperator Components

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


