package companies

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &ipRangesDataSource{}
	_ datasource.DataSourceWithConfigure = &ipRangesDataSource{}
)

func NewIpRangesDataSource() datasource.DataSource {
	return &ipRangesDataSource{}
}

type ipRangesDataSource struct {
	client *teclient.Client
}

func (d *ipRangesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_ranges"
}

func (d *ipRangesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves IP Ranges from TransparentEdge nodes.",
		MarkdownDescription: "Retrieves IP Ranges from TransparentEdge nodes.",

		Attributes: map[string]schema.Attribute{
			"ipv4_cidr_blocks": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "List of IPv4 CIDR blocks.",
				MarkdownDescription: "List of IPv4 CIDR blocks.",
			},
			"ipv6_cidr_blocks": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "List of IPv6 CIDR blocks.",
				MarkdownDescription: "List of IPv6 CIDR blocks.",
			},
		},
	}
}

func (d *ipRangesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state IPRanges

	cidr_ranges, err := d.client.GetIPRanges()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IP Ranges",
			fmt.Sprintf("Unexpected error trying to read ip ranges.\n%s\n", err.Error()),
		)
		return
	}

	ipv4_list := []string{}
	ipv6_list := []string{}
	for _, ip := range cidr_ranges {
		if strings.Contains(ip, ":") {
			ipv6_list = append(ipv6_list, ip)
		} else {
			ipv4_list = append(ipv4_list, ip)
		}
	}
	sort.Strings(ipv4_list)
	sort.Strings(ipv6_list)

	// Map response body to model
	ipv4_cidr_ranges, diag := types.ListValueFrom(ctx, types.StringType, ipv4_list)
	if diag != nil {
		resp.Diagnostics.Append(diag...)
		return
	}
	ipv6_cidr_ranges, diag := types.ListValueFrom(ctx, types.StringType, ipv6_list)
	if diag != nil {
		resp.Diagnostics.Append(diag...)
		return
	}
	state.Ipv4CidrBlocks = ipv4_cidr_ranges
	state.Ipv6CidrBlocks = ipv6_cidr_ranges

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ipRangesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
