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
	_ datasource.DataSource              = &siteVerifyDataSource{}
	_ datasource.DataSourceWithConfigure = &siteVerifyDataSource{}
)

// NewSiteVerifyDataSource is a helper function to simplify the provider implementation.
func NewSiteVerifyDataSource() datasource.DataSource {
	return &siteVerifyDataSource{}
}

// siteVerifyDataSource is the data source implementation.
type siteVerifyDataSource struct {
	client *teclient.Client
}

// siteVerifyModel maps schema data.
type siteVerifyModel struct {
	Domain              types.String `tfsdk:"domain"`
	VerificantionString types.String `tfsdk:"verification_string"`
}

// Metadata returns the data source type name.
func (d *siteVerifyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_siteverify"
}

// Schema defines the schema for the data source.
func (d *siteVerifyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required: true,
			},
			"verification_string": schema.StringAttribute{
				Computed: true,
				Optional: false,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *siteVerifyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data siteVerifyModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	domain := data.Domain.ValueString()
	verify_string := d.client.GetSiteVerifyString(domain)
	if verify_string == "" {
		resp.Diagnostics.AddError(
			"Unable to retrieve Site Verification string",
			"Could not retrieve the site verification string for the domain: "+domain,
		)
		return
	}

	data = siteVerifyModel{
		Domain:              types.StringValue(domain),
		VerificantionString: types.StringValue(verify_string),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Configure adds the provider configured client to the data source.
func (d *siteVerifyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}