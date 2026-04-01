package customtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/helpers"
)

type VCLCodeValue struct {
	basetypes.StringValue
}

var (
	_ basetypes.StringValuable                   = VCLCodeValue{}
	_ basetypes.StringValuableWithSemanticEquals = VCLCodeValue{}
)

func NewVCLCodeValue(value string) VCLCodeValue {
	return VCLCodeValue{StringValue: basetypes.NewStringValue(value)}
}

func NewVCLCodeNull() VCLCodeValue {
	return VCLCodeValue{StringValue: basetypes.NewStringNull()}
}

func NewVCLCodeUnknown() VCLCodeValue {
	return VCLCodeValue{StringValue: basetypes.NewStringUnknown()}
}

func (v VCLCodeValue) Equal(o attr.Value) bool {
	other, ok := o.(VCLCodeValue)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (VCLCodeValue) Type(_ context.Context) attr.Type {
	return VCLCodeType{}
}

func (v VCLCodeValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(VCLCodeValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected VCLCodeValue, got: "+newValuable.String(),
		)

		return false, diags
	}

	// VCL is semantically equal regardless of newlines or blank spaces.
	eq := helpers.VCLSemanticEquals(v.ValueString(), newValue.ValueString())

	return eq, diags
}
