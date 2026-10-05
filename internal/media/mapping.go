package media

import (
	"fmt"

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
		RestrictBitrate: p.RestrictBitrate.ValueBool(),
		CustomProfiles:  make([]any, 0),
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

// applyAPIModel maps the API response into the model, splitting custom_profiles by
// profile_type into Overlay/HLS. A custom profile of an unknown type is warned about, since
// it is not represented in the schema and the resource's next PUT would delete it.
func applyAPIModel(dst *TranscodingProfile, api *teclient.TranscodingProfileAPIModel) diag.Diagnostics {
	var diags diag.Diagnostics

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
	dst.RestrictBitrate = types.BoolValue(api.RestrictBitrate)

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
					"Transcoding profile %d has a custom profile of type %q that this provider does not support. "+
						"The API deletes custom profile types missing from an update, so the next update of this profile "+
						"through the transparentedge_transcoding_profile resource will delete it.",
					api.ID, cp.ProfileType,
				),
			)
		}
	}

	return diags
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
