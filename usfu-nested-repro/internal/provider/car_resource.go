package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	gen "terraform-provider-usfu/internal/generated/resource_car"
)

var _ resource.Resource = (*carResource)(nil)

func NewCarResource() resource.Resource { return &carResource{} }

type carResource struct{}

func (r *carResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_car"
}

// Schema uses the GENERATED schema (with its strict CustomType, NestedType) and
// attaches the stock objectplanmodifier.UseStateForUnknown() to the nested
// attribute. When `nested` is unknown in the plan with no prior state to
// substitute (i.e. create with `nested` omitted), the framework materializes a
// partial object and runs it through NestedType.ValueFromObject, which rejects it
// with "child is missing from object". The child type (scalar or object) is
// irrelevant; removing the modifier removes the crash.
func (r *carResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"color": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"engine": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"kind": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"volume": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
				},
				CustomType: gen.EngineType{
					ObjectType: types.ObjectType{
						AttrTypes: gen.EngineValue{}.AttributeTypes(ctx),
					},
				},
				// PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
				Optional: true,
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
	resp.Schema = s
}

// CRUD is a trivial echo — no client needed.
func (r *carResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data gen.CarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = types.StringValue("thing-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data gen.CarModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data gen.CarModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
