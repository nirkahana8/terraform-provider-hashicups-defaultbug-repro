package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestNestedObjectDefault_childClobbered shows that an object-level Default on a
// Computed SingleNestedAttribute is discarded for a computed child that has no
// default of its own.
//
// `nested` is omitted from config (null nested attribute) and carries an
// object-level Default of {child = "default-value"}. Per the maintainer behavior
// table in hashicorp/terraform-plugin-framework#726 (Default-on-nested = Yes,
// Default-on-child = No, config = null nested attribute) the plan SHOULD show
// nested = {child = "default-value"} and apply should succeed.
//
// Actual: PlanResourceChange applies the object default (TransformDefaults), then
// MarkComputedNilsAsUnknown re-marks nested.child unknown (it is Computed, null
// in config, and has no default of its own). `terraform plan` shows
// nested = {child = (known after apply)}, and apply fails with:
//
//	Provider returned invalid result object after apply
//	... still indicated an unknown value for defaultbug_thing.test.nested.child.
//
// Run: TF_ACC=1 go test ./... -run TestNestedObjectDefault_childClobbered -v
func TestNestedObjectDefault_childClobbered(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "defaultbug_thing" "test" {}`,
				Check: resource.TestCheckResourceAttr(
					"defaultbug_thing.test", "nested.child", "default-value",
				),
			},
		},
	})
}
