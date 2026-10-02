// Package media implements the /v1/media/ API domain (transcoding profiles).
package media

import "github.com/hashicorp/terraform-plugin-framework/types"

// TranscodingProfileSummary holds the fields of the TranscodingProfile serializer, the one
// the list endpoint returns. It is the element of the transcoding_profiles data source and
// is embedded in TranscodingProfile (the framework promotes the tfsdk tags of an embedded
// struct), so that both share a single mapping.
type TranscodingProfileSummary struct {
	ID              types.Int64  `tfsdk:"id"`
	Company         types.Int64  `tfsdk:"company"`
	Name            types.String `tfsdk:"name"`
	VideoFormat     types.String `tfsdk:"video_format"`
	VideoCodec      types.String `tfsdk:"video_codec"`
	VideoWidth      types.Int64  `tfsdk:"video_width"`
	VideoHeight     types.Int64  `tfsdk:"video_height"`
	VideoBitrate    types.Int64  `tfsdk:"video_bitrate"`
	VideoAspect     types.String `tfsdk:"video_aspect"`
	AudioBitrate    types.Int64  `tfsdk:"audio_bitrate"`
	AudioCodec      types.String `tfsdk:"audio_codec"`
	Segmentation    types.String `tfsdk:"segmentation"`
	RestrictBitrate types.Bool   `tfsdk:"restrict_bitrate"`
}

// TranscodingProfile is the transcoding_profile resource model: the flat fields of the
// listing plus the custom profiles, which only the detail serializer returns.
type TranscodingProfile struct {
	TranscodingProfileSummary

	Overlay *Overlay `tfsdk:"overlay"`
	HLS     *HLS     `tfsdk:"hls"`
}

// Overlay is the "overlay" custom profile (logo/watermark).
type Overlay struct {
	URL         types.String  `tfsdk:"url"`
	Position    types.String  `tfsdk:"position"`
	ScaleWidth  types.Int64   `tfsdk:"scale_width"`
	ScaleHeight types.Int64   `tfsdk:"scale_height"`
	Opacity     types.Float64 `tfsdk:"opacity"`
	Horizontal  types.String  `tfsdk:"horizontal"`
	Vertical    types.String  `tfsdk:"vertical"`
	OffsetX     types.Int64   `tfsdk:"offset_x"`
	OffsetY     types.Int64   `tfsdk:"offset_y"`
}

// HLS is the "hls" custom profile.
type HLS struct {
	Position     types.String `tfsdk:"position"`
	HLSTime      types.Int64  `tfsdk:"hls_time"`
	HLSListSize  types.Int64  `tfsdk:"hls_list_size"`
	MasterPlName types.String `tfsdk:"master_pl_name"`
	HLSFlags     types.String `tfsdk:"hls_flags"`
	PixFmt       types.String `tfsdk:"pix_fmt"`
	Framerate    types.String `tfsdk:"framerate"`
	H264Preset   types.String `tfsdk:"h264_preset"`
	H264Profile  types.String `tfsdk:"h264_profile"`
	H264Level    types.String `tfsdk:"h264_level"`
	Maxrate      types.String `tfsdk:"maxrate"`
	Bufsize      types.String `tfsdk:"bufsize"`
	BStrategy    types.Int64  `tfsdk:"b_strategy"`
	Refs         types.Int64  `tfsdk:"refs"`
	Coder        types.String `tfsdk:"coder"`
	ScThreshold  types.Int64  `tfsdk:"sc_threshold"`
}

// TranscodingProfiles is the transcoding_profiles data source model.
type TranscodingProfiles struct {
	Profiles []TranscodingProfileSummary `tfsdk:"profiles"`
}

// TranscodingAllowedValues is the transcoding_allowed_values data source model.
type TranscodingAllowedValues struct {
	VideoFormats types.List `tfsdk:"video_formats"`
	VideoCodecs  types.List `tfsdk:"video_codecs"`
	AudioCodecs  types.List `tfsdk:"audio_codecs"`
	VideoAspect  types.List `tfsdk:"video_aspect"`
}
