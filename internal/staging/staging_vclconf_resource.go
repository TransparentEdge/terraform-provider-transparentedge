package staging

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
	_ resource.Resource                = &stagingVclConfResource{}
	_ resource.ResourceWithConfigure   = &stagingVclConfResource{}
	_ resource.ResourceWithImportState = &stagingVclConfResource{}
	_ resource.ResourceWithModifyPlan  = &stagingVclConfResource{}
)

// helper function to simplify the provider implementation.
func NewStagingVclconfResource() resource.Resource {
	return &stagingVclConfResource{}
}

// resource implementation.
type stagingVclConfResource struct {
	client *teclient.Client
}

// maps schema data.
type stagingVclConfResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Company        types.Int64  `tfsdk:"company"`
	VCLCode        types.String `tfsdk:"vclcode"`
	UploadDate     types.String `tfsdk:"uploaddate"`
	ProductionDate types.String `tfsdk:"productiondate"`
	User           types.String `tfsdk:"user"`
}

// Metadata returns the resource type name.
func (r *stagingVclConfResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_staging_vclconf"
}

// Schema defines the schema for the resource.
func (r *stagingVclConfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the Staging VCL Config",
			},
			"company": schema.Int64Attribute{
				Computed:    true,
				Description: "Company ID that owns this Staging VCL Config",
			},
			"vclcode": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Verbatim of the VCL (Varnish Configuration Language) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully applied." +
					" You can know if a configuration is already in staging by running 'terraform plan' and checking the 'productiondate' field.",
			},
			"uploaddate": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the configuration was uploaded",
			},
			"productiondate": schema.StringAttribute{
				Computed:    true,
				Description: "Date when the configuration was fully applied in the CDN",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "User that created the configuration",
			},
		},
	}
}

// Create
func (r *stagingVclConfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stagingVclConfResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating new VCL configuration")
	newStagingVclconf := teclient.NewVCLConfAPIModel{
		VCLCode: plan.VCLCode.ValueString(),
	}
	stagingVclConfState, errCreate := r.client.CreateVclconf(newStagingVclconf, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating Staging VCL Configuration",
			fmt.Sprintf("Could not create the Staging VCL Configuration: %s", errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingVclConfState.ID))
	plan.Company = types.Int64Value(int64(stagingVclConfState.Company))
	// do not update the VCL Config since our API does some string modifications
	//plan.VCLCode = types.StringValue(stagingVclConfState.VCLCode)
	plan.UploadDate = types.StringValue(stagingVclConfState.UploadDate)
	plan.ProductionDate = types.StringValue(stagingVclConfState.ProductionDate)
	plan.User = types.StringValue(stagingVclConfState.CreatorUser.FirstName + " " + stagingVclConfState.CreatorUser.LastName + " <" + stagingVclConfState.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *stagingVclConfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Read resource information
func (r *stagingVclConfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state stagingVclConfResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stagingVclConf, err := r.client.GetActiveVCLConf(teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Staging VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(stagingVclConf.ID))
	state.Company = types.Int64Value(int64(stagingVclConf.Company))
	if helpers.SanitizeStringForDiff(stagingVclConf.VCLCode) != helpers.SanitizeStringForDiff(state.VCLCode.ValueString()) {
		state.VCLCode = types.StringValue(stagingVclConf.VCLCode)
	}
	state.UploadDate = types.StringValue(stagingVclConf.UploadDate)
	state.ProductionDate = types.StringValue(stagingVclConf.ProductionDate)
	state.User = types.StringValue(stagingVclConf.CreatorUser.FirstName + " " + stagingVclConf.CreatorUser.LastName + " <" + stagingVclConf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete
func (r *stagingVclConfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state stagingVclConfResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *stagingVclConfResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will only remove the resource from the Terraform state.\n"+
				"It will not call the API for deletion since VCL configurations cannot be deleted.",
		)
	}
}

// Configure adds the provider configured client to the resource.
func (r *stagingVclConfResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *stagingVclConfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// path.Root here is ignored
	// VCL Configs can be imported without issues, but they won't match perfectly
	// the configuration because of newlines and spaces
	resource.ImportStatePassthroughID(ctx, path.Root("user"), req, resp)
}
