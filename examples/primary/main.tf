terraform {
  required_providers {
    tree = {
      source = "local-development/tree"
    }
  }
}

provider "tree" {}

resource "tree_root" "one" {
  name = "One"
}
