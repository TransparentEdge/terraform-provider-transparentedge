package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
)

type vclCodeRequiresReplace struct{}

func VCLCodeRequiresReplace() planmodifier.String {
	return vclCodeRequiresReplace{}
}

func (vclCodeRequiresReplace) Description(_ context.Context) string {
	return "Requires replace if VCL code changes, ignoring inconsequential whitespace differences."
}

func (vclCodeRequiresReplace) MarkdownDescription(_ context.Context) string {
	return "Requires replace if VCL code changes, ignoring inconsequential whitespace differences."
}

func (vclCodeRequiresReplace) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If state is null, this is a new resource, no replace needed.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// If plan is null, this is a destroy, no replace needed.
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Requires replace is true if configurations differ semantically.
	eq := helpers.VCLSemanticEquals(req.StateValue.ValueString(), req.PlanValue.ValueString())
	resp.RequiresReplace = !eq
}
