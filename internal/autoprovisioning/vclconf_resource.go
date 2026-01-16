package autoprovisioning

import (
	"context"
	"fmt"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vclconfResource{}
	_ resource.ResourceWithConfigure   = &vclconfResource{}
	_ resource.ResourceWithImportState = &vclconfResource{}
)

// helper function to simplify the provider implementation.
func NewVclconfResource() resource.Resource {
	return &vclconfResource{}
}

// resource implementation.
type vclconfResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (r *vclconfResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vclconf"
}

// Schema defines the schema for the resource.
func (r *vclconfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages VCL Configuration.",
		MarkdownDescription: "Provides VCL Configuration resource. This allows to generate a new VCL configuration that replaces the current one.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the VCL Config.",
				MarkdownDescription: "ID of the VCL Config.",
			},
			"company": schema.Int64Attribute{
				Computed:            true,
				Description:         "Company ID that owns this VCL config.",
				MarkdownDescription: "Company ID that owns this VCL config.",
			},
			"vclcode": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Verbatim of the VCL (Varnish Configuration Language) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully applied." +
					" You can know if a configuration is already in production by running 'terraform plan' and checking the 'productiondate' field.",
				MarkdownDescription: "Verbatim of the VCL (_Varnish Configuration Language_) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully replicated in all the CDN edge nodes." +
					" You can check if a configuration is already in production by running `terraform plan` and checking the `productiondate` field.",
			},
			"uploaddate": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the configuration was uploaded.",
				MarkdownDescription: "Date when the configuration was uploaded.",
			},
			"productiondate": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the configuration was fully applied in the CDN.",
				MarkdownDescription: "Date when the configuration was fully applied in the CDN.",
			},
			"user": schema.StringAttribute{
				Computed:            true,
				Description:         "User that created the configuration.",
				MarkdownDescription: "User that created the configuration.",
			},
		},
	}
}

// Create
func (r *vclconfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan VCLConf
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating new VCL configuration")
	newVclconf := teclient.NewVCLConfAPIModel{
		VCLCode: plan.VCLCode.ValueString(),
	}
	vclconfState, errCreate := r.client.CreateVclconf(newVclconf, teclient.ProdEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating Production VCL Configuration",
			fmt.Sprintf("Could not create the vclconf: %s", errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(vclconfState.ID))
	plan.Company = types.Int64Value(int64(vclconfState.Company))
	// do not update the VCL Config since our API does some string modifications
	// plan.VCLCode = types.StringValue(vclconfState.VCLCode)
	plan.UploadDate = types.StringValue(vclconfState.UploadDate)
	plan.ProductionDate = types.StringValue(vclconfState.ProductionDate)
	plan.User = types.StringValue(vclconfState.CreatorUser.FirstName + " " + vclconfState.CreatorUser.LastName + " <" + vclconfState.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vclconfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Read resource information
func (r *vclconfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state VCLConf
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vclconf, err := r.client.GetActiveVCLConf(teclient.ProdEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(vclconf.ID))
	state.Company = types.Int64Value(int64(vclconf.Company))
	if helpers.SanitizeStringForDiff(vclconf.VCLCode) != helpers.SanitizeStringForDiff(state.VCLCode.ValueString()) {
		state.VCLCode = types.StringValue(vclconf.VCLCode)
	}
	state.UploadDate = types.StringValue(vclconf.UploadDate)
	state.ProductionDate = types.StringValue(vclconf.ProductionDate)
	state.User = types.StringValue(vclconf.CreatorUser.FirstName + " " + vclconf.CreatorUser.LastName + " <" + vclconf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete
func (r *vclconfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VCLConf
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting VCL configuration by creating placeholder config")

	placeholderVclconf := teclient.NewVCLConfAPIModel{
		VCLCode: `sub vcl_recv { set req.http.placeholder = "Modified by 'terraform destroy'"; }`,
	}

	_, errCreate := r.client.CreateVclconf(placeholderVclconf, teclient.ProdEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error deleting VCL Configuration",
			fmt.Sprintf("Could not create placeholder VCL configuration: %s", errCreate),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted VCL configuration with an empty placeholder")
}

// Configure adds the provider configured client to the resource.
func (r *vclconfResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *vclconfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// path.Root here is ignored
	// VCL Configs can be imported without issues, but they won't match perfectly
	// the configuration because of newlines and spaces
	resource.ImportStatePassthroughID(ctx, path.Root("user"), req, resp)
}
