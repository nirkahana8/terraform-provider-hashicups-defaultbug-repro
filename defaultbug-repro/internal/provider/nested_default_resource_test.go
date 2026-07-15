package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestNestedObjectDefault_childClobbered mirrors the pb_provider case: a
// grandparent `scope` is present (via `users`), the `locations` object is omitted
// from config but carries a classic object-level Default
// (objectdefault.StaticValue({is_any = true})), and locations' child `is_any` is
// Computed with no default of its own.
//
// Per the maintainer behavior table in
// hashicorp/terraform-plugin-framework#726 (Default-on-nested = Yes,
// Default-on-child = No, config = null nested attribute) the plan SHOULD yield
// the object default: scope.locations = {is_any = true}.
//
// Actual: TransformDefaults applies the object default, then
// MarkComputedNilsAsUnknown re-marks scope.locations.is_any unknown (Computed,
// null in config, no default of its own). `terraform plan` shows
// scope.locations = {is_any = (known after apply)}, and apply fails with:
//
//	Provider returned invalid result object after apply
//	... still indicated an unknown value for
//	defaultbug_thing.test.scope.locations.is_any.
//
// Run: TF_ACC=1 go test ./... -run TestNestedObjectDefault_childClobbered -v
func TestNestedObjectDefault_childClobbered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "defaultbug_thing" "test" {
  scope = { users = { is_any = true } }
}`,
				Check: resource.TestCheckResourceAttr(
					"defaultbug_thing.test", "scope.locations.is_any", "true",
				),
			},
		},
	})
}
