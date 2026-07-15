terraform {
  required_providers {
    usfu = { source = "registry.terraform.io/hashicorp/usfu" }
  }
}

# `nested` is a GENERATED (CustomType-backed) SingleNestedAttribute with the stock
# objectplanmodifier.UseStateForUnknown() attached in the resource's Schema().
# It is omitted here, so on CREATE it is unknown in the plan with no prior state.
#
# `terraform plan` fails with:
#   Error: Attribute Missing
#   child is missing from object
#
# It does NOT crash if you specify nested (e.g. nested = { child = "x" }), nor on
# update once state exists. The child type (scalar or object) is irrelevant.
resource "usfu_thing" "test" {}
