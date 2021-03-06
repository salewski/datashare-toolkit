# Google Cloud Platform Release Notes for Marketplace VM Solution

## 0.5.2 Release Notes
* Build - 7/31/2020

### Datashare
* Datashare release version is 0.4.1

### Updates
* Changed the Datashare Github release version to v0.4.1.
* Increased the waiter timeout to 900 seconds from 600 seconds.

### Bug Fixes
* When a user attempted to install Datashare via the Marketplace, it attempted to use the Datashare VM image from the user's project, 
but this image in included in the gcp-financial-services-public project. Modified the URL to point to the gcp-financial-services-public project.

## 0.5.1 Release Notes
* Build - 7/18/2020

### Datashare
* Datashare release version is 0.5.0.

### Updates
* Marketplace VM solution requires the input_adminEmail to be included in the schema even if it is not used. 
  Added this field to the schema to resolve the Marketplace submission error. 


## 0.5 Release Notes
* Build - 7/16/2020

### Datashare
* Datashare release version is 0.5.0.

### Launch Page
* Includes links to Datashare ingestion, UI and the API documentation.
* Includes link to [pricing estimate](https://cloud.google.com/products/calculator?id=31cf51a5-938b-4f9f-8cf3-131eb2b80ba7)

### Deployment Page  
* Custom VM image: **gcp-financial-services-debian-datashare-20200716**
* Default machine type: **small**
* Github Datashare release version: **master**
* Google Cloud Platform Region: **us-central1**
* Google Cloud Platform Storage Bucket Location: **US**
* Use Google RuntimeConfig Waiter to setup Datashare prereqs: **True**
* Deploy the Datashare API to Cloud Run/GKE: **True**
* GKE Zone: **us-central1-a**

### Deployment
The VM solution deployment manager completes the following
* Creates two Storage Buckets
  * install bucket
  * ingestion bucket
* Creates a 3 node GKE cluster (E2 nodes) with Cloud Run enabled
* Creates a small VM and executes startup script 
  * Zips the Datashare ingestion Cloud Function and uploads it to install bucket
  * Enables Google Cloud Platform [APIs](https://github.com/GoogleCloudPlatform/datashare-toolkit/tree/master/marketplace#google-apis-enabled)
  * Creates services accounts with appropriate permissions
  * [Modifies existing service accounts](https://github.com/GoogleCloudPlatform/datashare-toolkit/tree/master/marketplace#new-roles-added-to-service-accounts) to ensure the script can execute successfully
* Uses the [RuntimeConfig beta Waiter](https://cloud.google.com/deployment-manager/runtime-configurator/create-and-delete-runtimeconfig-resources) to wait for the VM startup script to execute before starting subsequent tasks.
* Submits a job to Cloud Build to build the UI container image
  * Deploys the UI container to Cloud Run
* Submits a job to Cloud Build to build the API container image
  * Deploys the API to Cloud Run on GKE

### Deployment Manager Suggested Next Steps page
Instructs user to
* validate the GKE deployment
* map their custom domain to the GKE cluster
* enable TLS on the GKE cluster; allows UI to send https request to the API

### Known Bugs

### Bug Fixes

### Updates
