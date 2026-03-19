package autoprovisioning

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certReqHTTPDataSource{}
	_ datasource.DataSourceWithConfigure = &certReqHTTPDataSource{}
)

// NewCertReqHTTPDataSource is a helper function to simplify the provider implementation.
func NewCertReqHTTPDataSource() datasource.DataSource {
	return &certReqHTTPDataSource{}
}

// data source implementation.
type certReqHTTPDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (*certReqHTTPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_http"
}

// Schema defines the schema for the data source.
func (*certReqHTTPDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

	apiModel, err := d.client.GetCertReqHTTP(int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving HTTP Certificate Request",
			err.Error(),
		)

		return
	}

	// Generate the list of domains from CommonName and SAN
	sortedDomains := helpers.SplitAndSort(apiModel.CommonName + "\n" + apiModel.SAN)
	domains, diags := types.SetValueFrom(ctx, types.StringType, sortedDomains)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	data.ID = types.Int64Value(int64(apiModel.ID))
	data.Domains = domains
	data.Standalone = types.BoolValue(apiModel.Standalone)
	data.CreatedAt = types.StringValue(apiModel.CreatedAt)
	data.UpdatedAt = types.StringValue(apiModel.UpdatedAt)

	if apiModel.CertificateID == nil {
		data.CertificateID = types.Int64Null()
	} else {
		data.CertificateID = types.Int64Value(int64(*apiModel.CertificateID))
	}

	if apiModel.Log == nil {
		data.StatusMessage = types.StringNull()
	} else {
		data.StatusMessage = types.StringValue(helpers.ParseCertReqLogString(*apiModel.Log))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Configure adds the provider configured client to the data source.
func (d *certReqHTTPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
