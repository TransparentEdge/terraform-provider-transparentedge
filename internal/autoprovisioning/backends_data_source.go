package autoprovisioning

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &backendsDataSource{}
	_ datasource.DataSourceWithConfigure = &backendsDataSource{}
)

// Helper function to simplify the provider implementation.
func NewBackendsDataSource() datasource.DataSource {
	return &backendsDataSource{}
}

// data source implementation.
type backendsDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *backendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backends"
}

// Schema defines the schema for the data source.
func (d *backendsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Backend listing.",
		MarkdownDescription: "Backend listing.",

		Attributes: map[string]schema.Attribute{
			"backends": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of all backends.",
				MarkdownDescription: "List of all backends.",
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
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
							Computed:            true,
							Description:         "Host header that the health check probe will send to the origin, for example: www.my-origin.com.",
							MarkdownDescription: "Host header that the health check probe will send to the origin, for example: `www.my-origin.com`.",
						},
						"hcpath": schema.StringAttribute{
							Computed:            true,
							Description:         "Path that the health check probe will use, for example: /favicon.ico.",
							MarkdownDescription: "Path that the health check probe will use, for example: `/favicon.ico`.",
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *backendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Backends

	backends, err := d.client.GetBackends(teclient.ProdEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Backends info",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, backend := range backends {
		backendState := Backend{
			ID:           types.Int64Value(int64(backend.ID)),
			Company:      types.Int64Value(int64(backend.Company)),
			Name:         types.StringValue(backend.Name),
			VclName:      types.StringValue("c" + strconv.Itoa(backend.Company) + "_" + backend.Name),
			Origin:       types.StringValue(backend.Origin),
			Ssl:          types.BoolValue(backend.Ssl),
			Port:         types.Int64Value(int64(backend.Port)),
			Headers:      types.StringValue(backend.Headers),
			HCHost:       types.StringValue(backend.HCHost),
			HCPath:       types.StringValue(backend.HCPath),
			HCStatusCode: types.Int64Value(int64(backend.HCStatusCode)),
			HCInterval:   types.Int64Value(int64(backend.HCInterval)),
			HCDisabled:   types.BoolValue(backend.HCDisabled),
		}

		state.Backends = append(state.Backends, backendState)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *backendsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
