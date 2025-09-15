package autoprovisioning

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

// Metadata returns the resource type name.
func (r *backendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Schema defines the schema for the resource.
func (r *backendResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages backend configuration.",
		MarkdownDescription: "Provides a Backend resource. This allows backends to be created, updated and deleted.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "ID of the backend.",
				MarkdownDescription: "ID of the backend.",
			},
			"company": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "Company ID that owns this backend.",
				MarkdownDescription: "Company ID that owns this backend.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(regexp.MustCompile("[^0-9][a-z0-9]+"), "The name must contain only lower case letters and numbers (cannot start with a number)"),
					),
				},
				Description:         "Name of the backend.",
				MarkdownDescription: "Name of the backend.",
			},
			"vclname": schema.StringAttribute{
				Computed:            true,
				Description:         "Final unique name of the backend to be referenced in VCL Code: 'c{company_id}_{name}'.",
				MarkdownDescription: "Final unique name of the backend to be referenced in VCL Code: `c{company_id}_{name}`.",
			},
			"origin": schema.StringAttribute{
				Required:            true,
				Description:         "IP or DNS name pointing to the origin backend, for example: 'my-origin.com'.",
				MarkdownDescription: "IP or DNS name pointing to the origin backend, for example: `my-origin.com`.",
			},
			"ssl": schema.BoolAttribute{
				Required:            true,
				Description:         "Use TLS encryption when contacting with the origin backend.",
				MarkdownDescription: "Use TLS encryption when contacting with the origin backend.",
			},
			"port": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(80, 65535),
				},
				Description:         "Port where the origin is listening to HTTP requests, for example: 80 or 443.",
				MarkdownDescription: "Port where the origin is listening to HTTP requests, for example: `80` or `443`.",
			},
			"headers": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.RegexMatches(
						regexp.MustCompile(`(?im)^(([a-z]\w*)[ ]*:[ ]*([^":]+))$`), "Extra headers must be in the format 'Key_1: Value_1\nKey_2: Value_2\n...\nKey_n: Value_n'"),
					),
				},
				Default:  stringdefault.StaticString(""),
				Description:         "Extra headers needed in order to validate backend status.",
				MarkdownDescription: "Extra headers needed in order to validate backend status.",
			},
			"hchost": schema.StringAttribute{
				Required:            true,
				Description:         "Host header that the health check probe will send to the origin, for example: www.my-origin.com.",
				MarkdownDescription: "Host header that the health check probe will send to the origin, for example: `www.my-origin.com`.",
			},
			"hcpath": schema.StringAttribute{
				Required:            true,
				Description:         "Path that the health check probe will use, for example: /favicon.ico.",
				MarkdownDescription: "Path that the health check probe will use, for example: `/favicon.ico`.",
			},
			"hcstatuscode": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(200, 599),
				},
				Description:         "Status code expected when the probe receives the HTTP health check response, for example: 200.",
				MarkdownDescription: "Status code expected when the probe receives the HTTP health check response, for example: `200`.",
			},
			"hcinterval": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(15, 300),
				},
				Default:             int64default.StaticInt64(40),
				Description:         "Interval in seconds within which the probes of each edge execute the HTTP request to validate the status of the backend.",
				MarkdownDescription: "Interval in seconds within which the probes of each edge execute the HTTP request to validate the status of the backend.",
			},
			"hcdisabled": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Description:         "Disable the health check probe.",
				MarkdownDescription: "Disable the health check probe.",
			},
		},
	}
}

// Create
func (r *backendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan Backend
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
		Headers:      plan.Headers.ValueString(),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
		HCInterval:   int(plan.HCInterval.ValueInt64()),
		HCDisabled:   plan.HCDisabled.ValueBool(),
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
	plan.VclName = types.StringValue("c" + strconv.Itoa(backendState.Company) + "_" + backendState.Name)
	plan.Origin = types.StringValue(backendState.Origin)
	plan.Ssl = types.BoolValue(backendState.Ssl)
	plan.Port = types.Int64Value(int64(backendState.Port))
	plan.Headers = types.StringValue(backendState.Headers)
	plan.HCHost = types.StringValue(backendState.HCHost)
	plan.HCPath = types.StringValue(backendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(backendState.HCStatusCode))
	plan.HCInterval = types.Int64Value(int64(backendState.HCInterval))
	plan.HCDisabled = types.BoolValue(backendState.HCDisabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *backendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Backend
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
		Headers:      plan.Headers.ValueString(),
		HCHost:       plan.HCHost.ValueString(),
		HCPath:       plan.HCPath.ValueString(),
		HCStatusCode: int(plan.HCStatusCode.ValueInt64()),
		HCInterval:   int(plan.HCInterval.ValueInt64()),
		HCDisabled:   plan.HCDisabled.ValueBool(),
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
	plan.VclName = types.StringValue("c" + strconv.Itoa(backendState.Company) + "_" + backendState.Name)
	plan.Origin = types.StringValue(backendState.Origin)
	plan.Ssl = types.BoolValue(backendState.Ssl)
	plan.Port = types.Int64Value(int64(backendState.Port))
	plan.Headers = types.StringValue(backendState.Headers)
	plan.HCHost = types.StringValue(backendState.HCHost)
	plan.HCPath = types.StringValue(backendState.HCPath)
	plan.HCStatusCode = types.Int64Value(int64(backendState.HCStatusCode))
	plan.HCInterval = types.Int64Value(int64(backendState.HCInterval))
	plan.HCDisabled = types.BoolValue(backendState.HCDisabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state Backend
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
			state.VclName = types.StringValue("c" + strconv.Itoa(backend.Company) + "_" + backend.Name)
			state.Origin = types.StringValue(backend.Origin)
			state.Ssl = types.BoolValue(backend.Ssl)
			state.Port = types.Int64Value(int64(backend.Port))
			state.Headers = types.StringValue(backend.Headers)
			state.HCHost = types.StringValue(backend.HCHost)
			state.HCPath = types.StringValue(backend.HCPath)
			state.HCStatusCode = types.Int64Value(int64(backend.HCStatusCode))
			state.HCInterval = types.Int64Value(int64(backend.HCInterval))
			state.HCDisabled = types.BoolValue(backend.HCDisabled)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Try to find by Name
	backend, err := r.client.GetBackendByName(state.Name.ValueString(), teclient.ProdEnv)
	if err == nil {
		if backend.Name == state.Name.ValueString() {
			state.ID = types.Int64Value(int64(backend.ID))
			state.Company = types.Int64Value(int64(backend.Company))
			state.Name = types.StringValue(backend.Name)
			state.VclName = types.StringValue("c" + strconv.Itoa(backend.Company) + "_" + backend.Name)
			state.Origin = types.StringValue(backend.Origin)
			state.Ssl = types.BoolValue(backend.Ssl)
			state.Port = types.Int64Value(int64(backend.Port))
			state.Headers = types.StringValue(backend.Headers)
			state.HCHost = types.StringValue(backend.HCHost)
			state.HCPath = types.StringValue(backend.HCPath)
			state.HCStatusCode = types.Int64Value(int64(backend.HCStatusCode))
			state.HCInterval = types.Int64Value(int64(backend.HCInterval))
			state.HCDisabled = types.BoolValue(backend.HCDisabled)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Not found
	resp.Diagnostics.AddError("Backend not found", "Backend '"+state.Name.ValueString()+"' doesn't exist in API")
}

// Delete
func (r *backendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Backend
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
