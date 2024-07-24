package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &AppvantageProvider{}
var _ provider.ProviderWithFunctions = &AppvantageProvider{}

type AppvantageProvider struct {
	version string
}

type AppvantageProviderModel struct {
}

func (p *AppvantageProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "appvantage"
	resp.Version = p.version
}

func (p *AppvantageProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *AppvantageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AppvantageProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *AppvantageProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *AppvantageProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AppvantageProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewSesSmtpPasswordV4Function,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AppvantageProvider{
			version: version,
		}
	}
}
