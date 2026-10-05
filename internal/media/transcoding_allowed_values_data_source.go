package media

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &transcodingAllowedValuesDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingAllowedValuesDataSource{}
)

// NewTranscodingAllowedValuesDataSource is a helper function to simplify the provider implementation.
func NewTranscodingAllowedValuesDataSource() datasource.DataSource {
	return &transcodingAllowedValuesDataSource{}
}

// data source implementation.
type transcodingAllowedValuesDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (*transcodingAllowedValuesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_allowed_values"
}

// Schema defines the schema for the data source.
func (*transcodingAllowedValuesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allowed values for the transcoding profile enums.",
		MarkdownDescription: "Allowed values for the transcoding profile enums.",

		Attributes: map[string]schema.Attribute{
			"video_formats": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Allowed values for video_format.",
				MarkdownDescription: "Allowed values for `video_format`.",
			},
			"video_codecs": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Allowed values for video_codec.",
				MarkdownDescription: "Allowed values for `video_codec`.",
			},
			"audio_codecs": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Allowed values for audio_codec.",
				MarkdownDescription: "Allowed values for `audio_codec`.",
			},
			"video_aspect": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Allowed values for video_aspect.",
				MarkdownDescription: "Allowed values for `video_aspect`.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transcodingAllowedValuesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TranscodingAllowedValues

	values, err := d.client.GetTranscodingAllowedValues()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Transcoding Allowed Values",
			err.Error(),
		)

		return
	}

	videoFormats, diags := types.ListValueFrom(ctx, types.StringType, values.VideoFormats)
	resp.Diagnostics.Append(diags...)

	videoCodecs, diags := types.ListValueFrom(ctx, types.StringType, values.VideoCodecs)
	resp.Diagnostics.Append(diags...)

	audioCodecs, diags := types.ListValueFrom(ctx, types.StringType, values.AudioCodecs)
	resp.Diagnostics.Append(diags...)

	videoAspect, diags := types.ListValueFrom(ctx, types.StringType, values.AspectRatio)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.VideoFormats = videoFormats
	state.VideoCodecs = videoCodecs
	state.AudioCodecs = audioCodecs
	state.VideoAspect = videoAspect

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *transcodingAllowedValuesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure transcoding allowed values client", "error while configuring API client")

		return
	}

	d.client = client
}
