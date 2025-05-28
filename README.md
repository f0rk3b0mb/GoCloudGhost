# GoCloudGhost


**GoCloudGhost** is a modular enumeration tool written in Go, focused on enumerating Cloud resources during security assessments. It currently ony supports azure. There will be more cloud provider support  in the future.

I developed it during a recent cloud penetration testing exercise where i had to do most of these manually.

## ✨ Features

- Enumerate blobs from Azure Storage accounts using access keys
- Enumerate Azure Management API resources using access tokens
- Retrieve storage account shared keys.
- List blob container contents with shared keys.
- Extensible modular architecture — more cloud modules coming soon

---

## 📦 Modules

### Management API Enumeration

### Enumerate Subscription info

```bash
GoCloudGhost management --token <jwt-accesss-key> --subscriptions
```

### Enumerate Storage Accounts

```bash
GoCloudGhost management --subscription <subscription-id> --token <jwt-accesss-key>
```

### Enumerate Resource Groups

```bash
GoCloudGhost management --subscription <sunscription-id> --token <jwt-accesss-key> --groups
```

### Enumerate Role Assignments 

```bash
GoCloudGhost management --subscription <sunscription-id> --token <jwt-accesss-key> --roles
```

### Blob Storage Enumeration

```bash
GoCloudGhost blob list --account <storage-account-name> --key <shared-key> --container <container-name>
```

### Blob Storage Item Download

```bash
GoCloudGhost blob download --account <storage-account-name> --key <shared-key> --container <container-name>

```

## 🔧 Installation
Clone the repository:

``` bash
git clone https://github.com/f0rk3b0mb/GoCloudGhost.git
cd GoCloudGhost
```
Build the binary:

```bash
go build -o GoCloudGhost
```

## 🛡️ Disclaimer
This tool is intended for authorized testing and educational purposes only. Ensure you have proper permissions before using it in any environment.

## 📌 Roadmap

✅ Blob Storage Enumeration

✅ Azure Management API enumeration

 AKS, App Services, Key Vault discovery

 Azure role/permission auditing

 Support for other cloud providers



## 🗒️ Licences

http://www.apache.org/licenses/

## 📫 Contact
Feel free to contribute or raise issues.

Reach out @f0rk3b0mb on twitter