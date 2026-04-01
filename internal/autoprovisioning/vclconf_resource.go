package autoprovisioning

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/customtypes"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/planmodifiers"
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
func (*vclconfResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required:   true,
				CustomType: customtypes.VCLCodeType{},
				PlanModifiers: []planmodifier.String{
					planmodifiers.VCLCodeRequiresReplace(),
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
			"comment": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Optional comment describing the changes introduced by this configuration.",
				MarkdownDescription: "Optional comment describing the changes introduced by this configuration.",
			},
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

	newConf := teclient.NewVCLConfAPIModel{
		VCLCode: plan.VCLCode.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	apiResp, errCreate := r.client.CreateVclconf(newConf, teclient.ProdEnv)
	if errCreate != nil {
		resp.Diagnostics.AddError(
			"Error creating Production VCL Configuration",
			fmt.Sprintf("Could not create the vclconf: %s", errCreate),
		)

		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(apiResp.ID))
	plan.Company = types.Int64Value(int64(apiResp.Company))
	plan.VCLCode = customtypes.NewVCLCodeValue(apiResp.VCLCode)
	plan.UploadDate = types.StringValue(apiResp.UploadDate)
	plan.ProductionDate = types.StringValue(apiResp.ProductionDate)
	plan.User = types.StringValue(apiResp.CreatorUser.FirstName + " " + apiResp.CreatorUser.LastName + " <" + apiResp.CreatorUser.Email + ">")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (*vclconfResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
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

	apiResp, err := r.client.GetActiveVCLConf(teclient.ProdEnv)
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

	// Do not update comment on Read, its just metadata, this also ensures compatibility between provider versions that did not have the comment field
	// only set if its null to set a default value.
	if state.Comment.IsNull() || state.Comment.IsUnknown() {
		state.Comment = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete.
func (*vclconfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VCLConf

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (*vclconfResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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
