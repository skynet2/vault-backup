# Vault-Backup 

![build workflow](https://github.com/skynet2/vault-backup/actions/workflows/release.yaml/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/skynet2/vault-backup/branch/master/graph/badge.svg?token=LAARF8BFLO)](https://codecov.io/gh/skynet2/vault-backup)
[![go-report](https://goreportcard.com/badge/github.com/skynet2/vault-backup?nocache=true)](https://goreportcard.com/report/github.com/skynet2/vault-backup)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/skynet2/vault-backup)](https://pkg.go.dev/github.com/skynet2/vault-backup?tab=doc)

## Description
Vault-Backup is an automatic backup solution for HashiCorp Vault. The project is designed to run as a Docker container or a standalone binary and supports the backup of Vault data to remote destinations such as AWS S3. The project also has Prometheus metrics support, which can be enabled using an environment variable.

## Features

- Automated HashiCorp Vault backups
- Supports AWS S3 as a backup destination
- Prometheus metrics support with counters for successful and failed backup attempts
- Automatically reads and applies environment variables from JSON files in a specified path when using the HashiCorp Vault Agent Sidecar Injector (see `VAULT_SECRETS_PATH`)
## Prerequisites

- A running HashiCorp Vault instance
- Access to an S3-compatible storage service (e.g., Amazon S3)
- (Optional) A running Prometheus instance for metrics collection

## Installation

### Docker

1. Pull the Docker image: `docker pull skydev/vault-backup`
2. Configure the necessary environment variables (see below for a list of supported variables).
3. Run the Docker container with the configured environment variables: `docker run --env-file <your-env-file> skydev/vault-backup`

[//]: # (### Binary)

[//]: # (1. Download the latest binary from the [GitHub releases page]&#40;https://github.com/skynet2/vault-backup/releases&#41;.)

[//]: # (2. Make the binary executable: `chmod +x vault-backup`)

[//]: # (3. Configure the necessary environment variables &#40;see below for a list of supported variables&#41;.)

[//]: # (4. Run the binary with the configured environment variables: `./vault-backup`)

## Usage

The following environment variables are used to configure the behavior of the Vault-Backup application:

- `VAULT_NAME`: Customize the name of your Vault instance. This is used primarily for Prometheus metrics.
- `PROMETHEUS_PUSH_GATEWAY_URL`: The URL to the Prometheus Pushgateway. This is optional and should only be set if you want to enable Prometheus metrics.
- `VAULT_URLS`: A comma-separated list of Vault URLs, e.g., `https://vault1:8200,https://vault2:8200`. The application will try to connect using these URLs.
- `VAULT_TOKEN`: A Vault token with the following access policy applied:
```hcl
path "sys/storage/raft/snapshot" {
    capabilities = ["read"]
}
```
- `VAULT_SECRETS_PATH`: Path to HashiCorp Vault files. This is useful when using the HashiCorp Vault Agent Sidecar Injector. For more information, visit [Vault Agent Sidecar Injector documentation](https://developer.hashicorp.com/vault/docs/platform/k8s/injector). Files are expected to be in JSON format. The application will read the JSON files and apply the environment variables. Expected format: `{"key": "value"}`. (On linux system by default /vault/secrets/). **Disabled by default**.
- S3 Configuration:
    - `S3_ACCESS_KEY`: Your S3 access key.
    - `S3_SECRET_KEY`: Your S3 secret key.
    - `S3_ENDPOINT`: The S3 endpoint URL.
    - `S3_REGION`: The S3 region.
    - `S3_DISABLE_SSL`: Disable SSL for S3 connections (true or false).
    - `S3_BUCKET`: The name of the S3 bucket to store the backups.

## Prometheus Metrics

The following Prometheus metrics are available when the `PROMETHEUS_PUSH_GATEWAY_URL` environment variable is set:

- `vault_backup_success_total`: A counter for the total number of successful backups.
- `vault_backup_errors_total`: A counter for the total number of failed backups. This metric includes a label `err_text` with the error text.

## Contributing

Contributions to the project are welcome. Please submit a pull request or create an issue to propose new features, report bugs, or suggest improvements.

## License

This project is licensed under the [MIT License](LICENSE).

## References

For more information on setting up Vault backups, refer to this guide: [DIY Vault Backup](https://shadow-soft.com/diy-vault-backup/)