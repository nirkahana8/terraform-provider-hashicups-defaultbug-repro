package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*thingResource)(nil)

// NewThingResource mirrors the pb_provider shape faithfully:
//
//	scope     (Optional+Computed, no default)
//	  users     (Optional+Computed) { is_any bool, no default }  -- set in config
//	  locations (Optional+Computed, OBJECT DEFAULT {is_any:true}) -- OMITTED in config
//	    is_any  (bool, Optional+Computed, NO default of its own)
//
// The "default" lives on the locations OBJECT (like objectdefault.StaticValue),
// and locations' child is_any is Computed with no default of its own.
func NewThingResource() resource.Resource { return &thingResource{} }

type thingResource struct{}

type thingModel struct {
	ID    types.String `tfsdk:"id"`
	Scope types.Object `tfsdk:"scope"`
}

var (
	usersAttrTypes     = map[string]attr.Type{"is_any": types.BoolType}
	locationsAttrTypes = map[string]attr.Type{"is_any": types.BoolType}
)

func (r *thingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_thing"
}

func (r *thingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"scope": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"is_any": schema.BoolAttribute{Optional: true, Computed: true},
						},
					},
					"locations": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						// OBJECT-level default, fully populated.
						Default: objectdefault.StaticValue(types.ObjectValueMust(
							locationsAttrTypes,
							map[string]attr.Value{"is_any": types.BoolValue(true)},
						)),
						Attributes: map[string]schema.Attribute{
							"is_any": schema.BoolAttribute{
								Optional: true,
								Computed: true,
								// NO default of its own.
							},
						},
					},
				},
			},
		},
	}
}

// CRUD is a trivial echo — no client needed.
func (r *thingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data thingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue("thing-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data thingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data thingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *thingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
