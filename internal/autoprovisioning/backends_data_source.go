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

// maps the data source schema data.
type backendsDataSourceModel struct {
	Backends []backendsModel `tfsdk:"backends"`
}

// maps schema data.
type backendsModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Company      types.Int64  `tfsdk:"company"`
	Name         types.String `tfsdk:"name"`
	Origin       types.String `tfsdk:"origin"`
	Ssl          types.Bool   `tfsdk:"ssl"`
	Port         types.Int64  `tfsdk:"port"`
	HCHost       types.String `tfsdk:"hchost"`
	HCPath       types.String `tfsdk:"hcpath"`
	HCStatusCode types.Int64  `tfsdk:"hcstatuscode"`
}

// Metadata returns the data source type name.
func (d *backendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backends"
}

// Schema defines the schema for the data source.
func (d *backendsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backends": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of all backends",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "ID of the backend",
						},
						"company": schema.Int64Attribute{
							Computed:    true,
							Description: "Company ID that owns this backend",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the backend",
						},
						"origin": schema.StringAttribute{
							Computed:    true,
							Description: "Origin is the IP or DNS address to the backend, for example: 'my-origin.com'",
						},
						"ssl": schema.BoolAttribute{
							Computed:    true,
							Description: "If the origin should be contacted using TLS encription.",
						},
						"port": schema.Int64Attribute{
							Computed:    true,
							Description: "Port where the origin is listening to HTTP requests, for example: 80 or 443",
						},
						"hchost": schema.StringAttribute{
							Computed:    true,
							Description: "Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com",
						},
						"hcpath": schema.StringAttribute{
							Computed:    true,
							Description: "Path that the healthcheck probe will used, for example: /favicon.ico",
						},
						"hcstatuscode": schema.Int64Attribute{
							Computed:    true,
							Description: "Status code expected when the probe receives the HTTP healthcheck response, for example: 200",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *backendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backendsDataSourceModel

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
		backendState := backendsModel{
			ID:           types.Int64Value(int64(backend.ID)),
			Company:      types.Int64Value(int64(backend.Company)),
			Name:         types.StringValue(backend.Name),
			Origin:       types.StringValue(backend.Origin),
			Ssl:          types.BoolValue(backend.Ssl),
			Port:         types.Int64Value(int64(backend.Port)),
			HCHost:       types.StringValue(backend.HCHost),
			HCPath:       types.StringValue(backend.HCPath),
			HCStatusCode: types.Int64Value(int64(backend.HCStatusCode)),
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
