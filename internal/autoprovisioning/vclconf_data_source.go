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
	_ datasource.DataSource              = &vclconfDataSource{}
	_ datasource.DataSourceWithConfigure = &vclconfDataSource{}
)

// Helper function to simplify the provider implementation.
func NewVclconfDataSource() datasource.DataSource {
	return &vclconfDataSource{}
}

// data source implementation.
type vclconfDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *vclconfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vclconf"
}

// Schema defines the schema for the data source.
func (d *vclconfDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "VCL Configuration listing.",
		MarkdownDescription: "VCL Configuration listing.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the VCL Config.",
				MarkdownDescription: "ID of the VCL Config.",
			},
			"company": schema.Int64Attribute{
				Computed:            true,
				Description:         "Company ID that owns this VCL config.",
				MarkdownDescription: "Company ID that owns this VCL config.",
			},
			"vclcode": schema.StringAttribute{
				Computed:            true,
				Description:         "Verbatim of the VCL code.",
				MarkdownDescription: "Verbatim of the VCL code.",
			},
			"uploaddate": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the configuration was uploaded.",
				MarkdownDescription: "Date when the configuration was uploaded.",
			},
			"productiondate": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the configuration was fully applied in the CDN.",
				MarkdownDescription: "Date when the configuration was fully applied in the CDN.",
			},
			"user": schema.StringAttribute{
				Computed:            true,
				Description:         "User that created the configuration.",
				MarkdownDescription: "User that created the configuration.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *vclconfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state VCLConf

	vclconf, err := d.client.GetActiveVCLConf(teclient.ProdEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(vclconf.ID))
	state.Company = types.Int64Value(int64(vclconf.ID))
	state.VCLCode = types.StringValue(vclconf.VCLCode)
	state.UploadDate = types.StringValue(vclconf.UploadDate)
	state.ProductionDate = types.StringValue(vclconf.ProductionDate)
	state.User = types.StringValue(vclconf.CreatorUser.FirstName + " " + vclconf.CreatorUser.LastName + " <" + vclconf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *vclconfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
