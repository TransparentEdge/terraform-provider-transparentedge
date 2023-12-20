package autoprovisioning

import (
	"context"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &crDNSProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &crDNSProviderDataSource{}
)

// NewCRDNSProvidersDataSource is a helper function to simplify the provider implementation.
func NewCRDNSProvidersDataSource() datasource.DataSource {
	return &crDNSProviderDataSource{}
}

// crDNSProviderDataSource is the data source implementation.
type crDNSProviderDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *crDNSProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cr_dns_providers"
}

// Schema defines the schema for the data source.
func (d *crDNSProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List of the available Certificate Request DNS Providers.",
		MarkdownDescription: "List of the available Certificate Request DNS Providers.",
		Attributes: map[string]schema.Attribute{
			"providers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Available DNS providers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "ID of the provider",
							MarkdownDescription: "ID of the provider",
						},
						"provider": schema.StringAttribute{
							Computed:            true,
							Description:         "DNS Provider",
							MarkdownDescription: "DNS Provider",
						},

						"keys": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Keys/parameters of the provider",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key_name_id": schema.Int64Attribute{
										Computed:            true,
										Description:         "ID of the key",
										MarkdownDescription: "ID of the key",
									},
									"key_name": schema.StringAttribute{
										Computed:            true,
										Description:         "Name of the key / parameter required by the provider",
										MarkdownDescription: "Name of the key / parameter required by the provider",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *crDNSProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var providers CertReqDNSProviders
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &providers)...)

	resp_providers, err := d.client.GetCRDNSProviders()
	if err != nil || resp_providers == nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve DNS Providers",
			"Unable to retrieve DNS Providers",
		)
		return
	}

	// Map response body to model
	for _, prov := range resp_providers {
		var keys []CertReqDNSKeys
		for _, key := range prov.Keys {
			keys = append(keys, CertReqDNSKeys{
				KeyNameID: types.Int64Value(int64(key.KeyNameID)),
				KeyName:   types.StringValue(key.KeyName)})
		}
		state := CertReqDNSProvider{
			ID:       types.Int64Value(int64(prov.ID)),
			Provider: types.StringValue(prov.Provider),
			Keys:     keys,
		}

		providers.Providers = append(providers.Providers, state)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &providers)...)
}

// Configure adds the provider configured client to the data source.
func (d *crDNSProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
