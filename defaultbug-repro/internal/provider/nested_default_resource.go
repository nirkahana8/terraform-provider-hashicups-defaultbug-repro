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

// NewThingResource is the repro resource: it exposes a Computed
// SingleNestedAttribute `nested` that carries an object-level Default. The
// nested object's child is Computed but has NO default of its own.
func NewThingResource() resource.Resource { return &thingResource{} }

type thingResource struct{}

type thingModel struct {
	ID     types.String `tfsdk:"id"`
	Nested types.Object `tfsdk:"nested"`
}

var nestedAttrTypes = map[string]attr.Type{"child": types.StringType}

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
			"nested": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				// Object-level Default, fully populated. Per the maintainer
				// behavior table in terraform-plugin-framework#726
				// (Default-on-nested = Yes, Default-on-child = No, config = null
				// nested attribute) this row is documented as yielding the
				// "single nested attribute default".
				Default: objectdefault.StaticValue(types.ObjectValueMust(
					nestedAttrTypes,
					map[string]attr.Value{"child": types.StringValue("default-value")},
				)),
				Attributes: map[string]schema.Attribute{
					"child": schema.StringAttribute{
						Optional: true,
						Computed: true,
						// Intentionally NO Default here — this is the crux: the
						// child is Computed with no default of its own.
					},
				},
			},
		},
	}
}

// CRUD is a trivial echo — no client needed. State is set from the plan.
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
