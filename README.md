# Terraform Provider Etcd

This project has been based on https://github.com/hashicorp/terraform-provider-hashicups

Its allow you to manage some etcd elements via terraform.


## How to use it?

Please check [/docs](https://github.com/cropalato/terraform-provider-etcd/tree/main/docs) directory.


## Self-signed SSL connection

If you are using a non trusted root CA, You should use SSL_CERT_DIR or SSL_CERT_FILE to choose a custom trusted bundle file.

Adding a dummy line to remove
