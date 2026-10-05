package media

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                     = &transcodingProfileResource{}
	_ resource.ResourceWithConfigure        = &transcodingProfileResource{}
	_ resource.ResourceWithImportState      = &transcodingProfileResource{}
	_ resource.ResourceWithConfigValidators = &transcodingProfileResource{}
	_ resource.ResourceWithModifyPlan       = &transcodingProfileResource{}
)

// NewTranscodingProfileResource is a helper function to simplify the provider implementation.
func NewTranscodingProfileResource() resource.Resource {
	return &transcodingProfileResource{}
}

// resource implementation.
type transcodingProfileResource struct {
	client *teclient.Client
}

// Metadata returns the resource type name.
func (*transcodingProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profile"
}

// Schema defines the schema for the resource.
func (*transcodingProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a transcoding profile.",
		MarkdownDescription: `Manages a transcoding profile, including its ` + "`overlay`" + ` (logo/watermark) and ` + "`hls`" + ` custom profiles.

The API enforces a few rules that this resource surfaces at plan time:

* ` + "`video_format`, `video_codec`, `audio_codec` and `video_aspect`" + ` must be one of the values listed by the ` + "`transparentedge_transcoding_allowed_values`" + ` data source.
* Only one custom profile of each type (` + "`overlay`, `hls`" + `) is allowed per transcoding profile.
* ` + "`overlay`" + ` cannot be combined with ` + "`video_width`" + ` or ` + "`video_height`" + `.

The update operation always replaces the full list of custom profiles: removing ` + "`overlay`" + ` or ` + "`hls`" + ` from the configuration deletes it from the API, and any custom profile added from the dashboard to a profile managed by this resource is removed on the next ` + "`apply`" + `.`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "ID of the transcoding profile.",
				MarkdownDescription: "ID of the transcoding profile.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"company": schema.Int64Attribute{
				Computed:            true,
				Description:         "Company ID that owns this transcoding profile.",
				MarkdownDescription: "Company ID that owns this transcoding profile.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the transcoding profile.",
				MarkdownDescription: "Name of the transcoding profile.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"video_format": schema.StringAttribute{
				Required:            true,
				Description:         "Output video container format, for example mp4 or mpegts. Must be one of the video_formats of the transparentedge_transcoding_allowed_values data source.",
				MarkdownDescription: "Output video container format, for example `mp4` or `mpegts`. Must be one of the `video_formats` of the `transparentedge_transcoding_allowed_values` data source.",
			},
			"video_codec": schema.StringAttribute{
				Required:            true,
				Description:         "Output video codec, for example h264 or libx264. Must be one of the video_codecs of the transparentedge_transcoding_allowed_values data source.",
				MarkdownDescription: "Output video codec, for example `h264` or `libx264`. Must be one of the `video_codecs` of the `transparentedge_transcoding_allowed_values` data source.",
			},
			"video_width": schema.Int64Attribute{
				Optional:            true,
				Description:         "Output video width, in pixels. Cannot be combined with overlay.",
				MarkdownDescription: "Output video width, in pixels. Cannot be combined with `overlay`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"video_height": schema.Int64Attribute{
				Optional:            true,
				Description:         "Output video height, in pixels. Cannot be combined with overlay.",
				MarkdownDescription: "Output video height, in pixels. Cannot be combined with `overlay`.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"video_bitrate": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1000),
				Description:         "Output video bitrate, in kbps.",
				MarkdownDescription: "Output video bitrate, in kbps.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"video_aspect": schema.StringAttribute{
				Optional:            true,
				Description:         "Output video aspect ratio, for example 16:9 or 4:3. Must be one of the video_aspect of the transparentedge_transcoding_allowed_values data source.",
				MarkdownDescription: "Output video aspect ratio, for example `16:9` or `4:3`. Must be one of the `video_aspect` of the `transparentedge_transcoding_allowed_values` data source.",
			},
			"audio_bitrate": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(96),
				Description:         "Output audio bitrate, in kbps.",
				MarkdownDescription: "Output audio bitrate, in kbps.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"audio_codec": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("libfdk_aac"),
				Description:         "Output audio codec, for example libfdk_aac or mp3. Must be one of the audio_codecs of the transparentedge_transcoding_allowed_values data source.",
				MarkdownDescription: "Output audio codec, for example `libfdk_aac` or `mp3`. Must be one of the `audio_codecs` of the `transparentedge_transcoding_allowed_values` data source.",
			},
			"restrict_bitrate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         "Restrict the output bitrate to the configured video_bitrate/audio_bitrate.",
				MarkdownDescription: "Restrict the output bitrate to the configured `video_bitrate`/`audio_bitrate`.",
			},
			"overlay": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Overlay (logo/watermark) custom profile. Cannot be combined with video_width or video_height.",
				MarkdownDescription: "Overlay (logo/watermark) custom profile. Cannot be combined with `video_width` or `video_height`.",
				Attributes:          overlaySchemaAttributes(),
			},
			"hls": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "HLS custom profile.",
				MarkdownDescription: "HLS custom profile.",
				Attributes:          hlsSchemaAttributes(),
			},
		},
	}
}

