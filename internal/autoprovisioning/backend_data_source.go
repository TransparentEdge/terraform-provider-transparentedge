package autoprovisioning

import (
	"context"
	"fmt"
	"strconv"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &backendDataSource{}
	_ datasource.DataSourceWithConfigure = &backendDataSource{}
)

// Helper function to simplify the provider implementation.
func NewBackendDataSource() datasource.DataSource {
	return &backendDataSource{}
}

// data source implementation.
type backendDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *backendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Schema defines the schema for the data source.
func (d *backendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Read a backend.",
		MarkdownDescription: "Read a backend.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the backend.",
				MarkdownDescription: "ID of the backend.",
			},
			"company": schema.Int64Attribute{
				Computed:            true,
				Description:         "Company ID that owns this backend.",
				MarkdownDescription: "Company ID that owns this backend.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the backend.",
				MarkdownDescription: "Name of the backend.",
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
func (d *backendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Backend

	// Read the config to the state
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	backend, err := d.client.GetBackendByName(state.Name.ValueString(), teclient.ProdEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to read the backend with name: %+v", state.Name),
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.Int64Value(int64(backend.ID))
	state.Company = types.Int64Value(int64(backend.Company))
	//state.Name = types.StringValue(backend.Name)
	state.VclName = types.StringValue("c" + strconv.Itoa(backend.Company) + "_" + backend.Name)
	state.Origin = types.StringValue(backend.Origin)
	state.Ssl = types.BoolValue(backend.Ssl)
	state.Port = types.Int64Value(int64(backend.Port))
	state.Headers = types.StringValue(backend.Headers)
	state.HCHost = types.StringValue(backend.HCHost)
	state.HCPath = types.StringValue(backend.HCPath)
	state.HCStatusCode = types.Int64Value(int64(backend.HCStatusCode))
	state.HCInterval = types.Int64Value(int64(backend.HCInterval))
	state.HCDisabled = types.BoolValue(backend.HCDisabled)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *backendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
