# Terraform Provider: Tree Example

This provider is a proof of concept and working example of the helper packages in [spilliams/tree-terraform-provider](https://github.com/spilliams/tree-terraform-provider).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `bin` directory.

To use the locally-built provider in a terraform configuration, modify your `~/.terraformrc` file to include ths following:

```hcl
provider_installation {
  dev_overrides {
    "local-development/tree" = "/path/to/clone/terraform-provider-tree-example/bin"
  }

  direct {}
}
```

Then your terraform configuration can use this code:

```hcl
terraform {
  required_providers {
    tree = {
      source = "local-development/tree"
    }
  }
}

provider "tree" {}
```

## Releasing the provider

### Update documentation

To generate or update documentation, run `make docs`.

### Run full test suite

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Goreleaser

```sh
goreleaser release
```

## Using the provider in production

### Managing the DynamoDB table

If you want to explicitly manage the DynamoDB table, you can do so with the module:

```hcl
provider "aws" {}

module "this" {
  source = "git@github.com:spilliams/tree-terraform-provider.git//terraform/module?ref=v0.2.0"

  enable_table_guardian = false
  kms_key_arn           = "foo"
  table_name            = "bar"
}
```
