package autoprovisioning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/customtypes"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vclconfResource{}
	_ resource.ResourceWithConfigure   = &vclconfResource{}
	_ resource.ResourceWithImportState = &vclconfResource{}
	_ resource.ResourceWithModifyPlan  = &vclconfResource{}
)

// NewVclconfResource is a helper function to simplify the provider implementation.
func NewVclconfResource() resource.Resource {
	return &vclconfResource{}
}

// resource implementation.
type vclconfResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (*vclconfResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vclconf"
}

// Schema defines the schema for the resource.
func (*vclconfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages VCL Configuration.",
		MarkdownDescription: "Provides VCL Configuration resource. This allows to generate a new VCL configuration that replaces the current one." +
			" Changing `vclcode` or `comment` uploads a new configuration version in place (no destroy/recreate)." +
			" Destroying the resource uploads an empty VCL configuration so that any backends referenced by the current code can be removed afterwards.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "ID of the VCL Config.",
				MarkdownDescription: "ID of the VCL Config.",
			},
			"company": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description:         "Company ID that owns this VCL config.",
				MarkdownDescription: "Company ID that owns this VCL config.",
			},
			"vclcode": schema.StringAttribute{
				Required:   true,
				CustomType: customtypes.VCLCodeType{},
				Description: "Verbatim of the VCL (Varnish Configuration Language) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully applied." +
					" You can know if a configuration is already in production by running 'terraform plan' and checking the 'productiondate' field.",
				MarkdownDescription: "Verbatim of the VCL (_Varnish Configuration Language_) code configuration to apply." +
					" After a successful code upload, it may take between 5 and 10 minutes for the new configuration to be fully replicated in all the CDN edge nodes." +
					" You can check if a configuration is already in production by running `terraform plan` and checking the `productiondate` field.",
			},
			"uploaddate": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Date when the configuration was uploaded.",
				MarkdownDescription: "Date when the configuration was uploaded.",
			},
			"productiondate": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Date when the configuration was fully applied in the CDN.",
				MarkdownDescription: "Date when the configuration was fully applied in the CDN.",
			},
			"user": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "User that created the configuration.",
				MarkdownDescription: "User that created the configuration.",
			},
			"comment": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Description:         "Optional comment describing the changes introduced by this configuration.",
				MarkdownDescription: "Optional comment describing the changes introduced by this configuration.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				CreateDescription: "If set, the provider will wait until the VCL configuration is fully deployed " +
					"across all CDN edge nodes before completing. Must be either null (don't wait) or a duration " +
					"greater than 5m, since propagation typically takes between 5 and 10 minutes (e.g. \"15m\").",
			}),
		},
	}
}

