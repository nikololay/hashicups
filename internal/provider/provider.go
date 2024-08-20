package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &hashicupsProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hashicupsProvider{
			version: version,
		}
	}
}

type hashicupsProvider struct {
	version string
}

// Metadata implements provider.Provider.
func (h *hashicupsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hashicups"
	resp.Version = h.version
}

// Schema implements provider.Provider.
func (h *hashicupsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure implements provider.Provider.
func (h *hashicupsProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

// DataSources implements provider.Provider.
func (h *hashicupsProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// Resources implements provider.Provider.
func (h *hashicupsProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}
