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
	_ datasource.DataSource              = &crDNSCNAMEVerifDataSource{}
	_ datasource.DataSourceWithConfigure = &crDNSCNAMEVerifDataSource{}
)

// NewCertReqDNSCNAMEVerifDataSource is a helper function to simplify the provider implementation.
func NewCertReqDNSCNAMEVerifDataSource() datasource.DataSource {
	return &crDNSCNAMEVerifDataSource{}
}

// crDNSCNAMEVerifDataSource is the data source implementation.
type crDNSCNAMEVerifDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *crDNSCNAMEVerifDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cr_dns_cname_verification"
}

// Schema defines the schema for the data source.
func (d *crDNSCNAMEVerifDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Certificate Request DNS CNAME Verification",
		MarkdownDescription: "Certificate Request DNS CNAME Verification",
		Attributes: map[string]schema.Attribute{
			"cname": schema.StringAttribute{
				Computed:            true,
				Description:         "The CNAME to configure the _acme-challenge.{domain} record in order to perform a DNS verification by CNAME.",
				MarkdownDescription: "The CNAME to configure the `_acme-challenge.{domain}` record in order to perform a DNS verification by CNAME.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *crDNSCNAMEVerifDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertReqCNAMEVerify
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	cname, err := d.client.GetDNSCNAMEVerification()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve DNS CNAME Verification",
			err.Error(),
		)
		return
	}

	data = CertReqCNAMEVerify{
		CNAME: types.StringValue(cname),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Configure adds the provider configured client to the data source.
func (d *crDNSCNAMEVerifDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
