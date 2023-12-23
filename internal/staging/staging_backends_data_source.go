package staging

import (
	"context"
	"strconv"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &stagingBackendsDataSource{}
	_ datasource.DataSourceWithConfigure = &stagingBackendsDataSource{}
)

// Helper function to simplify the provider implementation.
func NewStagingBackendsDataSource() datasource.DataSource {
	return &stagingBackendsDataSource{}
}

// data source implementation.
type stagingBackendsDataSource struct {
	client *teclient.Client
}

// maps the data source schema data.
type stagingBackendsDataSourceModel struct {
	StagingBackends []StagingBackend `tfsdk:"staging_backends"`
}

// Metadata returns the data source type name.
func (d *stagingBackendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_backends"
}

// Schema defines the schema for the data source.
func (d *stagingBackendsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Staging backend listing.",
		MarkdownDescription: "Staging backend listing.",

		Attributes: map[string]schema.Attribute{
			"staging_backends": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of all staging backends.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "ID of the staging backend.",
							MarkdownDescription: "ID of the staging backend.",
						},
						"company": schema.Int64Attribute{
							Computed:            true,
							Description:         "Company ID that owns this staging backend.",
							MarkdownDescription: "Company ID that owns this staging backend.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the staging backend.",
							MarkdownDescription: "Name of the staging backend.",
						},
						"vclname": schema.StringAttribute{
							Computed:            true,
							Description:         "Final unique name of the backend to be referenced in VCL Code: 'c{company_id}_{name}'.",
							MarkdownDescription: "Final unique name of the backend to be referenced in VCL Code: `c{company_id}_{name}`.",
						},
						"origin": schema.StringAttribute{
							Computed:            true,
							Description:         "IP or DNS name pointing to the origin backend, for example: 'my-origin.com'.",
							MarkdownDescription: "IP or DNS name pointing to the origin backend, for example: `my-origin.com`.",
						},
						"ssl": schema.BoolAttribute{
							Computed:            true,
							Description:         "Use TLS encription when contacting with the origin backend.",
							MarkdownDescription: "Use TLS encription when contacting with the origin backend.",
						},
						"port": schema.Int64Attribute{
							Computed:            true,
							Description:         "Port where the origin is listening to HTTP requests, for example: 80 or 443.",
							MarkdownDescription: "Port where the origin is listening to HTTP requests, for example: `80` or `443`.",
						},
						"hchost": schema.StringAttribute{
							Computed:    true,
							Description: "Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com.",
						},
						"hcpath": schema.StringAttribute{
							Computed:            true,
							Description:         "Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com.",
							MarkdownDescription: "Host header that the healthcheck probe will send to the origin, for example: `www.my-origin.com`.",
						},
						"hcstatuscode": schema.Int64Attribute{
							Computed:            true,
							Description:         "Status code expected when the probe receives the HTTP healthcheck response, for example: 200.",
							MarkdownDescription: "Status code expected when the probe receives the HTTP healthcheck response, for example: `200`.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *stagingBackendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state stagingBackendsDataSourceModel

	stagingBackends, err := d.client.GetBackends(teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Staging Backends info",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, stagingBackend := range stagingBackends {
		stagingBackendState := StagingBackend{
			ID:           types.Int64Value(int64(stagingBackend.ID)),
			Company:      types.Int64Value(int64(stagingBackend.Company)),
			Name:         types.StringValue(stagingBackend.Name),
			VclName:      types.StringValue("c" + strconv.Itoa(stagingBackend.Company) + "_" + stagingBackend.Name),
			Origin:       types.StringValue(stagingBackend.Origin),
			Ssl:          types.BoolValue(stagingBackend.Ssl),
			Port:         types.Int64Value(int64(stagingBackend.Port)),
			HCHost:       types.StringValue(stagingBackend.HCHost),
			HCPath:       types.StringValue(stagingBackend.HCPath),
			HCStatusCode: types.Int64Value(int64(stagingBackend.HCStatusCode)),
		}

		state.StagingBackends = append(state.StagingBackends, stagingBackendState)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *stagingBackendsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
