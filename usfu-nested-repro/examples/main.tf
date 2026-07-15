terraform {
  required_providers {
    usfu = { source = "registry.terraform.io/hashicorp/usfu" }
  }
}
resource "usfu_thing" "test" {}
