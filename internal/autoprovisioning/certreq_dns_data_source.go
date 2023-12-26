package autoprovisioning

import (
	"context"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certReqDNSDataSource{}
	_ datasource.DataSourceWithConfigure = &certReqDNSDataSource{}
)

// Helper function to simplify the provider implementation.
func NewCertReqDNSDataSource() datasource.DataSource {
	return &certReqDNSDataSource{}
}

// data source implementation.
type certReqDNSDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *certReqDNSDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_dns"
}

// Schema defines the schema for the data source.
func (d *certReqDNSDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "DNS Certificate Requests data source.",
		MarkdownDescription: `DNS Certificate Requests data source.

For detailed documentation (not Terraform-specific), please refer to this [link](https://docs.transparentedge.eu/getting-started/dashboard/auto-provisioning/ssl).`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				Description:         "ID of the DNS Certificate Request.",
				MarkdownDescription: "ID of the DNS Certificate Request.",
			},
			"domains": schema.SetAttribute{
				Computed:            true,
				Description:         "List of domains for which you want to request a certificate. You can include wildcard domains, such as `*.example.com`, to cover subdomains under a common domain.",
				MarkdownDescription: "List of domains for which you want to request a certificate. You can include wildcard domains, such as `*.example.com`, to cover subdomains under a common domain.",
				ElementType:         types.StringType,
			},
			"credential": schema.Int64Attribute{
				Computed:            true,
				Description:         "DNS Credential associated.",
				MarkdownDescription: "DNS Credential associated.",
			},
			"certificate_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
				MarkdownDescription: "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date of creation.",
				MarkdownDescription: "Date of creation.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date of last update.",
				MarkdownDescription: "Date of last update.",
			},
			"status_message": schema.StringAttribute{
				Computed:            true,
				Description:         "Indicates the current status message for the certificate request. This field will display a success message if the certificate is obtained successfully or an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
				MarkdownDescription: "Indicates the current status message for the certificate request. This field will display a success message if the certificate is obtained successfully or an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *certReqDNSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data CertReqDNS
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	api_model, err := d.client.GetCertReqDNS(int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving DNS Certificate Request",
			err.Error(),
		)
		return
	}

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.Domains)
	domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	data.ID = types.Int64Value(int64(api_model.ID))
	data.Domains = domains
	data.Credential = types.Int64Value(int64(api_model.Credential))
	data.CreatedAt = types.StringValue(api_model.CreatedAt)
	data.UpdatedAt = types.StringValue(api_model.UpdatedAt)
	if api_model.CertificateID == nil {
		data.CertificateID = types.Int64Null()
	} else {
		data.CertificateID = types.Int64Value(int64(*api_model.CertificateID))
	}
	if api_model.Log == nil {
		data.StatusMessage = types.StringNull()
	} else {
		data.StatusMessage = types.StringValue(helpers.ParseCertReqLogString(*api_model.Log))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Configure adds the provider configured client to the data source.
func (d *certReqDNSDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
