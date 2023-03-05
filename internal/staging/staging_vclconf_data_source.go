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
	_ datasource.DataSource              = &stagingVclConfDataSource{}
	_ datasource.DataSourceWithConfigure = &stagingVclConfDataSource{}
)

// Helper function to simplify the provider implementation.
func NewStagingVclconfDataSource() datasource.DataSource {
	return &stagingVclConfDataSource{}
}

// data source implementation.
type stagingVclConfDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *stagingVclConfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_vclconf"
}

// Schema defines the schema for the data source.
func (d *stagingVclConfDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
func (d *stagingVclConfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StagingVCLConf

	stagingVclConf, err := d.client.GetActiveVCLConf(teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Staging VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(stagingVclConf.ID))
	state.Company = types.Int64Value(int64(stagingVclConf.ID))
	state.VCLCode = types.StringValue(stagingVclConf.VCLCode)
	state.UploadDate = types.StringValue(stagingVclConf.UploadDate)
	state.ProductionDate = types.StringValue(stagingVclConf.ProductionDate)
	state.User = types.StringValue(stagingVclConf.CreatorUser.FirstName + " " + stagingVclConf.CreatorUser.LastName + " <" + stagingVclConf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *stagingVclConfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
