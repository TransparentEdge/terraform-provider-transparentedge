package transparentedge

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/autoprovisioning"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/companies"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
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
}

// Metadata returns the provider type name.
func (p *TransparentEdgeProvider) Metadata(ctx context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "transparentedge"
	resp.Version = p.version
	tflog.Info(ctx, "Version: "+p.version+", Commit: "+p.commit+", Date: "+time.Now().UTC().String())
}

// Schema defines the provider-level schema for configuration data.
func (p *TransparentEdgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A terraform provider for the CDN of Transparent Edge.",
		MarkdownDescription: "A terraform provider for the CDN of [Transparent Edge](https://www.transparentedge.eu/).",

		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional:            true,
				Description:         "URL of Transparent Edge API. default: 'https://api.transparentcdn.com'. May also be provided via TCDN_API_URL environment variable.",
				MarkdownDescription: "URL of Transparent Edge API. default: `https://api.transparentcdn.com`. May also be provided via `TCDN_API_URL` environment variable.",
			},
			"company_id": schema.Int64Attribute{
				Optional:            true,
				Description:         "Company ID number (for ex: 300). May also be provided via TCDN_COMPANY_ID environment variable.",
				MarkdownDescription: "Company ID number (for ex: `300`). May also be provided via `TCDN_COMPANY_ID` environment variable.",
			},
			"client_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Client ID (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_ID environment variable.",
				MarkdownDescription: "Client ID (`dashboard -> profile -> account options -> manage keys`). May also be provided via `TCDN_CLIENT_ID` environment variable.",
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Client Secret (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_SECRET environment variable.",
				MarkdownDescription: "Client Secret (`dashboard -> profile -> account options -> manage keys`). May also be provided via `TCDN_CLIENT_SECRET` environment variable.",
			},
			"insecure": schema.BoolAttribute{
				Optional:            true,
				Description:         "Ignore TLS certificate for 'api_url'. May also be provided via TCDN_INSECURE environment variable.",
				MarkdownDescription: "Ignore TLS certificate for `api_url`. May also be provided via `TCDN_INSECURE` environment variable.",
			},
			"auth": schema.BoolAttribute{
				Optional:            true,
				Description:         "Set to false if your configuration only consumes data sources that do not require authentication, such as `transparentedge_ip_ranges`.",
				MarkdownDescription: "Set to false if your configuration only consumes data sources that do not require authentication, such as `transparentedge_ip_ranges`.",
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
	Insecure     types.Bool   `tfsdk:"insecure"`
	Auth         types.Bool   `tfsdk:"auth"`
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

	companyid, _ := helpers.GetIntEnv("TCDN_COMPANY_ID", 0)
	insecure, _ := helpers.GetEnvBool("TCDN_INSECURE", false)

	auth := true

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
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}
	if !config.Auth.IsNull() {
		auth = config.Auth.ValueBool()
	}

	// Values that need conversion (if not set in the configuration)

	// Default values
	if api_url == "" {
		api_url = defaultApiUrl
	}

	if !auth {
		// Override unset values and let the client setup
		if companyid < 1 {
			companyid = 1
		}
		if clientid == "" {
			clientid = "noauth"
		}
		if clientsecret == "" {
			clientsecret = "noauth"
		}
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
	if !(strings.HasPrefix(api_url, "http://") || strings.HasPrefix(api_url, "https://")) {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Invalid API URL",
			"Please provide a valid API URL value, protocol (http or https) is required.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "tedge_api_url", api_url)
	ctx = tflog.SetField(ctx, "tedge_companyid", companyid)
	tflog.Debug(ctx, "Creating Transparent Edge API client")

	// Create a new client using the configuration values
	useragent := "terraform-provider-transparentedge/" + p.version
	client, err := teclient.NewClient(&api_url, &companyid, &clientid, &clientsecret, &insecure, &auth, &useragent)
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
		autoprovisioning.NewBackendDataSource,
		autoprovisioning.NewBackendsDataSource,
		autoprovisioning.NewVclconfDataSource,
		autoprovisioning.NewCertificatesDataSource,
		autoprovisioning.NewCertReqDNSProvidersDataSource,
		autoprovisioning.NewCertReqDNSCNAMEVerifDataSource,
		autoprovisioning.NewCertReqDNSCredentialDataSource,
		autoprovisioning.NewCertReqDNSDataSource,
		staging.NewStagingBackendDataSource,
		staging.NewStagingBackendsDataSource,
		staging.NewStagingVclconfDataSource,
		companies.NewIpRangesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *TransparentEdgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		autoprovisioning.NewSiteResource,
		autoprovisioning.NewBackendResource,
		autoprovisioning.NewVclconfResource,
		autoprovisioning.NewCustomCertificate,
		autoprovisioning.NewCertReqDNSCredentialResource,
		autoprovisioning.NewCertReqDNSResource,
		staging.NewStagingBackendResource,
		staging.NewStagingVclconfResource,
	}
}

func New(version string, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &TransparentEdgeProvider{
			version: version,
			commit:  commit,
		}
	}
}
