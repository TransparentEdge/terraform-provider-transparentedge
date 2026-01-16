package staging

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &stagingBackendDataSource{}
	_ datasource.DataSourceWithConfigure = &stagingBackendDataSource{}
)

// Helper function to simplify the provider implementation.
func NewStagingBackendDataSource() datasource.DataSource {
	return &stagingBackendDataSource{}
}

// data source implementation.
type stagingBackendDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *stagingBackendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_backend"
}

// Schema defines the schema for the data source.
func (d *stagingBackendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Read a staging backend.",
		MarkdownDescription: "Read a staging backend.",

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
				Required:            true,
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
				Description:         "Use TLS encryption when contacting with the origin backend.",
				MarkdownDescription: "Use TLS encryption when contacting with the origin backend.",
			},
			"port": schema.Int64Attribute{
				Computed:            true,
				Description:         "Port where the origin is listening to HTTP requests, for example: 80 or 443.",
				MarkdownDescription: "Port where the origin is listening to HTTP requests, for example: `80` or `443`.",
			},
			"headers": schema.StringAttribute{
				Computed:            true,
				Description:         "Extra headers needed in order to validate backend status.",
				MarkdownDescription: "Extra headers needed in order to validate backend status.",
			},
			"hchost": schema.StringAttribute{
				Computed:    true,
				Description: "Host header that the health check probe will send to the origin, for example: www.my-origin.com.",
			},
			"hcpath": schema.StringAttribute{
				Computed:            true,
				Description:         "Host header that the health check probe will send to the origin, for example: www.my-origin.com.",
				MarkdownDescription: "Host header that the health check probe will send to the origin, for example: `www.my-origin.com`.",
			},
			"hcstatuscode": schema.Int64Attribute{
				Computed:            true,
				Description:         "Status code expected when the probe receives the HTTP health check response, for example: 200.",
				MarkdownDescription: "Status code expected when the probe receives the HTTP health check response, for example: `200`.",
			},
			"hcinterval": schema.Int64Attribute{
				Computed:            true,
				Description:         "Interval in seconds within which the probes of each edge execute the HTTP request to validate the status of the backend.",
				MarkdownDescription: "Interval in seconds within which the probes of each edge execute the HTTP request to validate the status of the backend.",
			},
			"hcdisabled": schema.BoolAttribute{
				Computed:            true,
				Description:         "Whether the health check probe is disabled.",
				MarkdownDescription: "Whether the health check probe is disabled.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *stagingBackendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StagingBackend

	// Read the config to the state
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	stagingBackend, err := d.client.GetBackendByName(state.Name.ValueString(), teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read the backend with name: %+v", state.Name),
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.Int64Value(int64(stagingBackend.ID))
	state.Company = types.Int64Value(int64(stagingBackend.Company))
	// state.Name = types.StringValue(stagingBackend.Name)
	state.VclName = types.StringValue("c" + strconv.Itoa(stagingBackend.Company) + "_" + stagingBackend.Name)
	state.Origin = types.StringValue(stagingBackend.Origin)
	state.Ssl = types.BoolValue(stagingBackend.Ssl)
	state.Port = types.Int64Value(int64(stagingBackend.Port))
	state.Headers = types.StringValue(stagingBackend.Headers)
	state.HCHost = types.StringValue(stagingBackend.HCHost)
	state.HCPath = types.StringValue(stagingBackend.HCPath)
	state.HCStatusCode = types.Int64Value(int64(stagingBackend.HCStatusCode))
	state.HCInterval = types.Int64Value(int64(stagingBackend.HCInterval))
	state.HCDisabled = types.BoolValue(stagingBackend.HCDisabled)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *stagingBackendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
