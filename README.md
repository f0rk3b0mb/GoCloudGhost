# GoCloudGhost


**GoCloudGhost** is a modular enumeration tool written in Go, focused on enumerating Cloud resources during security assessments. It currently ony supports azure. There will be more cloud provider support  in the future.

I developed it during a recent cloud penetration testing exercise where i had to do most of these manually.

## ✨ Features

### Azure

- Azure service account token authentication
- Azure access token authentication
- Enumerate blobs from azure storage accounts using access keys
- Retrieve storage account shared keys.
- Enumerate subscription info
- Enumerate storage accounts
- Enumerate resource groups
- Enumerate role assignments  and definitions
- Enumerate keyvaults
- Blob storage enumeration
- Blob storage item download
- Extensible modular architecture — more cloud modules coming soon

### Gcp

- List GCP Compute Instances
- List GCP Storage Buckets
- Enumrate GCP Token Permission
- Extensible modular architecture — more cloud modules coming soon
---

## 📦 Modules

### Azure Management API Enumeration

## Service Account Authentication

```bash
GoCloudGhost azure auth --client-id <appId> --client-secret <secret> --tenant-id <tenantId>
```

### Enumerate Subscription info

#### Run this first to populate env variables

```bash
GoCloudGhost azure management --token <jwt-accesss-key> --subscriptions
```

### Enumerate Key Vaults

```bash
GoCloudGhost azure management --keyvaults
```

### Enumerate Policies

```bash
GoCloudGhost azure management --policies
```

### Enumerate Storage Accounts

```bash
GoCloudGhost azure management --storage
```

### Enumerate Resource Groups

```bash
GoCloudGhost  azure management --groups
```

### Enumerate Role Assignments 

```bash
GoCloudGhost azure management --roles
```

### Blob Storage Enumeration 

```bash
GoCloudGhost azure blob list --account <storage-account-name> --key <shared-key> --container <container-name>
```

### Blob Storage Item Download

```bash
GoCloudGhost azure blob download --account <storage-account-name> --key <shared-key> --container <container-name>

```

### GCP Api Enumeration

### List GCP Compute Instances

```bash
GoCloudGhost gcp list compute --token <oauth-token> --project-id <project-id>

```

### List GCP Storage Buckets

```bash
GoCloudGhost gcp list bucket --token <oauth-token> --project-id <project-id>

```


### Enumrate GCP Token Permission

```bash
GoCloudGhost gcp enum --token <ouath-token>  --project-id <project-id>

```



## 🔧 Installation

You can compile using the steps below or download the pre-compiled binay from [releases page](https://github.com/f0rk3b0mb/GoCloudGhost/releases)


- Clone the repository:

Youll need to have the golang compiler installed and in your PATH

Install compiler from [https://go.dev/doc/install](https://go.dev/doc/install)

``` bash
git clone https://github.com/f0rk3b0mb/GoCloudGhost.git
cd GoCloudGhost
```
- Build the binary:

```bash
go build .
```

## 🛡️ Disclaimer
This tool is intended for authorized testing and educational purposes only. Ensure you have proper permissions before using it in any environment.

## 📌 Roadmap

✅ Blob Storage Enumeration

✅ Azure Management API enumeration

✅ Gcp Token Permission

✅ Gcp Compute and Stoage bucker enumeration

Gcp Sercive Account Token Impersonate

AKS, App Services, Key Vault discovery

✅ Azure role/permission auditing

Support for other cloud providers



## 🗒️ Licences

http://www.apache.org/licenses/

## 📫 Contact
Feel free to contribute or raise issues.

Reach out @f0rk3b0mb on twitter