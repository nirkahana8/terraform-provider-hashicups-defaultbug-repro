terraform {
  required_providers {
    defaultbug = { source = "registry.terraform.io/hashicorp/defaultbug" }
  }
}

# `nested` is omitted entirely (null nested attribute). The resource schema sets
# an object-level Default of { child = "default-value" }.
#
# EXPECTED plan: nested = { child = "default-value" }
# ACTUAL   plan: nested = { child = (known after apply) }
resource "defaultbug_thing" "test" {}
