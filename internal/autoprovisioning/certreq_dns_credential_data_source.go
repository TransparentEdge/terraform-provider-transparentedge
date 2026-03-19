package autoprovisioning

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &crDNSCredentialDataSource{}
	_ datasource.DataSourceWithConfigure = &crDNSCredentialDataSource{}
)

// NewCertReqDNSCredentialDataSource is a helper function to simplify the provider implementation.
func NewCertReqDNSCredentialDataSource() datasource.DataSource {
	return &crDNSCredentialDataSource{}
}

// data source implementation.
type crDNSCredentialDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (*crDNSCredentialDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_dns_credential"
}

// Schema defines the schema for the data source.
func (*crDNSCredentialDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "DNS Credential data source.",
		MarkdownDescription: "DNS Credential data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				Description:         "ID of the DNS Credential.",
				MarkdownDescription: "ID of the DNS Credential.",
			},
			"alias": schema.StringAttribute{
				Computed:            true,
				Description:         "Alias of the DNS Credential.",
				MarkdownDescription: "Alias of the DNS Credential.",
			},
			"dns_provider": schema.StringAttribute{
				Computed:            true,
				Description:         "DNS Provider.",
				MarkdownDescription: "DNS Provider.",
			},
			"parameters": schema.MapAttribute{
				Computed:            true,
				Description:         "Keys/parameters of the provider.",
				MarkdownDescription: "Keys/parameters of the provider.",
				Sensitive:           true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *crDNSCredentialDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertReqDNSCredential
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	dnsCredential, err := d.client.GetCRDNSCredential(int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving DNS Credential.",
			err.Error(),
		)

		return
	}

	// Set state
	data.ID = types.Int64Value(int64(dnsCredential.ID))
	data.Alias = types.StringValue(dnsCredential.Alias)

	// Extract the parameters/keys obtained from the API into a map
	keys := make(map[string]attr.Value)

	var dnsProvider string

	for _, key := range dnsCredential.Creds {
		keys[key.KeyName] = types.StringValue(key.KeyValue)
		dnsProvider = key.Provider
	}

	// Transform the map into a Terraform type
	parameters, diags := types.MapValue(types.StringType, keys)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Parameters = parameters
	data.DNSProvider = types.StringValue(dnsProvider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Configure adds the provider configured client to the data source.
func (d *crDNSCredentialDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure", "error while configuring API client")

		return
	}

	d.client = client
}
