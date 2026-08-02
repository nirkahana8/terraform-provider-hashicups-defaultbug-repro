// Minimal provider for reproducing the object-level UseStateForUnknown crash on a
// generated (CustomType-backed) nested attribute. No configuration, no client.
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*usfuProvider)(nil)

type usfuProvider struct{ version string }

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &usfuProvider{version: version} }
}

func (p *usfuProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "usfu"
	resp.Version = p.version
}

func (p *usfuProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *usfuProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *usfuProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewCarResource}
}

func (p *usfuProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
