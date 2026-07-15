package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	gen "terraform-provider-usfu/internal/generated/resource_thing"
)

var _ resource.Resource = (*thingResource)(nil)

func NewThingResource() resource.Resource { return &thingResource{} }

type thingResource struct{}

func (r *thingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_thing"
}

// Schema uses the GENERATED schema (with its strict CustomType, NestedType) and
// attaches the stock objectplanmodifier.UseStateForUnknown() to the nested
// attribute. On plan the framework round-trips the object through
// NestedType.ValueFromObject, which rejects a partial object with
// "child is missing from object".
func (r *thingResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := gen.ThingResourceSchema(ctx)
	nested := s.Attributes["nested"].(schema.SingleNestedAttribute)
	nested.PlanModifiers = append(nested.PlanModifiers, objectplanmodifier.UseStateForUnknown())
	s.Attributes["nested"] = nested
	resp.Schema = s
}

// CRUD is a trivial echo — no client needed.
func (r *thingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data gen.ThingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue("thing-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data gen.ThingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data gen.ThingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
