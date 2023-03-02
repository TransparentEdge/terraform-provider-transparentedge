package transparentedge

import (
	"context"
	"os"
	"strconv"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/autoprovisioning"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &transparentedgeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &transparentedgeProvider{}
}

type transparentedgeProvider struct {
	version string
}

// Metadata returns the provider type name.
func (p *transparentedgeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "transparentedge"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *transparentedgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host_url": schema.StringAttribute{
				Optional:    true,
				Description: "URL of Transparent Edge API. default: 'https://api.transparentcdn.com'. May also be provided via TCDN_HOST_URL environment variable.",
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
				Description: "Ignore SSL certificate for 'host_url'. May also be provided via TCDN_VERIFY_SSL environment variable.",
			},
		},
	}
}

// maps provider schema data to a Go type.
type transparentedgeProviderModel struct {
	HostURL      types.String `tfsdk:"host_url"`
	CompanyId    types.Int64  `tfsdk:"company_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	VerifySSL    types.Bool   `tfsdk:"verify_ssl"`
}

func (p *transparentedgeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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

	host_url := os.Getenv("TCDN_HOST_URL")
	companyid, err := strconv.Atoi(os.Getenv("TCDN_COMPANY_ID"))
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("company_id"),
			"Invalid Company ID value",
			"Company ID must be an integer. Please ensure that you have it set with the TCDN_COMPANY_ID environment variable.",
		)
	}

	clientid := os.Getenv("TCDN_CLIENT_ID")
	clientsecret := os.Getenv("TCDN_CLIENT_SECRET")

	verifyssl_str, verifyssl_set := os.LookupEnv("TCDN_VERIFY_SSL")
	verifyssl := true
	if verifyssl_set {
		var err error
		verifyssl, err = strconv.ParseBool(verifyssl_str)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("verify_ssl"),
				"Invalid Verify SSL value",
				"Verify SSL must be a boolean. Please ensure that you have it set with the TCDN_VERIFY_SSL environment variable.",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Override with terraform configuration values
	if !config.HostURL.IsNull() {
		host_url = config.HostURL.ValueString()
	}
	if !config.CompanyId.IsNull() {
		companyid = int(config.CompanyId.ValueInt64())
	}
	if !config.ClientId.IsNull() {
		clientid = config.ClientId.ValueString()
	}
	if !config.ClientSecret.IsNull() {
		clientid = config.ClientSecret.ValueString()
	}
	if !config.VerifySSL.IsNull() {
		verifyssl = config.VerifySSL.ValueBool()
	}

	// Default values
	if host_url == "" {
		host_url = "https://api.transparentcdn.com"
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if companyid <= 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("company_id"),
			"Invalid Company ID",
			"Please provide a valid Company ID. Set the TCDN_COMPANY_ID environment variable.",
		)
	}
	if clientid == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Client ID",
			"Please provide a valid Client ID. Set the TCDN_CLIENT_ID environment variable.",
		)
	}
	if clientsecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Client Secret",
			"Please provide a valid Client Secret. Set the TCDN_CLIENT_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "tedge_host_url", host_url)
	ctx = tflog.SetField(ctx, "tedge_companyid", companyid)

	tflog.Debug(ctx, "Creating Transparent Edge API client")

	// Create a new client using the configuration values
	client, err := teclient.NewClient(&host_url, &companyid, &clientid, &clientsecret, verifyssl)
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
func (p *transparentedgeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		autoprovisioning.NewSitesDataSource,
		autoprovisioning.NewSiteVerifyDataSource,
		autoprovisioning.NewBackendsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *transparentedgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		autoprovisioning.NewSiteResource,
		autoprovisioning.NewBackendResource,
	}
}