// overlaySchemaAttributes returns the attributes of the overlay custom profile.
func overlaySchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"url": schema.StringAttribute{
			Required:            true,
			Description:         "URL of the image used as overlay.",
			MarkdownDescription: "URL of the image used as overlay.",
			Validators: []validator.String{
				urlValidator{},
			},
		},
		"position": schema.StringAttribute{
			Required:            true,
			Description:         "Whether the overlay filter is applied before or after the rest of the video filters. One of: before, after.",
			MarkdownDescription: "Whether the overlay filter is applied before or after the rest of the video filters. One of: `before`, `after`.",
			Validators: []validator.String{
				stringvalidator.OneOf("before", "after"),
			},
		},
		"scale_width": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(128),
			Description:         "Width, in pixels, the overlay image is scaled to.",
			MarkdownDescription: "Width, in pixels, the overlay image is scaled to.",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"scale_height": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(128),
			Description:         "Height, in pixels, the overlay image is scaled to.",
			MarkdownDescription: "Height, in pixels, the overlay image is scaled to.",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"opacity": schema.Float64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             float64default.StaticFloat64(0.8),
			Description:         "Opacity of the overlay image, from 0 (transparent) to 1 (opaque).",
			MarkdownDescription: "Opacity of the overlay image, from `0` (transparent) to `1` (opaque).",
			Validators: []validator.Float64{
				float64validator.Between(0, 1),
			},
		},
		"horizontal": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("right"),
			Description:         "Horizontal alignment of the overlay image. One of: left, center, right.",
			MarkdownDescription: "Horizontal alignment of the overlay image. One of: `left`, `center`, `right`.",
			Validators: []validator.String{
				stringvalidator.OneOf("left", "center", "right"),
			},
		},
		"vertical": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("top"),
			Description:         "Vertical alignment of the overlay image. One of: top, center, bottom.",
			MarkdownDescription: "Vertical alignment of the overlay image. One of: `top`, `center`, `bottom`.",
			Validators: []validator.String{
				stringvalidator.OneOf("top", "center", "bottom"),
			},
		},
		"offset_x": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(44),
			Description:         "Horizontal offset, in pixels, applied to the overlay image position.",
			MarkdownDescription: "Horizontal offset, in pixels, applied to the overlay image position.",
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"offset_y": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(44),
			Description:         "Vertical offset, in pixels, applied to the overlay image position.",
			MarkdownDescription: "Vertical offset, in pixels, applied to the overlay image position.",
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
	}
}

// hlsSchemaAttributes returns the attributes of the hls custom profile.
func hlsSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"hls_time": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(5),
			Description:         "Target segment duration, in seconds.",
			MarkdownDescription: "Target segment duration, in seconds.",
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"hls_list_size": schema.Int64Attribute{
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(0),
			Description:         "Maximum number of segments kept in the playlist. 0 keeps all segments.",
			MarkdownDescription: "Maximum number of segments kept in the playlist. `0` keeps all segments.",
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"master_pl_name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("master.m3u8"),
			Description:         "File name of the master playlist.",
			MarkdownDescription: "File name of the master playlist.",
		},
		"hls_flags": schema.StringAttribute{
			Optional:            true,
			Description:         "Value passed to ffmpeg's -hls_flags.",
			MarkdownDescription: "Value passed to ffmpeg's `-hls_flags`.",
		},
		"pix_fmt": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("yuv420p"),
			Description:         "Value passed to ffmpeg's -pix_fmt.",
			MarkdownDescription: "Value passed to ffmpeg's `-pix_fmt`.",
		},
		"framerate": schema.Int64Attribute{
			Optional:            true,
			Description:         "Output framerate.",
			MarkdownDescription: "Output framerate.",
		},
		"h264_preset": schema.StringAttribute{
			Optional:            true,
			Description:         "x264 encoding preset.",
			MarkdownDescription: "x264 encoding preset.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"ultrafast", "superfast", "veryfast", "faster", "fast",
					"medium", "slow", "slower", "veryslow",
				),
			},
		},
		"h264_profile": schema.StringAttribute{
			Optional:            true,
			Description:         "x264 profile. One of: baseline, main, high.",
			MarkdownDescription: "x264 profile. One of: `baseline`, `main`, `high`.",
			Validators: []validator.String{
				stringvalidator.OneOf("baseline", "main", "high"),
			},
		},
		"h264_level": schema.StringAttribute{
			Optional:            true,
			Description:         "x264 level, for example 4.1.",
			MarkdownDescription: "x264 level, for example `4.1`.",
		},
		"maxrate": schema.StringAttribute{
			Optional:            true,
			Description:         "Value passed to ffmpeg's -maxrate, for example 2M.",
			MarkdownDescription: "Value passed to ffmpeg's `-maxrate`, for example `2M`.",
		},
		"bufsize": schema.StringAttribute{
			Optional:            true,
			Description:         "Value passed to ffmpeg's -bufsize, for example 4M.",
			MarkdownDescription: "Value passed to ffmpeg's `-bufsize`, for example `4M`.",
		},
		"b_strategy": schema.Int64Attribute{
			Optional:            true,
			Description:         "Value passed to ffmpeg's -b_strategy.",
			MarkdownDescription: "Value passed to ffmpeg's `-b_strategy`.",
		},
		"refs": schema.Int64Attribute{
			Optional:            true,
			Description:         "Number of reference frames.",
			MarkdownDescription: "Number of reference frames.",
		},
		"coder": schema.Int64Attribute{
			Optional:            true,
			Description:         "Value passed to ffmpeg's -coder: 0 (CAVLC) or 1 (CABAC).",
			MarkdownDescription: "Value passed to ffmpeg's `-coder`: `0` (CAVLC) or `1` (CABAC).",
		},
		"sc_threshold": schema.Int64Attribute{
			Optional:            true,
			Description:         "Scene change detection threshold.",
			MarkdownDescription: "Scene change detection threshold.",
		},
	}
}

