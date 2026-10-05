package media

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &transcodingProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingProfilesDataSource{}
)

// NewTranscodingProfilesDataSource is a helper function to simplify the provider implementation.
func NewTranscodingProfilesDataSource() datasource.DataSource {
	return &transcodingProfilesDataSource{}
}

// data source implementation.
type transcodingProfilesDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (*transcodingProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profiles"
}

// Schema defines the schema for the data source.
func (*transcodingProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Transcoding profile listing. Use it to discover the id of existing profiles (for example to generate import blocks).",
		MarkdownDescription: "Transcoding profile listing. Use it to discover the `id` of existing profiles (for example to generate `import {}` blocks).",

		Attributes: map[string]schema.Attribute{
			"profiles": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of all transcoding profiles.",
				MarkdownDescription: "List of all transcoding profiles.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "ID of the transcoding profile.",
							MarkdownDescription: "ID of the transcoding profile.",
						},
						"company": schema.Int64Attribute{
							Computed:            true,
							Description:         "Company ID that owns this transcoding profile.",
							MarkdownDescription: "Company ID that owns this transcoding profile.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the transcoding profile.",
							MarkdownDescription: "Name of the transcoding profile.",
						},
						"video_format": schema.StringAttribute{
							Computed:            true,
							Description:         "Output video container format.",
							MarkdownDescription: "Output video container format.",
						},
						"video_codec": schema.StringAttribute{
							Computed:            true,
							Description:         "Output video codec.",
							MarkdownDescription: "Output video codec.",
						},
						"video_width": schema.Int64Attribute{
							Computed:            true,
							Description:         "Output video width, in pixels.",
							MarkdownDescription: "Output video width, in pixels.",
						},
						"video_height": schema.Int64Attribute{
							Computed:            true,
							Description:         "Output video height, in pixels.",
							MarkdownDescription: "Output video height, in pixels.",
						},
						"video_bitrate": schema.Int64Attribute{
							Computed:            true,
							Description:         "Output video bitrate, in kbps.",
							MarkdownDescription: "Output video bitrate, in kbps.",
						},
						"video_aspect": schema.StringAttribute{
							Computed:            true,
							Description:         "Output video aspect ratio.",
							MarkdownDescription: "Output video aspect ratio.",
						},
						"audio_bitrate": schema.Int64Attribute{
							Computed:            true,
							Description:         "Output audio bitrate, in kbps.",
							MarkdownDescription: "Output audio bitrate, in kbps.",
						},
						"audio_codec": schema.StringAttribute{
							Computed:            true,
							Description:         "Output audio codec.",
							MarkdownDescription: "Output audio codec.",
						},
						"restrict_bitrate": schema.BoolAttribute{
							Computed:            true,
							Description:         "Restrict the output bitrate to the configured video_bitrate/audio_bitrate.",
							MarkdownDescription: "Restrict the output bitrate to the configured `video_bitrate`/`audio_bitrate`.",
						},
						"overlay": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "Overlay (logo/watermark) custom profile, null when not configured.",
							MarkdownDescription: "Overlay (logo/watermark) custom profile, `null` when not configured.",
							Attributes:          computedAttributes(overlaySchemaAttributes()),
						},
						"hls": schema.SingleNestedAttribute{
							Computed:            true,
							Description:         "HLS custom profile, null when not configured.",
							MarkdownDescription: "HLS custom profile, `null` when not configured.",
							Attributes:          computedAttributes(hlsSchemaAttributes()),
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transcodingProfilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TranscodingProfiles

	profiles, err := d.client.GetTranscodingProfiles()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Transcoding Profiles",
			err.Error(),
		)

		return
	}

	// Start from an empty list so that a company without profiles yields [] instead of
	// null, which would break "for p in ...profiles" expressions.
	state.Profiles = []TranscodingProfile{}

	for i := range profiles {
		var item TranscodingProfile

		resp.Diagnostics.Append(applyAPIModel(&item, &profiles[i])...)

		state.Profiles = append(state.Profiles, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *transcodingProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure transcoding profiles client", "error while configuring API client")

		return
	}

	d.client = client
}

// computedAttributes turns the resource's custom profile attributes into computed data source
// attributes, so that both schemas share the same names and descriptions.
func computedAttributes(attrs map[string]rschema.Attribute) map[string]schema.Attribute {
	out := make(map[string]schema.Attribute, len(attrs))

	for name, attr := range attrs {
		switch attr.(type) {
		case rschema.StringAttribute:
			out[name] = schema.StringAttribute{Computed: true, Description: attr.GetDescription(), MarkdownDescription: attr.GetMarkdownDescription()}
		case rschema.Int64Attribute:
			out[name] = schema.Int64Attribute{Computed: true, Description: attr.GetDescription(), MarkdownDescription: attr.GetMarkdownDescription()}
		case rschema.Float64Attribute:
			out[name] = schema.Float64Attribute{Computed: true, Description: attr.GetDescription(), MarkdownDescription: attr.GetMarkdownDescription()}
		default:
			panic(fmt.Sprintf("computedAttributes: unsupported attribute type %T for %q", attr, name))
		}
	}

	return out
}
