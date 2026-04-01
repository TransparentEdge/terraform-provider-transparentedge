package autoprovisioning

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/customtypes"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vclconfDataSource{}
	_ datasource.DataSourceWithConfigure = &vclconfDataSource{}
)

// NewVclconfDataSource is a helper function to simplify the provider implementation.
func NewVclconfDataSource() datasource.DataSource {
	return &vclconfDataSource{}
}

// data source implementation.
type vclconfDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (*vclconfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vclconf"
}

// Schema defines the schema for the data source.
func (*vclconfDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				CustomType:          customtypes.VCLCodeType{},
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
			"comment": schema.StringAttribute{
				Computed:            true,
				Description:         "Optional comment describing the changes introduced by this configuration.",
				MarkdownDescription: "Optional comment describing the changes introduced by this configuration.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *vclconfDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state VCLConf

	apiResp, err := d.client.GetActiveVCLConf(teclient.ProdEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read VclConf info",
			err.Error(),
		)

		return
	}

	state.ID = types.Int64Value(int64(apiResp.ID))
	state.Company = types.Int64Value(int64(apiResp.ID))
	state.VCLCode = customtypes.NewVCLCodeValue(apiResp.VCLCode)
	state.UploadDate = types.StringValue(apiResp.UploadDate)
	state.ProductionDate = types.StringValue(apiResp.ProductionDate)
	state.User = types.StringValue(apiResp.CreatorUser.FirstName + " " + apiResp.CreatorUser.LastName + " <" + apiResp.CreatorUser.Email + ">")
	state.Comment = types.StringValue(apiResp.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *vclconfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
