package media

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
)

// toAPIModel builds the create/update payload from the plan. custom_profiles is always a
// non-nil slice (even empty) so that removing overlay/hls from the config actually
// deletes it: the API replaces the full list on every PUT.
func toAPIModel(p TranscodingProfile) teclient.NewTranscodingProfileAPIModel {
	api := teclient.NewTranscodingProfileAPIModel{
		Name:            p.Name.ValueString(),
		VideoFormat:     p.VideoFormat.ValueString(),
		VideoCodec:      p.VideoCodec.ValueString(),
		VideoWidth:      int64ToIntPtr(p.VideoWidth),
		VideoHeight:     int64ToIntPtr(p.VideoHeight),
		VideoBitrate:    int(p.VideoBitrate.ValueInt64()),
		VideoAspect:     stringToPtr(p.VideoAspect),
		AudioBitrate:    int(p.AudioBitrate.ValueInt64()),
		AudioCodec:      stringToPtr(p.AudioCodec),
		Segmentation:    stringToPtr(p.Segmentation),
		RestrictBitrate: p.RestrictBitrate.ValueBool(),
		CustomProfiles:  make([]teclient.CustomProfileRequest, 0),
	}

	if p.Overlay != nil {
		api.CustomProfiles = append(api.CustomProfiles, teclient.OverlayProfileAPIModel{
			ProfileType: "overlay",
			Position:    p.Overlay.Position.ValueString(),
			URL:         stringToPtr(p.Overlay.URL),
			ScaleWidth:  int64ToIntPtr(p.Overlay.ScaleWidth),
			ScaleHeight: int64ToIntPtr(p.Overlay.ScaleHeight),
			Opacity:     float64ToPtr(p.Overlay.Opacity),
			Horizontal:  stringToPtr(p.Overlay.Horizontal),
			Vertical:    stringToPtr(p.Overlay.Vertical),
			OffsetX:     int64ToIntPtr(p.Overlay.OffsetX),
			OffsetY:     int64ToIntPtr(p.Overlay.OffsetY),
		})
	}

	if p.HLS != nil {
		api.CustomProfiles = append(api.CustomProfiles, teclient.HLSProfileAPIModel{
			ProfileType:  "hls",
			HLSTime:      int64ToIntPtr(p.HLS.HLSTime),
			HLSListSize:  int64ToIntPtr(p.HLS.HLSListSize),
			MasterPlName: stringToPtr(p.HLS.MasterPlName),
			HLSFlags:     stringToPtr(p.HLS.HLSFlags),
			PixFmt:       stringToPtr(p.HLS.PixFmt),
			Framerate:    int64ToIntPtr(p.HLS.Framerate),
			H264Preset:   stringToPtr(p.HLS.H264Preset),
			H264Profile:  stringToPtr(p.HLS.H264Profile),
			H264Level:    stringToPtr(p.HLS.H264Level),
			Maxrate:      stringToPtr(p.HLS.Maxrate),
			Bufsize:      stringToPtr(p.HLS.Bufsize),
			BStrategy:    int64ToIntPtr(p.HLS.BStrategy),
			Refs:         int64ToIntPtr(p.HLS.Refs),
			Coder:        int64ToIntPtr(p.HLS.Coder),
			ScThreshold:  int64ToIntPtr(p.HLS.ScThreshold),
		})
	}

	return api
}

// applyFlatAPIModel maps the fields of the TranscodingProfile serializer, the ones both the
// list and the detail endpoints return.
func applyFlatAPIModel(dst *TranscodingProfileSummary, api *teclient.TranscodingProfileAPIModel) {
	dst.ID = types.Int64Value(int64(api.ID))
	dst.Company = types.Int64Value(int64(api.Company))
	dst.Name = types.StringValue(api.Name)
	dst.VideoFormat = types.StringValue(api.VideoFormat)
	dst.VideoCodec = types.StringValue(api.VideoCodec)
	dst.VideoWidth = intPtrToInt64(api.VideoWidth)
	dst.VideoHeight = intPtrToInt64(api.VideoHeight)
	dst.VideoBitrate = types.Int64Value(int64(api.VideoBitrate))
	dst.VideoAspect = strPtrToString(api.VideoAspect)
	dst.AudioBitrate = types.Int64Value(int64(api.AudioBitrate))
	dst.AudioCodec = strPtrToString(api.AudioCodec)
	dst.Segmentation = normalizeSegmentation(api.Segmentation, dst.Segmentation)
	dst.RestrictBitrate = types.BoolValue(api.RestrictBitrate)
}

