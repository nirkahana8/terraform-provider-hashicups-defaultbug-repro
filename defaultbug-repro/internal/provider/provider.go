// Minimal provider for reproducing the nested-object Default bug in
// hashicorp/terraform-plugin-framework. It has no configuration and no client,
// so the repro is self-contained (no HashiCups server required).
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*bugProvider)(nil)

type bugProvider struct{ version string }

// New returns the provider factory used by both main.go and the acceptance test.
func New(version string) func() provider.Provider {
	return func() provider.Provider { return &bugProvider{version: version} }
}

func (p *bugProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "defaultbug"
	resp.Version = p.version
}

func (p *bugProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *bugProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *bugProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewThingResource}
}

func (p *bugProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
