package autoprovisioning

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

const (
	certreqHTTPCreateTimeout time.Duration = 10 * time.Minute
	certreqHTTPCreateRetry   time.Duration = 20 * time.Second
)

type certreqHTTPResource struct {
	client *teclient.Client
}

var (
	_ resource.Resource                = &certreqHTTPResource{}
	_ resource.ResourceWithConfigure   = &certreqHTTPResource{}
	_ resource.ResourceWithImportState = &certreqHTTPResource{}
	_ resource.ResourceWithModifyPlan  = &certreqHTTPResource{}
)

func NewCertReqHTTPResource() resource.Resource {
	return &certreqHTTPResource{}
}

func (r *certreqHTTPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_http"
}

func (r *certreqHTTPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages HTTP Certificate Requests.",
		MarkdownDescription: `Manages HTTP Certificate Requests.

This resource enables the creation of certificate requests using the HTTP challenge. This challenge requires that the domain's DNS already points to the CDN.

For detailed documentation (not Terraform-specific), please refer to this [link](https://docs.transparentedge.eu/getting-started/dashboard/auto-provisioning/ssl).`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the HTTP Certificate Request.",
				MarkdownDescription: "ID of the HTTP Certificate Request.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domains": schema.SetAttribute{
				Required:            true,
				Description:         "List of domains for which you want to request a certificate. You can not include wildcard domains, such as `*.example.com`, use DNS Certificate Requests instead.",
				MarkdownDescription: "List of domains for which you want to request a certificate. You can **not** include wildcard domains, such as `*.example.com`, use DNS Certificate Requests instead.",
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 200),
				},
			},
			"certificate_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
				MarkdownDescription: "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
			},
			"standalone": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Description:         "When set to `true`, this indicates that the certificate's domains should be treated as standalone and not merged into an existing certificate, either immediately or during future renewals.",
				MarkdownDescription: "When set to `true`, this indicates that the certificate's domains should be treated as standalone and not merged into an existing certificate, either immediately or during future renewals.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date of creation.",
				MarkdownDescription: "Date of creation.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date of last update.",
				MarkdownDescription: "Date of last update.",
			},
			"status_message": schema.StringAttribute{
				Computed:            true,
				Description:         "Indicates the current status message for the certificate request. This field will display a success message if the certificate is obtained successfully or an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
				MarkdownDescription: "Indicates the current status message for the certificate request. This field will display a success message if the certificate is obtained successfully or an error message if the request fails. When the certificate request is in progress or if there is no status message available, it will be represented as `null`.",
			},
		},
	}
}

func (r *certreqHTTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve current Plan
	var plan CertReqHTTP
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the Certificate Request in the API
	domains := make([]string, 0, len(plan.Domains.Elements()))
	diags = plan.Domains.ElementsAs(ctx, &domains, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	standalone := false
	if !plan.Standalone.IsNull() && plan.Standalone.ValueBool() {
		standalone = true
	}
	api_model, err := r.client.CreateHTTPCertReq(map[string]interface{}{
		"domains":    domains,
		"standalone": standalone,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP Certificate Request",
			err.Error(),
		)
		return
	}

	// Wait until the CR is complete
	updated_api_model := r.WaitCompleteHTTPCertReq(ctx, api_model)
	if updated_api_model != nil {
		if updated_api_model.CertificateID != nil {
			// Certificate was generated
			api_model = updated_api_model
		} else if updated_api_model.Log != nil && *updated_api_model.Log != "" {
			// An error happened
			resp.Diagnostics.AddError(
				"Error generating the certificate",
				helpers.ParseCertReqLogString(*updated_api_model.Log),
			)
			return
		}
	}

	// Save API response in the state

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.CommonName + "\n" + api_model.SAN)
	new_domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	plan.ID = types.Int64Value(int64(api_model.ID))
	plan.Domains = new_domains
	plan.Standalone = types.BoolValue(api_model.Standalone)
	plan.CreatedAt = types.StringValue(api_model.CreatedAt)
	plan.UpdatedAt = types.StringValue(api_model.UpdatedAt)
	if api_model.CertificateID == nil {
		plan.CertificateID = types.Int64Null()
	} else {
		plan.CertificateID = types.Int64Value(int64(*api_model.CertificateID))
	}
	if api_model.Log == nil {
		plan.StatusMessage = types.StringNull()
	} else {
		plan.StatusMessage = types.StringValue(helpers.ParseCertReqLogString(*api_model.Log))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *certreqHTTPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CertReqHTTP
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current
	api_model, err := r.client.GetCertReqHTTP(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving HTTP Certificate Request",
			err.Error(),
		)
		return
	}

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.CommonName + "\n" + api_model.SAN)
	domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(api_model.ID))
	state.Domains = domains
	state.Standalone = types.BoolValue(api_model.Standalone)
	state.CreatedAt = types.StringValue(api_model.CreatedAt)
	state.UpdatedAt = types.StringValue(api_model.UpdatedAt)
	if api_model.CertificateID == nil {
		state.CertificateID = types.Int64Null()
	} else {
		state.CertificateID = types.Int64Value(int64(*api_model.CertificateID))
	}
	if api_model.Log == nil {
		state.StatusMessage = types.StringNull()
	} else {
		state.StatusMessage = types.StringValue(helpers.ParseCertReqLogString(*api_model.Log))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *certreqHTTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *certreqHTTPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *certreqHTTPResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will only remove the resource from the Terraform state.\n"+
				"It will not call the API for deletion since HTTP Certificate Requests cannot be deleted.",
		)
	}
}

func (r *certreqHTTPResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *certreqHTTPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid identifier", "ID must be a valid number.")
		return
	}
	if id <= 0 {
		resp.Diagnostics.AddError("Invalid identifier", "ID must be a valid number greater than 0.")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *certreqHTTPResource) WaitCompleteHTTPCertReq(ctx context.Context, certreq *teclient.CertReqHTTP) *teclient.CertReqHTTP {
	// Read the state of the HTTP Certificate Request until it's processed or timeout
	remainingTime := certreqHTTPCreateTimeout.Seconds()

	for {
		api_model, err := r.client.GetCertReqHTTP(certreq.ID)
		if err == nil {
			if api_model.CertificateID != nil || (api_model.Log != nil && *api_model.Log != "") {
				return &api_model
			}
		}

		if remainingTime <= 0 {
			break
		}
		time.Sleep(certreqHTTPCreateRetry)
		remainingTime -= (certreqHTTPCreateRetry.Seconds() + 1)
		tflog.Info(ctx, fmt.Sprintf("Waiting for the HTTP Certificate Request %d to be completed.", certreq.ID))
	}

	return nil
}