// applyAPIModel maps the detail API response into the resource model, splitting
// custom_profiles by profile_type into Overlay/HLS. A custom profile of an unknown type is
// warned about, since the next PUT would silently delete it (it is not represented in the
// schema). The warning is worded for the resource on purpose: this is the only caller that
// writes custom_profiles back.
func applyAPIModel(dst *TranscodingProfile, api *teclient.TranscodingProfileAPIModel) diag.Diagnostics {
	var diags diag.Diagnostics

	applyFlatAPIModel(&dst.TranscodingProfileSummary, api)

	dst.Overlay = nil
	dst.HLS = nil

	for _, cp := range api.CustomProfiles {
		switch cp.ProfileType {
		case "overlay":
			dst.Overlay = &Overlay{
				URL:         strPtrToString(cp.URL),
				Position:    types.StringValue(cp.Position),
				ScaleWidth:  intPtrToInt64(cp.ScaleWidth),
				ScaleHeight: intPtrToInt64(cp.ScaleHeight),
				Opacity:     floatPtrToFloat64(cp.Opacity),
				Horizontal:  strPtrToString(cp.Horizontal),
				Vertical:    strPtrToString(cp.Vertical),
				OffsetX:     intPtrToInt64(cp.OffsetX),
				OffsetY:     intPtrToInt64(cp.OffsetY),
			}
		case "hls":
			dst.HLS = &HLS{
				HLSTime:      intPtrToInt64(cp.HLSTime),
				HLSListSize:  intPtrToInt64(cp.HLSListSize),
				MasterPlName: strPtrToString(cp.MasterPlName),
				HLSFlags:     strPtrToString(cp.HLSFlags),
				PixFmt:       strPtrToString(cp.PixFmt),
				Framerate:    intPtrToInt64(cp.Framerate),
				H264Preset:   strPtrToString(cp.H264Preset),
				H264Profile:  strPtrToString(cp.H264Profile),
				H264Level:    strPtrToString(cp.H264Level),
				Maxrate:      strPtrToString(cp.Maxrate),
				Bufsize:      strPtrToString(cp.Bufsize),
				BStrategy:    intPtrToInt64(cp.BStrategy),
				Refs:         intPtrToInt64(cp.Refs),
				Coder:        intPtrToInt64(cp.Coder),
				ScThreshold:  intPtrToInt64(cp.ScThreshold),
			}
		default:
			diags.AddWarning(
				"Unmanaged custom profile found",
				fmt.Sprintf(
					"Transcoding profile %d has a custom profile of type %q that this provider does not know how to represent. "+
						"The next apply will delete it, since the API replaces the full custom_profiles list on every update.",
					api.ID, cp.ProfileType,
				),
			)
		}
	}

	return diags
}

// normalizeSegmentation avoids a perpetual diff: the API can return " " for an unset
// segmentation instead of null. A whitespace value that the configuration asked for is kept
// as it is (the API enum accepts " "), otherwise writing it would make the state differ from
// the plan; prior is the value already in the plan or state, if any.
func normalizeSegmentation(s *string, prior types.String) types.String {
	if s == nil {
		return types.StringNull()
	}

	if strings.TrimSpace(*s) == "" && !isWhitespaceString(prior) {
		return types.StringNull()
	}

	return types.StringValue(*s)
}

// isWhitespaceString reports whether v holds a known, non-null, whitespace-only string.
func isWhitespaceString(v types.String) bool {
	if v.IsNull() || v.IsUnknown() {
		return false
	}

	return strings.TrimSpace(v.ValueString()) == ""
}

func int64ToIntPtr(v types.Int64) *int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	i := int(v.ValueInt64())

	return &i
}

func stringToPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	s := v.ValueString()

	return &s
}

func float64ToPtr(v types.Float64) *float64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	f := v.ValueFloat64()

	return &f
}

func intPtrToInt64(v *int) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*v))
}

func strPtrToString(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}

	return types.StringValue(*v)
}

func floatPtrToFloat64(v *float64) types.Float64 {
	if v == nil {
		return types.Float64Null()
	}

	return types.Float64Value(*v)
}
