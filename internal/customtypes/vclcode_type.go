package customtypes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// VCLCodeType is a custom type to handle VCL code differences between user state and API responses.
// Implemented following: https://developer.hashicorp.com/terraform/plugin/framework/handling-data/types/custom#developing-custom-types.
type VCLCodeType struct {
	basetypes.StringType
}

var _ basetypes.StringTypable = VCLCodeType{}

func (t VCLCodeType) Equal(o attr.Type) bool {
	other, ok := o.(VCLCodeType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (VCLCodeType) String() string {
	return "VCLCodeType"
}

func (VCLCodeType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return VCLCodeValue{StringValue: in}, nil
}

func (t VCLCodeType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (VCLCodeType) ValueType(_ context.Context) attr.Value {
	return VCLCodeValue{}
}
