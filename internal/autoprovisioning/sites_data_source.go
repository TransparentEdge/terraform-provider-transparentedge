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
	_ datasource.DataSource              = &sitesDataSource{}
	_ datasource.DataSourceWithConfigure = &sitesDataSource{}
)

// NewSitesDataSource is a helper function to simplify the provider implementation.
func NewSitesDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

// sitesDataSource is the data source implementation.
type sitesDataSource struct {
	client *teclient.Client
}

// sitesDataSourceModel maps the data source schema data.
type sitesDataSourceModel struct {
	Sites []sitesModel `tfsdk:"sites"`
}

// sitesModel maps schema data.
type sitesModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Url     types.String `tfsdk:"url"`
	Company types.Int64  `tfsdk:"company"`
	Ssl     types.Bool   `tfsdk:"ssl"`
	Active  types.Bool   `tfsdk:"active"`
}

// Metadata returns the data source type name.
func (d *sitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

// Schema defines the schema for the data source.
func (d *sitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sites": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of all active sites",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "ID of the site",
						},
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "Domain in FDQN form, i.e: 'www.example.com'",
						},
						"company": schema.Int64Attribute{
							Computed:    true,
							Description: "Company ID that owns this domain",
						},
						"ssl": schema.BoolAttribute{
							Computed:    true,
							Description: "If SSL is active (deprecated)",
						},
						"active": schema.BoolAttribute{
							Computed:    true,
							Description: "Active status in the CDN",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sitesDataSourceModel

	sites, err := d.client.GetSites()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Transparent Edge Sites",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, site := range sites {
		siteState := sitesModel{
			ID:      types.Int64Value(int64(site.ID)),
			Url:     types.StringValue(site.Url),
			Company: types.Int64Value(int64(site.Company)),
			Ssl:     types.BoolValue(site.Ssl),
			Active:  types.BoolValue(site.Active),
		}

		state.Sites = append(state.Sites, siteState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *sitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
