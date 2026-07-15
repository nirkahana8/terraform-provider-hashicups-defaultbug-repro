terraform {
  required_providers {
    defaultbug = { source = "registry.terraform.io/hashicorp/defaultbug" }
  }
}

# scope is present (via users). locations is OMITTED but carries a classic
# object-level Default of { is_any = true }. locations.is_any is Computed with no
# default of its own.
#
# EXPECTED plan: scope.locations = { is_any = true }
# ACTUAL   plan: scope.locations = { is_any = (known after apply) }
resource "defaultbug_thing" "test" {
  scope = { users = { is_any = true } }
}
