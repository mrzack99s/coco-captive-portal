# CoCo Captive Portal

CoCo Captive Portal is a network security application. it's extremely simple, lightning fast installer.

## Download

Please check the [releases](https://github.com/mrzack99s/coco-captive-portal/releases/) page.

## Requirements

- NIC: 2 interfaces
- CPU: Minumum require 2 physical cores (4 vCores)
- Memory: Minumum require 2 GiB
- OS
  - Ubuntu: 18.04, 20.04
  - Debian: 10, 11
- IP Address of an interface and prepare routing table

## Authentication protocol supports

- LDAP
  - Google LDAP (G-Suite)
  - Active Directory
- Radius

## Getting Started

- Download CoCo Installer

```bash
curl -L https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco-installer -o coco-installer
```

- Grant it to be executable

```bash
sudo chmod +x coco-installer
```

- Install CoCo Captive Portal

```bash
sudo ./coco-installer up
```

> Install with lastest version

```bash
sudo ./coco-installer up --latest
```

> Install with ignore some verification

```bash
sudo ./coco-installer up --ignore
```

> Uninstall CoCo Captive Portal

```bash
sudo ./coco-installer purge
```

## Certificates

- Need a certificates in **_certs_** directory
  - For Auth Endpoint
    - **authfullchain.pem**
    - **authprivkey.pem**
  - For Auth Endpoint
    - **operatorfullchain.pem**
    - **operatorprivkey.pem**
  - For LDAP (Optional)
    - **ldapchain.pem**
    - **ldapprivkey.pem**
- to auto generate self-signed certificate

```bash
/var/coco-captive-portal/coco gencert
```

## API

- API Token will store in file **_app_credentials.yaml_**
- For renew api token

```bash
/var/coco-captive-portal/coco renew-api-token
```

- API Documents, open url to your endpoint then **/docs**
  - Ex. https://coco-captive-portal.local/docs

## Manual Install

### Devtool Requirements

- Make
- gcc - build essential
- Libpcap
- Redis 7
- Golang >= 1.17.xx
- NodeJS >= 16.xx
- Yarn

```bash
git clone https://github.com/mrzack99s/coco-captive-portal

cd coco-captive-portal

# Install devtools and libs
sudo make install-dev-tools

# Build CoCo Captive Portal
sudo make build

# Build UI
sudo make build-node

# Create directory
mkdir -p /var/coco-captive-portal

# Copy
sudo cp -r dist-auth-ui /var/coco-captive-portal
sudo cp -r dist-operator-ui /var/coco-captive-portal
sudo cp coco /var/coco-captive-portal
sudo cp coco-captive-portal.service /etc/systemd/system

# Reload service
sudo systemctl daemon-reload
sudo systemctl enable --now redis-server
sudo systemctl enable --now coco-captive-portal
```

## License

Copyright (c) 2022 [CoCo Captive Portal](https://github.com/mrzack99s/coco-captive-portal)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
