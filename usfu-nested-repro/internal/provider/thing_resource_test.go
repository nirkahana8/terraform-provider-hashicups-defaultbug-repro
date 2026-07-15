package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestObjectUseStateForUnknown_crash shows that attaching the stock
// objectplanmodifier.UseStateForUnknown() to a GENERATED, CustomType-backed
// nested attribute crashes at plan time when that attribute is omitted from
// config (null + Computed => unknown in the plan) and its child is itself an
// object. The framework materializes a partial parent object and runs it through
// the generated NestedType.ValueFromObject, which rejects it with:
//
//	Error: Attribute Missing
//	sub is missing from object
//
// Note: it does NOT crash when `nested` is fully specified in config — only when
// it is unknown in the plan. Run:
//
//	TF_ACC=1 go test ./... -run TestObjectUseStateForUnknown_crash -v
func TestObjectUseStateForUnknown_crash(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// nested omitted => unknown in plan => crash.
				Config: `resource "usfu_thing" "test" {}`,
			},
		},
	})
}
