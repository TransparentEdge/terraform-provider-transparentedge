package staging

import (
	"context"
	"fmt"
	"regexp"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &stagingbackendResource{}
	_ resource.ResourceWithConfigure   = &stagingbackendResource{}
	_ resource.ResourceWithImportState = &stagingbackendResource{}
)

// helper function to simplify the provider implementation.
func NewStagingBackendResource() resource.Resource {
	return &stagingbackendResource{}
}

// resource implementation.
type stagingbackendResource struct {
	client *teclient.Client
}

// maps schema data.
type stagingbackendResourceModel struct {
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

// Metadata returns the resource type name.
func (r *stagingbackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stagingbackend"
}

// Schema defines the schema for the resource.
func (r *stagingbackendResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the staging backend",
			},
			"company": schema.Int64Attribute{
				Computed:    true,
				Description: "Company ID that owns this staging backend",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(regexp.MustCompile("[^0-9][a-z0-9]+"), "The name must contain only lower case letters and numbers (cannot start with a number)"),
					),
				},
				Description: "Name of the staging backend",
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
func (r *stagingbackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stagingbackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating stagingbackend: "+plan.Name.ValueString())
	newStagingBackend := teclient.NewBackendAPIModel{
		Name:         plan.Name.ValueString(),
		Origin:       plan.Origin.ValueString(),
		Ssl:          plan.Ssl.ValueBool(),
		Port:         int(plan.Port.ValueInt64()),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
	}
	stagingbackendState, errCreate := r.client.CreateBackend(newStagingBackend, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating staging backend",
			fmt.Sprintf("Could not create the staging backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingbackendState.ID))
	plan.Company = types.Int64Value(int64(stagingbackendState.Company))
	plan.Name = types.StringValue(stagingbackendState.Name)
	plan.Origin = types.StringValue(stagingbackendState.Origin)
	plan.Ssl = types.BoolValue(stagingbackendState.Ssl)
	plan.Port = types.Int64Value(int64(stagingbackendState.Port))
	plan.HCHost = types.StringValue(stagingbackendState.HCHost)
	plan.HCPath = types.StringValue(stagingbackendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(stagingbackendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *stagingbackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan stagingbackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating stagingbackend: "+plan.Name.ValueString())
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
	stagingbackendState, errCreate := r.client.UpdateBackend(newStagingBackend, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error updating staging backend",
			fmt.Sprintf("Could not update the staging backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingbackendState.ID))
	plan.Company = types.Int64Value(int64(stagingbackendState.Company))
	plan.Name = types.StringValue(stagingbackendState.Name)
	plan.Origin = types.StringValue(stagingbackendState.Origin)
	plan.Ssl = types.BoolValue(stagingbackendState.Ssl)
	plan.Port = types.Int64Value(int64(stagingbackendState.Port))
	plan.HCHost = types.StringValue(stagingbackendState.HCHost)
	plan.HCPath = types.StringValue(stagingbackendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(stagingbackendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *stagingbackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state stagingbackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to find by ID
	if !state.ID.IsNull() {
		if stagingbackend, err := r.client.GetBackend(int(state.ID.ValueInt64()), teclient.StagingEnv); err == nil {
			state.ID = types.Int64Value(int64(stagingbackend.ID))
			state.Company = types.Int64Value(int64(stagingbackend.Company))
			state.Name = types.StringValue(stagingbackend.Name)
			state.Origin = types.StringValue(stagingbackend.Origin)
			state.Ssl = types.BoolValue(stagingbackend.Ssl)
			state.Port = types.Int64Value(int64(stagingbackend.Port))
			state.HCHost = types.StringValue(stagingbackend.HCHost)
			state.HCPath = types.StringValue(stagingbackend.HCPath)
			state.HCStatusCode = types.Int64Value(int64(stagingbackend.HCStatusCode))
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Try to find by Name
	stagingbackends, err := r.client.GetBackends(teclient.StagingEnv)
	if err == nil {
		for _, stagingbackend := range stagingbackends {
			if stagingbackend.Name == state.Name.ValueString() {
				state.ID = types.Int64Value(int64(stagingbackend.ID))
				state.Company = types.Int64Value(int64(stagingbackend.Company))
				state.Name = types.StringValue(stagingbackend.Name)
				state.Origin = types.StringValue(stagingbackend.Origin)
				state.Ssl = types.BoolValue(stagingbackend.Ssl)
				state.Port = types.Int64Value(int64(stagingbackend.Port))
				state.HCHost = types.StringValue(stagingbackend.HCHost)
				state.HCPath = types.StringValue(stagingbackend.HCPath)
				state.HCStatusCode = types.Int64Value(int64(stagingbackend.HCStatusCode))
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
				return
			}
		}
	}

	// Not found
	resp.Diagnostics.AddError("StagingBackend not found", "Staging Backend '"+state.Name.ValueString()+"' doesn't exist in API")
}

// Delete
func (r *stagingbackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state stagingbackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 204 on successful delete
	tflog.Info(ctx, "Deleting stagingbackend: '"+state.Name.ValueString()+"' with id: "+state.ID.String())
	if err := r.client.DeleteBackend(int(state.ID.ValueInt64()), teclient.StagingEnv); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting a stagingbackend",
			"Could not delete the stagingbackend: "+state.Name.ValueString()+"\n"+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *stagingbackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *stagingbackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
