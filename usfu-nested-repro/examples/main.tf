terraform {
  required_providers {
    usfu = { source = "registry.terraform.io/hashicorp/usfu" }
  }
}
resource "usfu_car" "audi" {
  name  = "A8"
  color = "white"
  # engine = {
  #   kind = "TSFI"
  #   volume = 3
  # }
}
