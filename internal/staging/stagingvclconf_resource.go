package staging

import (
	"context"
	"fmt"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &stagingvclconfResource{}
	_ resource.ResourceWithConfigure   = &stagingvclconfResource{}
	_ resource.ResourceWithImportState = &stagingvclconfResource{}
	_ resource.ResourceWithModifyPlan  = &stagingvclconfResource{}
)

// helper function to simplify the provider implementation.
func NewStagingVclconfResource() resource.Resource {
	return &stagingvclconfResource{}
}

// resource implementation.
type stagingvclconfResource struct {
	client *teclient.Client
}

// maps schema data.
type stagingvclconfResourceModel struct {
	ID             types.Int64  `tfsdk:"id"`
	Company        types.Int64  `tfsdk:"company"`
	VCLCode        types.String `tfsdk:"vclcode"`
	UploadDate     types.String `tfsdk:"uploaddate"`
	ProductionDate types.String `tfsdk:"productiondate"`
	User           types.String `tfsdk:"user"`
}

// Metadata returns the resource type name.
func (r *stagingvclconfResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stagingvclconf"
}

// Schema defines the schema for the resource.
func (r *stagingvclconfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "ID of the staging vclconf",
			},
			"company": schema.Int64Attribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Company ID that owns this staging vclconf",
			},
			"vclcode": schema.StringAttribute{
				Required: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Verbatim of the VCL (Varnish Configuration Language) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully applied." +
					" You can know if a configuration is already in staging by running 'terraform plan' and checking the 'productiondate' field.",
			},
			"uploaddate": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Date when the configuration was uploaded",
			},
			"productiondate": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Date when the configuration was fully applied in the CDN",
			},
			"user": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "User that created the configuration",
			},
		},
	}
}

// Create
func (r *stagingvclconfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stagingvclconfResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating new VCL configuration")
	newStagingVclconf := teclient.NewVCLConfAPIModel{
		VCLCode: plan.VCLCode.ValueString(),
	}
	stagingvclconfState, errCreate := r.client.CreateVclconf(newStagingVclconf, teclient.StagingEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating staging vclconf",
			fmt.Sprintf("Could not create the staging vclconf: %s", errCreate),
		)
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(stagingvclconfState.ID))
	plan.Company = types.Int64Value(int64(stagingvclconfState.Company))
	// do not update the VCL Config since our API does some string modifications
	//plan.VCLCode = types.StringValue(stagingvclconfState.VCLCode)
	plan.UploadDate = types.StringValue(stagingvclconfState.UploadDate)
	plan.ProductionDate = types.StringValue(stagingvclconfState.ProductionDate)
	plan.User = types.StringValue(stagingvclconfState.CreatorUser.FirstName + " " + stagingvclconfState.CreatorUser.LastName + " <" + stagingvclconfState.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *stagingvclconfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Read resource information
func (r *stagingvclconfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state stagingvclconfResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stagingvclconf, err := r.client.GetActiveVCLConf(teclient.StagingEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Staging VclConf info",
			err.Error(),
		)
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(stagingvclconf.ID))
	state.Company = types.Int64Value(int64(stagingvclconf.Company))
	if helpers.SanitizeStringForDiff(stagingvclconf.VCLCode) != helpers.SanitizeStringForDiff(state.VCLCode.ValueString()) {
		state.VCLCode = types.StringValue(stagingvclconf.VCLCode)
	}
	state.UploadDate = types.StringValue(stagingvclconf.UploadDate)
	state.ProductionDate = types.StringValue(stagingvclconf.ProductionDate)
	state.User = types.StringValue(stagingvclconf.CreatorUser.FirstName + " " + stagingvclconf.CreatorUser.LastName + " <" + stagingvclconf.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete
func (r *stagingvclconfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state stagingvclconfResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *stagingvclconfResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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
func (r *stagingvclconfResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *stagingvclconfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// path.Root here is ignored
	// VCL Configs can be imported without issues, but they won't match perfectly
	// the configuration because of newlines and spaces
	resource.ImportStatePassthroughID(ctx, path.Root("user"), req, resp)
}
