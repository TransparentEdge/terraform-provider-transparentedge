package companies

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

var (
	_ datasource.DataSource              = &ipRangesDataSource{}
	_ datasource.DataSourceWithConfigure = &ipRangesDataSource{}
)

func NewIPRangesDataSource() datasource.DataSource {
	return &ipRangesDataSource{}
}

type ipRangesDataSource struct {
	client *teclient.Client
}

func (*ipRangesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_ranges"
}

func (*ipRangesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *ipRangesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state IPRanges

	cidrRanges, err := d.client.GetIPRanges()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IP Ranges",
			fmt.Sprintf("Unexpected error trying to read ip ranges.\n%s\n", err.Error()),
		)

		return
	}

	ipv4List := []string{}
	ipv6List := []string{}

	for _, ip := range cidrRanges {
		if strings.Contains(ip, ":") {
			ipv6List = append(ipv6List, ip)
		} else {
			ipv4List = append(ipv4List, ip)
		}
	}

	slices.Sort(ipv4List)
	slices.Sort(ipv6List)

	// Map response body to model
	ipv4CIDRRanges, diag := types.ListValueFrom(ctx, types.StringType, ipv4List)
	if diag != nil {
		resp.Diagnostics.Append(diag...)

		return
	}

	ipv6CIDRRanges, diag := types.ListValueFrom(ctx, types.StringType, ipv6List)
	if diag != nil {
		resp.Diagnostics.Append(diag...)

		return
	}

	state.Ipv4CidrBlocks = ipv4CIDRRanges
	state.Ipv6CidrBlocks = ipv6CIDRRanges

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *ipRangesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure", "error while configuring API client")

		return
	}

	d.client = client
}
