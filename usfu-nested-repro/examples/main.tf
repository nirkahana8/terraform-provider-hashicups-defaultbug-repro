terraform {
  required_providers {
    usfu = { source = "registry.terraform.io/hashicorp/usfu" }
  }
}

# `nested` is a GENERATED (CustomType-backed) SingleNestedAttribute with the stock
# objectplanmodifier.UseStateForUnknown() attached in the resource's Schema().
# Its child `sub` is itself an object. `nested` is omitted here, so it is unknown
# in the plan.
#
# `terraform plan` fails with:
#   Error: Attribute Missing
#   sub is missing from object
#
# (It does NOT crash if you fully specify nested, e.g. nested = { sub = { ref_id = "x" } }.)
resource "usfu_thing" "test" {}
