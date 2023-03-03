package staging

import (
	"context"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &stagingvclconfDataSource{}
	_ datasource.DataSourceWithConfigure = &stagingvclconfDataSource{}
)

// Helper function to simplify the provider implementation.
func NewStagingVclconfDataSource() datasource.DataSource {
	return &stagingvclconfDataSource{}
}

// data source implementation.
type stagingvclconfDataSource struct {
	client *teclient.Client
}

// maps the data source schema data.
type stagingvclconfDataSourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Company        types.Int64  `tfsdk:"company"`
	VCLCode        types.String `tfsdk:"vclcode"`
	UploadDate     types.String `tfsdk:"uploaddate"`
	ProductionDate types.String `tfsdk:"productiondate"`
	User           types.String `tfsdk:"user"`
}

// Metadata returns the data source type name.
func (d *stagingvclconfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stagingvclconf"
}

// Schema defines the schema for the data source.
func (d *stagingvclconfDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Staging VCL Config",
			},
			"company": schema.Int64Attribute{
				Computed:    true,
				Description: "Company ID that owns this Staging VCL config",
			},
			"vclcode": schema.StringAttribute{
				Computed:    true,
				Description: "Verbatim of the VCL code",
			},
			"uploaddate": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the configuration was uploaded",
			},
			"productiondate": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the configuration was fully applied in the CDN",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "User that created the configuration",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *stagingvclconfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state stagingvclconfDataSourceModel

	stagingvclconf, err := d.client.GetActiveVCLConf(teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Staging VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(stagingvclconf.ID))
	state.Company = types.Int64Value(int64(stagingvclconf.ID))
	state.VCLCode = types.StringValue(stagingvclconf.VCLCode)
	state.UploadDate = types.StringValue(stagingvclconf.UploadDate)
	state.ProductionDate = types.StringValue(stagingvclconf.ProductionDate)
	state.User = types.StringValue(stagingvclconf.CreatorUser.FirstName + " " + stagingvclconf.CreatorUser.LastName + " <" + stagingvclconf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *stagingvclconfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
