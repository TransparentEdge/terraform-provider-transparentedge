package autoprovisioning

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	certreqDNSCreateTimeout time.Duration = 12 * time.Minute
	certreqDNSCreateRetry   time.Duration = 20 * time.Second
)

type certreqDNSResource struct {
	client *teclient.Client
}

var (
	_ resource.Resource                = &certreqDNSResource{}
	_ resource.ResourceWithConfigure   = &certreqDNSResource{}
	_ resource.ResourceWithImportState = &certreqDNSResource{}
)

func NewCertReqDNSResource() resource.Resource {
	return &certreqDNSResource{}
}

func (r *certreqDNSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_dns"
}

func (r *certreqDNSResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages DNS Certificate Requests.",
		MarkdownDescription: `Manages DNS Certificate Requests.

This resource enables the creation of certificate requests using various challenges, including:
- DNS Challenge
- DNS Challenge by CNAME

For detailed documentation (not Terraform-specific), please refer to this [link](https://docs.transparentedge.eu/getting-started/dashboard/auto-provisioning/ssl).`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the DNS Certificate Request.",
				MarkdownDescription: "ID of the DNS Certificate Request.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"domains": schema.SetAttribute{
				Required:            true,
				Description:         "List of domains for which you want to request a certificate. You can include wildcard domains, such as `*.example.com`, to cover subdomains under a common domain.",
				MarkdownDescription: "List of domains for which you want to request a certificate. You can include wildcard domains, such as `*.example.com`, to cover subdomains under a common domain.",
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeBetween(1, 200),
				},
			},
			"credential": schema.Int64Attribute{
				Required:            true,
				Description:         "DNS Credential associated.",
				MarkdownDescription: "DNS Credential associated.",
			},
			"certificate_id": schema.Int64Attribute{
				Computed:            true,
				Description:         "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
				MarkdownDescription: "Certificate ID. It will be `null` in case of failure or when the certificate request is in progress.",
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

func (r *certreqDNSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve current Plan
	var plan CertReqDNS
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
	api_model, err := r.client.CreateDNSCertReq(map[string]interface{}{
		"domains":               strings.Join(domains, "\n"),
		"credential":            plan.Credential.ValueInt64(),
		"certificate_authority": 1, // Let's Encrypt
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DNS Certificate Request",
			err.Error(),
		)
		return
	}

	// Wait until the CR is complete
	updated_api_model := r.WaitCompleteDNSCertReq(ctx, api_model)
	if updated_api_model != nil {
		api_model = updated_api_model
	}

	// Save API response in the state

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.Domains)
	new_domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	plan.ID = types.Int64Value(int64(api_model.ID))
	plan.Domains = new_domains
	plan.Credential = types.Int64Value(int64(api_model.Credential))
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

func (r *certreqDNSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CertReqDNS
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only the credential can be updated. Modifying the domains requires replace.
	err := r.client.UpdateDNSCertReq(int(plan.ID.ValueInt64()), int(plan.Credential.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating DNS Certificate Request",
			err.Error(),
		)
		return
	}

	// Get current
	api_model, err := r.client.GetCertReqDNS(int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving DNS Certificate Request",
			err.Error(),
		)
		return
	}

	// Save API response in the state

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.Domains)
	new_domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	plan.ID = types.Int64Value(int64(api_model.ID))
	plan.Domains = new_domains
	plan.Credential = types.Int64Value(int64(api_model.Credential))
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

func (r *certreqDNSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CertReqDNS
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current
	api_model, err := r.client.GetCertReqDNS(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving DNS Certificate Request",
			err.Error(),
		)
		return
	}

	// Generate the list of domains
	sorted_domains := helpers.SplitAndSort(api_model.Domains)
	domains, diags := types.SetValueFrom(ctx, types.StringType, sorted_domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state
	state.ID = types.Int64Value(int64(api_model.ID))
	state.Domains = domains
	state.Credential = types.Int64Value(int64(api_model.Credential))
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

func (r *certreqDNSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state CertReqDNS
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the resource
	if err := r.client.DeteleDNSCertReq(int(state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting the DNS Certificate Request",
			"Could not delete the DNS Certificate Request with id: "+state.ID.String()+"\n"+err.Error(),
		)
		return
	}
}

func (r *certreqDNSResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *certreqDNSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *certreqDNSResource) WaitCompleteDNSCertReq(ctx context.Context, certreq *teclient.CertReqDNS) *teclient.CertReqDNS {
	// Read the state of the DNS Certificate Request until it's processed or timeout
	remainingTime := certreqDNSCreateTimeout.Seconds()

	for {
		api_model, err := r.client.GetCertReqDNS(certreq.ID)
		if err == nil {
			if api_model.CertificateID != nil || api_model.Log != nil {
				return &api_model
			}
		}

		if remainingTime <= 0 {
			break
		}
		time.Sleep(certreqDNSCreateRetry)
		remainingTime -= (certreqDNSCreateRetry.Seconds() + 1)
		tflog.Info(ctx, fmt.Sprintf("Waiting for the DNS Certificate Request %d to be completed.", certreq.ID))
	}

	return nil
}
