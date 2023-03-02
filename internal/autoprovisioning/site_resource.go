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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

// siteModel maps schema data.
type siteResourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Domain   types.String   `tfsdk:"domain"`
	ID       types.Int64    `tfsdk:"id"`
	Active   types.Bool     `tfsdk:"active"`
}

// Metadata returns the resource type name.
func (r *siteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Schema defines the schema for the resource.
func (r *siteResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "ID of the site",
			},
			"active": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Active status in the CDN",
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "Domain in FDQN form, i.e: 'www.example.com'",
			},
		},
	}
}

// Create new site
func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maxTimeout, err := plan.Timeouts.Create(ctx, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying timeouts",
			"Could not apply the timeouts configuration.",
		)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, maxTimeout)
	defer cancel()

	tflog.Info(ctx, "Creating site: "+plan.Domain.ValueString())
	siteState, errCreate := r.HelperCreateSite(plan.Domain.ValueString(), maxTimeout, ctx)
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
	var plan siteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maxTimeout, err := plan.Timeouts.Create(ctx, 5*time.Minute)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error applying timeouts",
			"Could not apply the timeouts configuration.",
		)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, maxTimeout)
	defer cancel()

	siteState, errCreate := r.HelperCreateSite(plan.Domain.ValueString(), maxTimeout, ctx)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error updating site",
			fmt.Sprintf("Could not update the site '%s': %s", plan.Domain.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = siteState.ID
	plan.Domain = siteState.Domain
	plan.Active = siteState.Active
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state siteResourceModel
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
	var state siteResourceModel
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
func (r *siteResource) HelperCreateSite(domain string, maxTimeout time.Duration, ctx context.Context) (*siteResourceModel, error) {
	createState := &sdkresource.StateChangeConf{
		Pending: []string{
			"verify_error",
		},
		Target: []string{
			"site_verified",
		},
		ContinuousTargetOccurence: 1, // How many times the Target has to be reached to continue
		Timeout:                   maxTimeout,
		Delay:                     1 * time.Second,  // Delay before starting
		MinTimeout:                30 * time.Second, // Delay between retries
		Refresh: func() (any, string, error) {
			siteCreate := teclient.SiteNewAPIModel{
				Url: domain,
			}
			newSite, verify_error, err := r.client.CreateSite(siteCreate)
			if err == nil {
				siteState := siteResourceModel{
					Domain: types.StringValue(newSite.Url),
					ID:     types.Int64Value(int64(newSite.ID)),
					Active: types.BoolValue(true), // API is not sending back "active" field
				}
				return siteState, "site_verified", nil
			} else if verify_error {
				// We cannot return the error inmediately or the refresh function will exit
				// instead, return the error in the interface position and nil on the error
				return err, "verify_error", nil
			}
			return nil, "create_error", fmt.Errorf("Could not create the site: %s; Error: %s", domain, err.Error())
		},
	}

	siteState, errState := createState.WaitForStateContext(ctx)
	// Verify error
	if siteState, ok := siteState.(error); ok {
		return nil, siteState // here siteState is an error
	}
	if errState != nil {
		return nil, errState
	}

	if siteState, ok := siteState.(siteResourceModel); ok {
		return &siteState, errState
	}

	// Cannot assert the response matches the model...
	return nil, fmt.Errorf("Could not create the site, unknown error.")
}
