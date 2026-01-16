package autoprovisioning

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &crDNSCredentialResource{}
	_ resource.ResourceWithConfigure   = &crDNSCredentialResource{}
	_ resource.ResourceWithImportState = &crDNSCredentialResource{}
)

// helper function to simplify the provider implementation.
func NewCertReqDNSCredentialResource() resource.Resource {
	return &crDNSCredentialResource{}
}

// resource implementation.
type crDNSCredentialResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (r *crDNSCredentialResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certreq_dns_credential"
}

// Schema defines the schema for the resource.
func (r *crDNSCredentialResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Generates DNS Credentials.",
		MarkdownDescription: "Provides DNS Credential resource. This allows to create, update and delete DNS Credentials used in [DNS Certificate Requests](https://docs.transparentedge.eu/getting-started/dashboard/auto-provisioning/ssl).",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the DNS Credential.",
				MarkdownDescription: "ID of the DNS Credential.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"alias": schema.StringAttribute{
				Required:            true,
				Description:         "Alias for the DNS Credential.",
				MarkdownDescription: "Alias for the DNS Credential.",
			},
			"dns_provider": schema.StringAttribute{
				Computed:            true,
				Description:         "DNS Provider.",
				MarkdownDescription: "DNS Provider.",
			},
			"parameters": schema.MapAttribute{
				Required:            true,
				Description:         "Keys/parameters of the provider.",
				MarkdownDescription: "Keys/parameters of the provider.",
				Sensitive:           true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Create
func (r *crDNSCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CertReqDNSCredential
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if parameters are set
	if plan.Parameters.IsNull() {
		resp.Diagnostics.AddError(
			"Parameters cannot be empty.",
			"Please provide the parameters required for the DNS provider of your choice.",
		)
		return
	}

	// Read the parameters from the Terraform plan and transform to API struct
	parameters := make(map[string]string, len(plan.Parameters.Elements()))
	diags = plan.Parameters.ElementsAs(ctx, &parameters, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var creds []teclient.NewCRDNSCreds
	for key, value := range parameters {
		creds = append(creds, teclient.NewCRDNSCreds{KeyName: key, KeyValue: value})
	}

	new_data := teclient.NewCRDNSCredential{
		Alias: plan.Alias.ValueString(),
		Creds: creds,
	}

	new_state, err := r.client.CreateDNSCredential(new_data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating the credential",
			err.Error(),
		)
		return
	}

	// Extract the parameters/keys obtained from the API into a map
	new_keys := make(map[string]attr.Value)
	dns_provider := "Unknown"
	for _, key := range new_state.Creds {
		new_keys[key.KeyName] = types.StringValue(key.KeyValue)
		dns_provider = key.Provider
	}

	// Transform the map into a Terraform type
	new_parameters, diags := types.MapValue(types.StringType, new_keys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(new_state.ID))
	plan.Alias = types.StringValue(new_state.Alias)
	plan.Parameters = new_parameters
	plan.DNSProvider = types.StringValue(dns_provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *crDNSCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CertReqDNSCredential
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if parameters are set
	if plan.Parameters.IsNull() {
		resp.Diagnostics.AddError(
			"Parameters cannot be empty.",
			"Please provide the parameters required for the DNS provider of your choice.",
		)
		return
	}

	// Read the parameters from the Terraform plan and transform to API struct
	parameters := make(map[string]string, len(plan.Parameters.Elements()))
	diags = plan.Parameters.ElementsAs(ctx, &parameters, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var creds []teclient.NewCRDNSCreds
	for key, value := range parameters {
		creds = append(creds, teclient.NewCRDNSCreds{KeyName: key, KeyValue: value})
	}

	new_data := teclient.NewCRDNSCredential{
		Alias: plan.Alias.ValueString(),
		Creds: creds,
	}

	new_state, err := r.client.UpdateDNSCredential(new_data, int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating the credential",
			err.Error(),
		)
		return
	}

	// Extract the parameters/keys obtained from the API into a map
	new_keys := make(map[string]attr.Value)
	dns_provider := "Unknown"
	for _, key := range new_state.Creds {
		new_keys[key.KeyName] = types.StringValue(key.KeyValue)
		dns_provider = key.Provider
	}

	// Transform the map into a Terraform type
	new_parameters, diags := types.MapValue(types.StringType, new_keys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	plan.ID = types.Int64Value(int64(new_state.ID))
	plan.Alias = types.StringValue(new_state.Alias)
	plan.Parameters = new_parameters
	plan.DNSProvider = types.StringValue(dns_provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information
func (r *crDNSCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CertReqDNSCredential
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credential, err := r.client.GetCRDNSCredential(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failure retrieving the credential from the API",
			err.Error(),
		)
		return
	}

	// Extract the parameters/keys obtained from the API into a map
	new_keys := make(map[string]attr.Value)
	dns_provider := "Unknown"
	for _, key := range credential.Creds {
		new_keys[key.KeyName] = types.StringValue(key.KeyValue)
		dns_provider = key.Provider
	}

	// Transform the map into a Terraform type
	new_parameters, diags := types.MapValue(types.StringType, new_keys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	state.ID = types.Int64Value(int64(credential.ID))
	state.Alias = types.StringValue(credential.Alias)
	state.Parameters = new_parameters
	state.DNSProvider = types.StringValue(dns_provider)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete
func (r *crDNSCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CertReqDNSCredential
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 204 on successful delete
	if err := r.client.DeleteCRCredential(int(state.ID.ValueInt64())); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting the credential",
			"Could not delete the credential with id: "+state.ID.String()+"\n"+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *crDNSCredentialResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*teclient.Client)
}

func (r *crDNSCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
