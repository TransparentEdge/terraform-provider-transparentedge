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
	_ datasource.DataSource              = &certReqHTTPDataSource{}
	_ datasource.DataSourceWithConfigure = &certReqHTTPDataSource{}
)

// Helper function to simplify the provider implementation.
func NewCertReqHTTPDataSource() datasource.DataSource {
	return &certReqHTTPDataSource{}
}

// data source implementation.
type certReqHTTPDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *certReqHTTPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_http"
}

// Schema defines the schema for the data source.
func (d *certReqHTTPDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "HTTP Certificate Requests data source.",
		MarkdownDescription: `HTTP Certificate Requests data source.

For detailed documentation (not Terraform-specific), please refer to this [link](https://docs.transparentedge.eu/getting-started/dashboard/auto-provisioning/ssl).`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:            true,
				Description:         "ID of the HTTP Certificate Request.",
				MarkdownDescription: "ID of the HTTP Certificate Request.",
			},
			"domains": schema.SetAttribute{
				Computed:            true,
				Description:         "List of domains for which you want to request a certificate. You can not include wildcard domains, such as `*.example.com`, use DNS Certificate Requests instead.",
				MarkdownDescription: "List of domains for which you want to request a certificate. You can **not** include wildcard domains, such as `*.example.com`, use DNS Certificate Requests instead.",
				ElementType:         types.StringType,
			},
			"certificate_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Certificate associated.",
				MarkdownDescription: "Certificate associated.",
			},
			"standalone": schema.BoolAttribute{
				Computed:            true,
				Description:         "When set to `true`, this indicates that the certificate's domains should be treated as standalone and not merged into an existing certificate, either immediately or during future renewals.",
				MarkdownDescription: "When set to `true`, this indicates that the certificate's domains should be treated as standalone and not merged into an existing certificate, either immediately or during future renewals.",
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
				Description:         "Indicates the current status message for the certificate request. This field will display an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
				MarkdownDescription: "Indicates the current status message for the certificate request. This field will display an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *certReqHTTPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data CertReqHTTP
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	api_model, err := d.client.GetCertReqHTTP(int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving HTTP Certificate Request",
			err.Error(),
		)
		return
	}

	// Generate the list of domains from CommonName and SAN
	sorted_domains := helpers.SplitAndSort(api_model.CommonName + "\n" + api_model.SAN)
	domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	data.ID = types.Int64Value(int64(api_model.ID))
	data.Domains = domains
	data.Standalone = types.BoolValue(api_model.Standalone)
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
func (d *certReqHTTPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
