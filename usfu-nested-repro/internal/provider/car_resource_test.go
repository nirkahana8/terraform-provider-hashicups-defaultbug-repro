package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestObjectUseStateForUnknown_crash shows that attaching the stock
// objectplanmodifier.UseStateForUnknown() to a GENERATED, CustomType-backed
// nested attribute crashes at plan time when that attribute is omitted from
// config on CREATE (null + Computed => unknown in the plan, with no prior state
// to substitute). The framework materializes a partial object and runs it through
// the generated NestedType.ValueFromObject, which rejects it with:
//
//	Error: Attribute Missing
//	child is missing from object
//
// It does NOT crash when `nested` is present in config, nor on update (prior state
// substitutes), nor without the modifier. The child type (scalar or object) is
// irrelevant. Run:
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
