# Terraform Provider Etcd

This project has been based on https://github.com/hashicorp/terraform-provider-hashicups

Its allow you to manage some etcd elements via terraform.


## What This provider can do

### User

TODO

### Role

TODO

### Permission

TODO

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-etcd
````

## Test sample configuration

First, build and install the provider.

```shell
$ move Makefile.tmpl Makefile && make install
```

Then, go to examples directory.

```shell
$ cd examples
```

Create your main.tf using  main.tf.tmpl as reference and run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

### Self-signed SSL connection

If you are using a non trusted root CA, You should use SSL_CERT_DIR or SSL_CERT_FILE to choose a custom trusted bundle file.

