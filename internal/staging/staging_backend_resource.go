package staging

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &stagingBackendResource{}
	_ resource.ResourceWithConfigure   = &stagingBackendResource{}
	_ resource.ResourceWithImportState = &stagingBackendResource{}
)

// helper function to simplify the provider implementation.
func NewStagingBackendResource() resource.Resource {
	return &stagingBackendResource{}
}

// resource implementation.
type stagingBackendResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (r *stagingBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_backend"
}

// Schema defines the schema for the resource.
func (r *stagingBackendResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages staging backend configuration",
		MarkdownDescription: "Manages staging backend configuration",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "ID of the staging backend",
				MarkdownDescription: "ID of the staging backend",
			},
			"company": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "Company ID that owns this staging backend",
				MarkdownDescription: "Company ID that owns this staging backend",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(regexp.MustCompile("[^0-9][a-z0-9]+"), "The name must contain only lower case letters and numbers (cannot start with a number)"),
					),
				},
				Description:         "Name of the staging backend",
				MarkdownDescription: "Name of the staging backend",
			},
			"vclname": schema.StringAttribute{
				Computed:            true,
				Description:         "Final unique name of the backend to be referencen in VCL Code: 'c{company_id}_{name}'",
				MarkdownDescription: "Final unique name of the backend to be referencen in VCL Code: `c{company_id}_{name}`",
			},
			"origin": schema.StringAttribute{
				Required:    true,
				Description: "Origin is the IP or DNS address to the origin backend, for example: 'my-origin.com'",
			},
			"ssl": schema.BoolAttribute{
				Required:    true,
				Description: "If the origin should be contacted using TLS encription.",
			},
			"port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(80, 65535),
				},
				Description: "Port where the origin is listening to HTTP requests, for example: 80 or 443",
			},
			"hchost": schema.StringAttribute{
				Required:    true,
				Description: "Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com",
			},
			"hcpath": schema.StringAttribute{
				Required:    true,
				Description: "Path that the healthcheck probe will used, for example: /favicon.ico",
			},
			"hcstatuscode": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(200, 599),
				},
				Description: "Status code expected when the probe receives the HTTP healthcheck response, for example: 200",
			},
		},
	}
}

// Create
func (r *stagingBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan StagingBackend
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Staging Backend: "+plan.Name.ValueString())
	newStagingBackend := teclient.NewBackendAPIModel{
		Name:         plan.Name.ValueString(),
		Origin:       plan.Origin.ValueString(),
		Ssl:          plan.Ssl.ValueBool(),
		Port:         int(plan.Port.ValueInt64()),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
	}
	stagingBackendState, errCreate := r.client.CreateBackend(newStagingBackend, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating staging backend",
			fmt.Sprintf("Could not create the staging backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingBackendState.ID))
	plan.Company = types.Int64Value(int64(stagingBackendState.Company))
	plan.Name = types.StringValue(stagingBackendState.Name)
	plan.VclName = types.StringValue("c" + strconv.Itoa(stagingBackendState.Company) + "_" + stagingBackendState.Name)
	plan.Origin = types.StringValue(stagingBackendState.Origin)
	plan.Ssl = types.BoolValue(stagingBackendState.Ssl)
	plan.Port = types.Int64Value(int64(stagingBackendState.Port))
	plan.HCHost = types.StringValue(stagingBackendState.HCHost)
	plan.HCPath = types.StringValue(stagingBackendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(stagingBackendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *stagingBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StagingBackend
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating Staging Backend: "+plan.Name.ValueString())
	newStagingBackend := teclient.BackendAPIModel{
		ID:           int(plan.ID.ValueInt64()),
		Company:      int(plan.Company.ValueInt64()),
		Name:         plan.Name.ValueString(),
		Origin:       plan.Origin.ValueString(),
		Ssl:          plan.Ssl.ValueBool(),
		Port:         int(plan.Port.ValueInt64()),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
	}
	stagingBackendState, errCreate := r.client.UpdateBackend(newStagingBackend, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error updating staging backend",
			fmt.Sprintf("Could not update the staging backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingBackendState.ID))
	plan.Company = types.Int64Value(int64(stagingBackendState.Company))
	plan.Name = types.StringValue(stagingBackendState.Name)
	plan.VclName = types.StringValue("c" + strconv.Itoa(stagingBackendState.Company) + "_" + stagingBackendState.Name)
	plan.Origin = types.StringValue(stagingBackendState.Origin)
	plan.Ssl = types.BoolValue(stagingBackendState.Ssl)
	plan.Port = types.Int64Value(int64(stagingBackendState.Port))
	plan.HCHost = types.StringValue(stagingBackendState.HCHost)
	plan.HCPath = types.StringValue(stagingBackendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(stagingBackendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *stagingBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state StagingBackend
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to find by ID
	if !state.ID.IsNull() {
		if stagingBackend, err := r.client.GetBackend(int(state.ID.ValueInt64()), teclient.StagingEnv); err == nil {
			state.ID = types.Int64Value(int64(stagingBackend.ID))
			state.Company = types.Int64Value(int64(stagingBackend.Company))
			state.Name = types.StringValue(stagingBackend.Name)
			state.VclName = types.StringValue("c" + strconv.Itoa(stagingBackend.Company) + "_" + stagingBackend.Name)
			state.Origin = types.StringValue(stagingBackend.Origin)
			state.Ssl = types.BoolValue(stagingBackend.Ssl)
			state.Port = types.Int64Value(int64(stagingBackend.Port))
			state.HCHost = types.StringValue(stagingBackend.HCHost)
			state.HCPath = types.StringValue(stagingBackend.HCPath)
			state.HCStatusCode = types.Int64Value(int64(stagingBackend.HCStatusCode))
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Try to find by Name
	stagingBackends, err := r.client.GetBackends(teclient.StagingEnv)
	if err == nil {
		for _, stagingBackend := range stagingBackends {
			if stagingBackend.Name == state.Name.ValueString() {
				state.ID = types.Int64Value(int64(stagingBackend.ID))
				state.Company = types.Int64Value(int64(stagingBackend.Company))
				state.Name = types.StringValue(stagingBackend.Name)
				state.VclName = types.StringValue("c" + strconv.Itoa(stagingBackend.Company) + "_" + stagingBackend.Name)
				state.Origin = types.StringValue(stagingBackend.Origin)
				state.Ssl = types.BoolValue(stagingBackend.Ssl)
				state.Port = types.Int64Value(int64(stagingBackend.Port))
				state.HCHost = types.StringValue(stagingBackend.HCHost)
				state.HCPath = types.StringValue(stagingBackend.HCPath)
				state.HCStatusCode = types.Int64Value(int64(stagingBackend.HCStatusCode))
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
				return
			}
		}
	}

	// Not found
	resp.Diagnostics.AddError("Staging Backend not found", "Staging Backend '"+state.Name.ValueString()+"' doesn't exist in API")
}

// Delete
func (r *stagingBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StagingBackend
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 204 on successful delete
	tflog.Info(ctx, "Deleting Staging Backend: '"+state.Name.ValueString()+"' with id: "+state.ID.String())
	if err := r.client.DeleteBackend(int(state.ID.ValueInt64()), teclient.StagingEnv); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting a Staging Backend",
			"Could not delete the Staging Backend: "+state.Name.ValueString()+"\n"+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *stagingBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *stagingBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
