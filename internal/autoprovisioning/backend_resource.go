package autoprovisioning

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &backendResource{}
	_ resource.ResourceWithConfigure   = &backendResource{}
	_ resource.ResourceWithImportState = &backendResource{}
)

// helper function to simplify the provider implementation.
func NewBackendResource() resource.Resource {
	return &backendResource{}
}

// resource implementation.
type backendResource struct {
	client *teclient.Client
}

// maps schema data.
type backendResourceModel struct {
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
func (r *backendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Schema defines the schema for the resource.
func (r *backendResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "ID of the backend",
			},
			"company": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Company ID that owns this backend",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(regexp.MustCompile("[^0-9][a-z0-9]+"), "The name must contain only lower case letters and numbers (cannot start with a number)"),
					),
				},
				Description: "Name of the backend",
			},
			"origin": schema.StringAttribute{
				Required:    true,
				Description: "Origin is the IP or DNS address to the backend, for example: 'my-origin.com'",
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
					int64validator.Between(200, 499),
				},
				Description: "Status code expected when the probe receives the HTTP healthcheck response, for example: 200",
			},
		},
	}
}

// Create
func (r *backendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating backend: "+plan.Name.ValueString())
	newBackend := teclient.NewBackendAPIModel{
		Name:         plan.Name.ValueString(),
		Origin:       plan.Origin.ValueString(),
		Ssl:          plan.Ssl.ValueBool(),
		Port:         int(plan.Port.ValueInt64()),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
	}
	backendState, errCreate := r.client.CreateBackend(newBackend, teclient.ProdEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating backend",
			fmt.Sprintf("Could not create the backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(backendState.ID))
	plan.Company = types.Int64Value(int64(backendState.Company))
	plan.Name = types.StringValue(backendState.Name)
	plan.Origin = types.StringValue(backendState.Origin)
	plan.Ssl = types.BoolValue(backendState.Ssl)
	plan.Port = types.Int64Value(int64(backendState.Port))
	plan.HCHost = types.StringValue(backendState.HCHost)
	plan.HCPath = types.StringValue(backendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(backendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *backendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating backend: "+plan.Name.ValueString())
	newBackend := teclient.BackendAPIModel{
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
	backendState, errCreate := r.client.UpdateBackend(newBackend, teclient.ProdEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error updating backend",
			fmt.Sprintf("Could not update the backend '%s': %s", plan.Name.ValueString(), errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(backendState.ID))
	plan.Company = types.Int64Value(int64(backendState.Company))
	plan.Name = types.StringValue(backendState.Name)
	plan.Origin = types.StringValue(backendState.Origin)
	plan.Ssl = types.BoolValue(backendState.Ssl)
	plan.Port = types.Int64Value(int64(backendState.Port))
	plan.HCHost = types.StringValue(backendState.HCHost)
	plan.HCPath = types.StringValue(backendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(backendState.HCStatusCode))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to find by ID
	if !state.ID.IsNull() {
		if backend, err := r.client.GetBackend(int(state.ID.ValueInt64()), teclient.ProdEnv); err == nil {
			state.ID = types.Int64Value(int64(backend.ID))
			state.Company = types.Int64Value(int64(backend.Company))
			state.Name = types.StringValue(backend.Name)
			state.Origin = types.StringValue(backend.Origin)
			state.Ssl = types.BoolValue(backend.Ssl)
			state.Port = types.Int64Value(int64(backend.Port))
			state.HCHost = types.StringValue(backend.HCHost)
			state.HCPath = types.StringValue(backend.HCPath)
			state.HCStatusCode = types.Int64Value(int64(backend.HCStatusCode))
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Try to find by Name
	backends, err := r.client.GetBackends(teclient.ProdEnv)
	if err == nil {
		for _, backend := range backends {
			if backend.Name == state.Name.ValueString() {
				state.ID = types.Int64Value(int64(backend.ID))
				state.Company = types.Int64Value(int64(backend.Company))
				state.Name = types.StringValue(backend.Name)
				state.Origin = types.StringValue(backend.Origin)
				state.Ssl = types.BoolValue(backend.Ssl)
				state.Port = types.Int64Value(int64(backend.Port))
				state.HCHost = types.StringValue(backend.HCHost)
				state.HCPath = types.StringValue(backend.HCPath)
				state.HCStatusCode = types.Int64Value(int64(backend.HCStatusCode))
				resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
				return
			}
		}
	}

	// Not found
	resp.Diagnostics.AddError("Backend not found", "Backend '"+state.Name.ValueString()+"' doesn't exist in API")
}

// Delete
func (r *backendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 204 on successful delete
	tflog.Info(ctx, "Deleting backend: '"+state.Name.ValueString()+"' with id: "+state.ID.String())
	if err := r.client.DeleteBackend(int(state.ID.ValueInt64()), teclient.ProdEnv); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting a backend",
			"Could not delete the backend: "+state.Name.ValueString()+"\n"+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *backendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *backendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