// Create.
func (r *vclconfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan VCLConf

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating new VCL configuration")

	r.pushVCLConf(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update pushes a new VCL configuration version whenever vclcode or comment actually
// change. If only client-side attributes (e.g. timeouts) changed, no API call is made.
func (r *vclconfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan VCLConf

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if helpers.VCLSemanticEquals(state.VCLCode.ValueString(), plan.VCLCode.ValueString()) &&
		state.Comment.ValueString() == plan.Comment.ValueString() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

		return
	}

	tflog.Info(ctx, "Updating VCL configuration by uploading a new version")

	r.pushVCLConf(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information.
func (r *vclconfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state VCLConf

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.GetActiveVCLConf(apiEnv)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read VclConf info",
			err.Error(),
		)

		return
	}

	state.ID = types.Int64Value(int64(apiResp.ID))
	state.Company = types.Int64Value(int64(apiResp.Company))
	state.VCLCode = customtypes.NewVCLCodeValue(apiResp.VCLCode)
	state.UploadDate = types.StringValue(apiResp.UploadDate)
	state.ProductionDate = types.StringValue(apiResp.ProductionDate)
	state.User = types.StringValue(apiResp.CreatorUser.FirstName + " " + apiResp.CreatorUser.LastName + " <" + apiResp.CreatorUser.Email + ">")
	// Comment is only ever set through Create/Update and the suffix appended there is stripped
	// by the API client, so syncing it here always matches, including right after an Import.
	state.Comment = types.StringValue(apiResp.Comment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete uploads an empty VCL configuration so any backends referenced by the
// current code stop being referenced and can then be deleted through the API.
func (r *vclconfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VCLConf

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Emptying VCL configuration so any referenced backends can be deleted")

	emptyConf := teclient.NewVCLConfAPIModel{
		VCLCode: helpers.EmptyVCLCode,
		Comment: "Emptied by 'terraform destroy'",
	}

	_, errCreate := r.client.CreateVclconf(emptyConf, apiEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error emptying VCL Configuration",
			fmt.Sprintf("Could not upload an empty VCL configuration: %s", errCreate),
		)
	}
}

// ModifyPlan marks the computed attributes as unknown whenever a new VCL configuration
// version is going to be uploaded (i.e. vclcode or comment change), since the API always
// assigns fresh values (id, dates, ...) to every uploaded version.
func (*vclconfResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will upload an empty VCL configuration as the new active version, "+
				"so that any backends referenced by the current VCL code can be deleted afterwards.\n"+
				"Previous VCL configuration history entries are never removed from the API.",
		)

		return
	}

	// Nothing to compare against yet, this is a resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	var state, plan VCLConf

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	unchanged := !plan.VCLCode.IsUnknown() && !plan.Comment.IsUnknown() &&
		helpers.VCLSemanticEquals(state.VCLCode.ValueString(), plan.VCLCode.ValueString()) &&
		state.Comment.ValueString() == plan.Comment.ValueString()

	if unchanged {
		tflog.Info(ctx, "ModifyPlan is unchanged")

		return
	}

	tflog.Info(ctx, "Mark computed attributes as unknown")

	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("company"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("uploaddate"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("productiondate"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("user"), types.StringUnknown())...)
}

// Configure adds the provider configured client to the resource.
func (r *vclconfResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure", "error while configuring API client")

		return
	}

	r.client = client
}

func (*vclconfResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// path.Root here is ignored
	// VCL Configs can be imported without issues, but they won't match perfectly
	// the configuration because of newlines and spaces
	resource.ImportStatePassthroughID(ctx, path.Root("user"), req, resp)
}

// pushVCLConf uploads plan's VCLCode/Comment as a new VCL configuration version and
// populates the computed attributes with the API response. Used by both Create and
// Update, since every upload produces a brand new history entry in the API.
func (r *vclconfResource) pushVCLConf(ctx context.Context, plan *VCLConf, diags *diag.Diagnostics) {
	newConf := teclient.NewVCLConfAPIModel{
		VCLCode: plan.VCLCode.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	apiResp, errCreate := r.client.CreateVclconf(newConf, apiEnv)
	if errCreate != nil {
		diags.AddError(
			"Error uploading Production VCL Configuration",
			fmt.Sprintf("Could not upload the vclconf: %s", errCreate),
		)

		return
	}

	plan.ID = types.Int64Value(int64(apiResp.ID))
	plan.Company = types.Int64Value(int64(apiResp.Company))
	plan.VCLCode = customtypes.NewVCLCodeValue(apiResp.VCLCode)
	plan.UploadDate = types.StringValue(apiResp.UploadDate)
	plan.ProductionDate = types.StringValue(apiResp.ProductionDate)
	plan.User = types.StringValue(apiResp.CreatorUser.FirstName + " " + apiResp.CreatorUser.LastName + " <" + apiResp.CreatorUser.Email + ">")

	createTimeout, timeoutDiags := plan.Timeouts.Create(ctx, 0)
	diags.Append(timeoutDiags...)

	if diags.HasError() || createTimeout == 0 {
		return
	}

	// Wait for the configuration to be deployed (when productiondate field is set)
	// if the user specified a value for the create timeout.
	pollCtx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

poll:
	for {
		select {
		case <-pollCtx.Done():
			diags.AddWarning(
				"Timeout waiting for VCL deployment",
				"The configuration was uploaded but did not reach production within the expected time.",
			)

			break poll

		case <-time.After(10 * time.Second):
			vclconf, err := r.client.GetVCLConfByID(apiEnv, apiResp.ID)
			if err != nil {
				continue
			}

			if vclconf.ProductionDate != "" && vclconf.ID == apiResp.ID {
				plan.ProductionDate = types.StringValue(vclconf.ProductionDate)

				break poll
			}

			tflog.Info(ctx, "VCL configuration not yet in production, waiting...")
		}
	}
}
