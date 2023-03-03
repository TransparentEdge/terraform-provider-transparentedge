package transparentedge

import (
	"context"
	"os"
	"strconv"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/autoprovisioning"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/staging"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultApiUrl = "https://api.transparentcdn.com"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &TransparentEdgeProvider{}
)

type TransparentEdgeProvider struct {
	version string
	commit  string
	date    string
}

// Metadata returns the provider type name.
func (p *TransparentEdgeProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "transparentedge"
	resp.Version = p.version
	tflog.Info(ctx, "Version: "+p.version+", Commit: "+p.commit+", Date: "+p.date)
}

// Schema defines the provider-level schema for configuration data.
func (p *TransparentEdgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL of Transparent Edge API. default: 'https://api.transparentcdn.com'. May also be provided via TCDN_API_URL environment variable.",
			},
			"company_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Company ID number (for ex: 300). May also be provided via TCDN_COMPANY_ID environment variable.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Client ID (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_ID environment variable.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Client Secret (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_SECRET environment variable.",
			},
			"verify_ssl": schema.BoolAttribute{
				Optional:    true,
				Description: "Ignore SSL certificate for 'api_url'. May also be provided via TCDN_VERIFY_SSL environment variable.",
			},
		},
	}
}

// maps provider schema data to a Go type.
type transparentedgeProviderModel struct {
	ApiURL       types.String `tfsdk:"api_url"`
	CompanyId    types.Int64  `tfsdk:"company_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	VerifySSL    types.Bool   `tfsdk:"verify_ssl"`
}

func (p *TransparentEdgeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Transparent Edge API client")

	// Retrieve provider data from configuration
	var config transparentedgeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	api_url := os.Getenv("TCDN_API_URL")
	clientid := os.Getenv("TCDN_CLIENT_ID")
	clientsecret := os.Getenv("TCDN_CLIENT_SECRET")

	companyid := 0
	companyid, _ = strconv.Atoi(os.Getenv("TCDN_COMPANY_ID"))
	verifyssl, _ := strconv.ParseBool(os.Getenv("TCDN_VERIFY_SSL"))

	// Override with terraform configuration values
	if !config.ApiURL.IsNull() {
		api_url = config.ApiURL.ValueString()
	}
	if !config.CompanyId.IsNull() {
		companyid = int(config.CompanyId.ValueInt64())
	}
	if !config.ClientId.IsNull() {
		clientid = config.ClientId.ValueString()
	}
	if !config.ClientSecret.IsNull() {
		clientsecret = config.ClientSecret.ValueString()
	}
	if !config.VerifySSL.IsNull() {
		verifyssl = config.VerifySSL.ValueBool()
	}

	// Values that need conversion (if not set in the configuration)

	// Default values
	if api_url == "" {
		api_url = defaultApiUrl
	}

	if companyid < 1 {
		resp.Diagnostics.AddAttributeError(
			path.Root("company_id"),
			"Invalid Company ID value",
			"Company ID is an integer greater than 0.",
		)
	}
	if clientid == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Client ID",
			"Please provide a valid Client ID.",
		)
	}
	if clientsecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Client Secret",
			"Please provide a valid Client Secret.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "tedge_api_url", api_url)
	ctx = tflog.SetField(ctx, "tedge_companyid", companyid)
	tflog.Debug(ctx, "Creating Transparent Edge API client")

	// Create a new client using the configuration values
	client, err := teclient.NewClient(&api_url, &companyid, &clientid, &clientsecret, verifyssl)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Transparent Edge API Client",
			"An unexpected error occurred when creating the Transparent Edge API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Client Error: "+err.Error(),
		)
		return
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Transparent Edge API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *TransparentEdgeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		autoprovisioning.NewSitesDataSource,
		autoprovisioning.NewSiteVerifyDataSource,
		autoprovisioning.NewBackendsDataSource,
		autoprovisioning.NewVclconfDataSource,
		autoprovisioning.NewCertificatesDataSource,
		staging.NewStagingBackendsDataSource,
		staging.NewStagingVclconfDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *TransparentEdgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		autoprovisioning.NewSiteResource,
		autoprovisioning.NewBackendResource,
		autoprovisioning.NewVclconfResource,
		staging.NewStagingBackendResource,
		staging.NewStagingVclconfResource,
	}
}

func New(version string, commit string, date string) func() provider.Provider {
	return func() provider.Provider {
		return &TransparentEdgeProvider{
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
