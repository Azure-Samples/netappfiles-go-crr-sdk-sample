---
page_type: sample
languages:
- go
products:
- azure
- azure-netapp-files
description: "This sample project demonstrates how to enable cross-region replication (CRR) on NFSv3 volume using Azure Go SDK with Microsoft.NetApp resource provider."
---


# Azure NetApp Files CRR SDK Sample for Go

This sample demonstrates how to enable cross-region replication (CRR) on NFSv3 volume using Azure Go SDK for Microsoft.NetApp resource provider. This process is identical for NFS v4.1 volumes with exception of protocol type used.

In this sample application we perform the following operations:

* Creation
  * Primary NetApp account
    * Primary capacity pool
      * Primary NFS v3 volume
  * Secondary NetApp account
    * Secondary capacity pool
      * Secondary NFS v3 Data Replication volume with reference to the primary volume Resource ID
* Authorize primary volume with secondary volume Resource ID
* Clean up created resources (not enabled by default) 

If you don't already have a Microsoft Azure subscription, you can get a FREE trial account [here](http://go.microsoft.com/fwlink/?LinkId=330212).

## Prerequisites

1. Go installed \(if not installed yet, follow the [official instructions](https://golang.org/dl/)\)
2. Azure Subscription.
3. Subscription needs to have Azure NetApp Files resource provider registered. For more information, see [Register for NetApp Resource Provider](https://docs.microsoft.com/en-us/azure/azure-netapp-files/azure-netapp-files-register).
4. Resource Group(s) created
5. Virtual Networks (both for primary and secondary volumes) with a delegated subnet to Microsoft.Netapp/volumes resource. For more information, see [Guidelines for Azure NetApp Files network planning](https://docs.microsoft.com/en-us/azure/azure-netapp-files/azure-netapp-files-network-topologies).
6. Adjust variable contents within `var()` block at `example.go` file to match your environment
7. For this sample Go console application to work, authentication is needed. The chosen method for this sample is using service principals:
    * Within an [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/quickstart) session, make sure you're logged on at the subscription where you want to be associated with the service principal by default: 

      ```bash
      az account show
      ```

      If this is not the correct subscription, use: 

      ```bash
      az account set -s <subscription name or id>  
      ```

    * Create a service principal using Azure CLI: 

      ```bash
      az ad sp create-for-rbac --sdk-auth
      ```

      >Note: This command will automatically assign RBAC contributor role to the service principal at subscription level. You can narrow down the scope to the specific resource group where your tests will create the resources.

    * Copy the output content, paste it in a file called azureauth.json, and secure it with file system permissions (make sure it is not inside of any repo).
    * Set an environment variable pointing to the file path you just created. Here is an example with Powershell and bash:

      Powershell

      ```powershell
      [Environment]::SetEnvironmentVariable("AZURE_AUTH_LOCATION", "C:\sdksample\azureauth.json", "User")
      ```

      Bash

      ```bash
      export AZURE_AUTH_LOCATION=/sdksamples/azureauth.json
      ```

    >Note: For other Azure Active Directory authentication methods for Go, see [Authentication methods in the Azure SDK for Go](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization).

## What does example.go do

This sample project demonstrates how to enable cross-region replication in Azure NetApp Files for an NFSv3 enabled volume. (Note that this process is the same for NFSv4.1 volumes.) Similar to other examples, the authentication method is based on a service principal. This project will create two NetApp accounts in different regions, each with a capacity pool. A single volume using the Premium service level as the Primary volume, and the Standard service level in the secondary region with Data Protection object.

The process of enabling cross-region replication involves creating the primary resources, including primary volume.  Then we continue to create the secondary resources, but the secondary volume needs to contain the Data Replication Object. After this step, we authorize the replication from the primary volume, referencing the resource ID of the secondary volume.

In addition, we use some non-sensitive information from the *file-based authentication* file where we initially get the subscription ID. This information is used for the test to check if the subnets provided exist before creating any Azure NetApp Files resources, failing execution if they're missing.

Authentication is made on each operation where we obtain an authorizer to pass to each client we instantiate (in Azure Go SDK for NetAppFiles each resource has its own client). For more information about the authentication process used, see the [Use file-based authentication](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization#use-file-based-authentication) section in [Authentication methods in the Azure SDK for Go](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization).

The last step is the cleanup process (which is not enabled by default; you need to change variable `shouldCleanUp` to `true` at `example.go` file `var()` section to clean up). The process must delete all resources in the reverse order, following the hierarchy; otherwise, we can't remove resources that have nested resources. Before removing the secondary volume, we need to remove the data replication object on it. If there is an error during the application execution, the cleanup process does not take place, and you need to manually perform this task.
The cleanup process uses a function called `WaitForNoANFResource`, while other parts of the code uses `WaitForANFResource`.  Currently, this behavior is required in order to work around the ARM behavior that reports that the object was deleted when in fact its deletion is still in progress (similarly, stating that volume is fully created while the creation is still completing). Also, we will see functions called `GetAnf<resource type>`; these functions were created in this sample to get the name of the resource without its hierarchy represented in the `<resource type>.name` property, which cannot be used directly in other methods of Azure NetApp Files client like `get`.

>Note: see [Resource limits for Azure NetApp Files](https://docs.microsoft.com/en-us/azure/azure-netapp-files/azure-netapp-files-resource-limits) to understand Azure NetApp Files limits.

## Contents

| File/folder                 | Description                                                                                                      |
|-----------------------------|------------------------------------------------------------------------------------------------------------------|
| `media\`                       | Folder that contains screenshots.                                                                                              |
| `netappfiles-go-crr-sdk-sample\`                       | Sample source code folder.                                                                                              |
| `netappfiles-go-crr-sdk-sample\example.go`            | Sample main file.                                                                                                |
| `netappfiles-go-crr-sdk-sample\go.mod`            |The go.mod file defines the module’s module path, which is also the import path used for the root directory, and its dependency requirements, which are the other modules needed for a successful build.|
| `netappfiles-go-crr-sdk-sample\go.sum`            | The go.sum file contains hashes for each of the modules and it's versions used in this sample|
| `netappfiles-go-crr-sdk-sample\internal\`       | Folder that contains all internal packages dedicated to this sample.                |
| `netappfiles-go-crr-sdk-sample\internal\iam\iam.go` | Package that allows us to get the `authorizer` object from Azure Active Directory by using the `NewAuthorizerFromFile` function. |
| `netappfiles-go-crr-sdk-sample\internal\models\models.go`       | Provides models for this sample, e.g. `AzureAuthInfo` models the authorization file.                   |
| `netappfiles-go-crr-sdk-sample\internal\sdkutils\sdkutils.go`       | Contains all functions that directly uses the SDK and some helper functions.                   |
| `netappfiles-go-crr-sdk-sample\internal\uri\uri.go`       | Provides various functions to parse resource IDs and get information or perform validations.                   |
| `netappfiles-go-crr-sdk-sample\internal\utils\utils.go`       | Provides generic functions.                   |
| `.gitignore`                | Define what to ignore at commit time.                                                                            |
| `CHANGELOG.md`              | List of changes to the sample.                                                                                   |
| `CONTRIBUTING.md`           | Guidelines for contributing to the sample.                                                                       |
| `README.md`                 | This README file.                                                                                                |
| `LICENSE`                   | The license for the sample.                                                                                      |
| `CODE_OF_CONDUCT.md`        | Microsoft's Open Source Code of Conduct.                                                                         |

## How to run

1. Go to your GOPATH folder and create the following path: 
    ```powershell
    # PowerShell example
    cd $env:GOPATH/src
    mkdir ./github.com/Azure-Samples
    ```

    ```bash
    # Bash example
    cd $GOPATH/src
    mkdir -p ./github.com/Azure-Samples
    ```
2. Clone the sample locally: 
    ```bash
    cd github.com/Azure-Samples
    git clone https://github.com/Azure-Samples/netappfiles-go-crr-sdk-sample.git
    ```
3. Change folder to **netappfiles-go-crr-sdk-sample/netappfiles-go-crr-sdk-sample**: 
    ```bash
    cd netappfiles-go-crr-sdk-sample/netappfiles-go-crr-sdk-sample
    ```
4. Make sure you have the `azureauth.json` and its environment variable with the path to it defined. See [prerequisites](#Prerequisites).
6. Edit file **example.go** `var()` block and change the variables contents as appropriate (names are self-explanatory).
7. Run the sample: 
    ```bash
    go run .
    ```

Sample output
![e2e execution](./media/e2e-go.png)

## References

* [Cross-region replication of Azure NetApp Files volumes](https://docs.microsoft.com/en-us/azure/azure-netapp-files/cross-region-replication-introduction)
* [Authentication methods in the Azure SDK for Go](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization)
* [Azure SDK for Go Samples](https://github.com/Azure-Samples/azure-sdk-for-go-samples) - contains other resource types samples
* [Resource limits for Azure NetApp Files](https://docs.microsoft.com/en-us/azure/azure-netapp-files/azure-netapp-files-resource-limits)
* [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/cloud-shell/quickstart)
* [Azure NetApp Files documentation](https://docs.microsoft.com/en-us/azure/azure-netapp-files/)
* [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go) 
