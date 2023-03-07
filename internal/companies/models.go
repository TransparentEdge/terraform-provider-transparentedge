package companies

import "github.com/hashicorp/terraform-plugin-framework/types"

type IPRanges struct {
	Ipv4CidrBlocks types.List `tfsdk:"ipv4_cidr_blocks"`
	Ipv6CidrBlocks types.List `tfsdk:"ipv6_cidr_blocks"`
}