// ConfigValidators enforces the API rules that would otherwise fail at apply time with a 400:
// overlay cannot be combined with video_width/video_height ("-filter_complex" vs "-vf"). The "only one custom profile of
// each type" rule is enforced by the schema itself (overlay/hls are single objects).
func (*transcodingProfileResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(path.MatchRoot("overlay"), path.MatchRoot("video_width")),
		resourcevalidator.Conflicting(path.MatchRoot("overlay"), path.MatchRoot("video_height")),
	}
}

// ModifyPlan checks the enum attributes against the API's allowed values, so that an
// unsupported value fails at plan time instead of with a 400 at apply time. Only the
// attributes that change are checked, which skips the request on no-op plans.
func (r *transcodingProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state TranscodingProfile

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if !req.State.Raw.IsNull() {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	changed := map[string]string{}

	for attr, v := range map[string]struct{ planned, prior types.String }{
		"video_format": {plan.VideoFormat, state.VideoFormat},
		"video_codec":  {plan.VideoCodec, state.VideoCodec},
		"audio_codec":  {plan.AudioCodec, state.AudioCodec},
		"video_aspect": {plan.VideoAspect, state.VideoAspect},
	} {
		if !v.planned.IsNull() && !v.planned.IsUnknown() && !v.planned.Equal(v.prior) {
			changed[attr] = v.planned.ValueString()
		}
	}

	if len(changed) == 0 {
		return
	}

	values, err := r.client.GetTranscodingAllowedValues()
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Transcoding Allowed Values", err.Error())

		return
	}

	allowed := map[string][]string{
		"video_format": values.VideoFormats,
		"video_codec":  values.VideoCodecs,
		"audio_codec":  values.AudioCodecs,
		"video_aspect": values.AspectRatio,
	}

	for attr, v := range changed {
		if !slices.Contains(allowed[attr], v) {
			resp.Diagnostics.AddAttributeError(
				path.Root(attr),
				"Invalid value",
				fmt.Sprintf("%q is not allowed by the API, allowed values: %s", v, strings.Join(allowed[attr], ", ")),
			)
		}
	}
}

// Create.
func (r *transcodingProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TranscodingProfile

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateTranscodingProfile(toAPIModel(plan))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating the transcoding profile",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(applyAPIModel(&plan, created)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *transcodingProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TranscodingProfile

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateTranscodingProfile(toAPIModel(plan), int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating the transcoding profile",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(applyAPIModel(&plan, updated)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read resource information.
func (r *transcodingProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TranscodingProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.GetTranscodingProfile(int(state.ID.ValueInt64()))
	if err != nil {
		// Deleted outside of Terraform (from the dashboard, for example): drop it from the
		// state so that the next plan recreates it, instead of failing every plan, apply and
		// destroy until someone runs 'terraform state rm'.
		if errors.Is(err, teclient.ErrTranscodingProfileNotFound) {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"Failure retrieving the transcoding profile from the API",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(applyAPIModel(&state, profile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete.
func (r *transcodingProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TranscodingProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTranscodingProfile(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting the transcoding profile",
			"Could not delete the transcoding profile with id: "+state.ID.String()+"\n"+err.Error(),
		)

		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *transcodingProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*teclient.Client)
	if !ok {
		resp.Diagnostics.AddError("Unable to configure", "error while configuring API client")

		return
	}

	r.client = client
}

func (*transcodingProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
