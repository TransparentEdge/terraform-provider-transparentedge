package autoprovisioning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	defaultCreateTimeout    time.Duration = 5 * time.Minute
	delayBetweenCreateRetry time.Duration = 30 * time.Second
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &siteResource{}
	_ resource.ResourceWithConfigure   = &siteResource{}
	_ resource.ResourceWithImportState = &siteResource{}
)

// NewSiteResource is a helper function to simplify the provider implementation.
func NewSiteResource() resource.Resource {
	return &siteResource{}
}

// siteResource is the resource implementation.
type siteResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema defines the schema for the resource.
func (r *siteResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages company sites (domains)",
		MarkdownDescription: "Manages company sites (domains)",

		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the site",
				MarkdownDescription: "ID of the site",
			},
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Domain in FDQN form, i.e: 'www.example.com'",
				MarkdownDescription: "Domain in FDQN form, i.e: `www.example.com`",
			},
			"active": schema.BoolAttribute{
				Computed:            true,
				Description:         "Internal value that indicates if the site is active in the CDN",
				MarkdownDescription: "Internal value that indicates if the site is active in the CDN",
			},
		},
	}
}

// Create new site
func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Site
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maxTimeout, err := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying timeouts",
			"Could not apply the timeouts configuration.",
		)
		return
	}

	tflog.Info(ctx, "Creating site: "+plan.Domain.ValueString())
	siteState, errCreate := r.HelperCreateSite(ctx, plan.Domain.ValueString(), maxTimeout)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			fmt.Sprintf("Could not create the site '%s': %s", plan.Domain.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = siteState.ID
	plan.Domain = siteState.Domain
	plan.Active = siteState.Active
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Read resource information
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Site
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to find by ID
	if !state.ID.IsNull() {
		if siteAPI, err := r.client.GetSite(int(state.ID.ValueInt64())); err == nil {
			if siteAPI.Active {
				state.ID = types.Int64Value(int64(siteAPI.ID))
				state.Domain = types.StringValue(siteAPI.Url)
				state.Active = types.BoolValue(siteAPI.Active)
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
				return
			}
		}
	}

	// Try to find by Domain
	sites, err := r.client.GetSites()
	if err == nil {
		for _, siteAPI := range sites {
			if siteAPI.Url == state.Domain.ValueString() && siteAPI.Active {
				state.ID = types.Int64Value(int64(siteAPI.ID))
				state.Domain = types.StringValue(siteAPI.Url)
				state.Active = types.BoolValue(siteAPI.Active)
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
				return
			}
		}
	}

	// Not found or inactive
	if !strings.Contains(state.Domain.ValueString(), "<inactive>") {
		state.Domain = types.StringValue(state.Domain.ValueString() + " <inactive>")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes the site and removes the terraform plan on success
func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Site
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 204 on successful delete
	tflog.Info(ctx, "Deleting site: "+state.Domain.ValueString()+" with id: "+state.ID.String())
	if err := r.client.DeleteSite(int(state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting a site",
			"Could not delete the site: "+state.Domain.ValueString(),
		)
		return
	}

	// Sites are disabled, not deleted
	resp.Diagnostics.AddWarning(
		"Site Resource Destruction Considerations",
		"This action will disable the site and remove it from the Terraform state, but it does not permanently delete the site.\n"+
			"Remember, sites can only be re-enabled within the same company.\n"+
			"If there's a future requirement to reassign this site to a different company, please contact our support team for help.\n"+
			"If this action aligns with your current intentions, you may disregard this warning.",
	)
}

// Configure adds the provider configured client to the resource.
func (r *siteResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *siteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Helpers
func (r *siteResource) HelperCreateSite(ctx context.Context, domain string, maxTimeout time.Duration) (*Site, error) {
	var err error = nil
	remainingTimeForVerification := maxTimeout.Seconds()
	siteCreate := teclient.SiteNewAPIModel{Url: domain}
	siteState := Site{}

	site := &teclient.SiteAPIModel{}
	verify_error := false

	for {
		site, verify_error, err = r.client.CreateSite(siteCreate)
		if err == nil {
			siteState.ID = types.Int64Value(int64(site.ID))
			siteState.Domain = types.StringValue(site.Url)
			siteState.Active = types.BoolValue(true)
			break
		}

		// break if err != nil and it wasn't a verification error
		if !verify_error {
			break
		}

		// break if we exhausted the remaining time
		if remainingTimeForVerification <= 0 {
			break
		}
		time.Sleep(delayBetweenCreateRetry)
		remainingTimeForVerification -= (delayBetweenCreateRetry.Seconds() + 5)
		tflog.Info(ctx, "Retry site verification for "+domain)
	}

	if err == nil {
		return &siteState, nil
	}

	return nil, err
}
