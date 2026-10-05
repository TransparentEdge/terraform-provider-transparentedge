package media

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceModelRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res := &transcodingProfileResource{}
	schemaResp := &resource.SchemaResponse{}
	res.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	state := tfsdk.State{Schema: schemaResp.Schema}

	in := TranscodingProfile{
		ID:              types.Int64Value(7),
		Company:         types.Int64Value(1),
		Name:            types.StringValue("hd"),
		VideoFormat:     types.StringValue("mp4"),
		VideoCodec:      types.StringValue("h264"),
		VideoWidth:      types.Int64Null(),
		VideoHeight:     types.Int64Null(),
		VideoBitrate:    types.Int64Value(1000),
		VideoAspect:     types.StringValue("16:9"),
		AudioBitrate:    types.Int64Value(96),
		AudioCodec:      types.StringNull(),
		RestrictBitrate: types.BoolValue(false),
		Overlay: &Overlay{
			URL:         types.StringValue("https://example.com/logo.png"),
			Position:    types.StringValue("after"),
			ScaleWidth:  types.Int64Value(128),
			ScaleHeight: types.Int64Value(128),
			Opacity:     types.Float64Value(0.8),
			Horizontal:  types.StringValue("right"),
			Vertical:    types.StringValue("top"),
			OffsetX:     types.Int64Value(44),
			OffsetY:     types.Int64Value(44),
		},
	}

	if diags := state.Set(ctx, &in); diags.HasError() {
		t.Fatalf("set: %v", diags)
	}

	var out TranscodingProfile
	if diags := state.Get(ctx, &out); diags.HasError() {
		t.Fatalf("get: %v", diags)
	}

	if out.ID != in.ID || out.Name != in.Name || out.VideoAspect != in.VideoAspect {
		t.Fatalf("fields lost: %+v", out)
	}

	if out.Overlay == nil || out.Overlay.OffsetX != types.Int64Value(44) {
		t.Fatalf("overlay lost: %+v", out.Overlay)
	}
}

func TestDataSourceModelRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ds := &transcodingProfilesDataSource{}
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, schemaResp)

	state := tfsdk.State{Schema: schemaResp.Schema}

	in := TranscodingProfiles{Profiles: []TranscodingProfile{{
		ID:              types.Int64Value(7),
		Company:         types.Int64Value(1),
		Name:            types.StringValue("hd"),
		VideoFormat:     types.StringValue("mp4"),
		VideoCodec:      types.StringValue("h264"),
		VideoWidth:      types.Int64Null(),
		VideoHeight:     types.Int64Null(),
		VideoBitrate:    types.Int64Value(1000),
		VideoAspect:     types.StringNull(),
		AudioBitrate:    types.Int64Value(96),
		AudioCodec:      types.StringNull(),
		RestrictBitrate: types.BoolValue(false),
		HLS: &HLS{
			HLSTime:      types.Int64Value(5),
			HLSListSize:  types.Int64Value(0),
			MasterPlName: types.StringValue("master.m3u8"),
			HLSFlags:     types.StringNull(),
			PixFmt:       types.StringValue("yuv420p"),
			Framerate:    types.Int64Null(),
			H264Preset:   types.StringNull(),
			H264Profile:  types.StringNull(),
			H264Level:    types.StringNull(),
			Maxrate:      types.StringNull(),
			Bufsize:      types.StringNull(),
			BStrategy:    types.Int64Null(),
			Refs:         types.Int64Null(),
			Coder:        types.Int64Null(),
			ScThreshold:  types.Int64Null(),
		},
	}}}

	if diags := state.Set(ctx, &in); diags.HasError() {
		t.Fatalf("set: %v", diags)
	}

	var out TranscodingProfiles
	if diags := state.Get(ctx, &out); diags.HasError() {
		t.Fatalf("get: %v", diags)
	}

	if len(out.Profiles) != 1 || out.Profiles[0].Name != types.StringValue("hd") {
		t.Fatalf("profile lost: %+v", out.Profiles)
	}

	if out.Profiles[0].Overlay != nil || out.Profiles[0].HLS == nil || out.Profiles[0].HLS.HLSTime != types.Int64Value(5) {
		t.Fatalf("custom profiles lost: %+v", out.Profiles[0])
	}
}

func TestPayloadSendsExplicitNulls(t *testing.T) {
	t.Parallel()

	plan := TranscodingProfile{
		Name:            types.StringValue("hd"),
		VideoFormat:     types.StringValue("mp4"),
		VideoCodec:      types.StringValue("h264"),
		VideoWidth:      types.Int64Null(),
		VideoHeight:     types.Int64Null(),
		VideoBitrate:    types.Int64Value(1000),
		VideoAspect:     types.StringNull(),
		AudioBitrate:    types.Int64Value(96),
		AudioCodec:      types.StringValue("libfdk_aac"),
		RestrictBitrate: types.BoolValue(false),
		Overlay: &Overlay{
			URL:         types.StringValue("https://example.com/logo.png"),
			Position:    types.StringValue("after"),
			ScaleWidth:  types.Int64Value(128),
			ScaleHeight: types.Int64Value(128),
			Opacity:     types.Float64Value(0.8),
			Horizontal:  types.StringValue("right"),
			Vertical:    types.StringValue("top"),
			OffsetX:     types.Int64Value(44),
			OffsetY:     types.Int64Value(44),
		},
	}

	body, err := json.Marshal(toAPIModel(plan))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	payload := map[string]any{}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// A cleared field has to reach the API as an explicit null, or its previous value
	// survives the update. audio_codec is not in the list: the API rejects a null there, so
	// the schema gives it a default instead.
	for _, key := range []string{"video_aspect", "video_width", "video_height"} {
		if v, ok := payload[key]; !ok || v != nil {
			t.Fatalf("%s: want explicit null, got %v (present: %t)", key, v, ok)
		}
	}

	profiles, ok := payload["custom_profiles"].([]any)
	if !ok || len(profiles) != 1 {
		t.Fatalf("custom_profiles: want one element, got %v", payload["custom_profiles"])
	}

	overlay, ok := profiles[0].(map[string]any)
	if !ok {
		t.Fatalf("custom_profiles[0]: want an object, got %v", profiles[0])
	}

	if v := overlay["horizontal"]; v != "right" {
		t.Fatalf("overlay.horizontal: want the planned default, got %v", v)
	}

	// The overlay payload must not carry the hls fields of the response union.
	if _, present := overlay["hls_time"]; present {
		t.Fatalf("overlay payload leaks hls fields: %v", overlay)
	}
}

func TestHLSPayload(t *testing.T) {
	t.Parallel()

	plan := TranscodingProfile{
		Name:        types.StringValue("hls"),
		VideoFormat: types.StringValue("mpegts"),
		VideoCodec:  types.StringValue("h264"),
		AudioCodec:  types.StringValue("libfdk_aac"),
		HLS: &HLS{
			HLSTime:      types.Int64Value(5),
			HLSListSize:  types.Int64Value(0),
			MasterPlName: types.StringValue("master.m3u8"),
			HLSFlags:     types.StringNull(),
			PixFmt:       types.StringValue("yuv420p"),
			Framerate:    types.Int64Value(25),
			H264Preset:   types.StringNull(),
			H264Profile:  types.StringNull(),
			H264Level:    types.StringNull(),
			Maxrate:      types.StringNull(),
			Bufsize:      types.StringNull(),
			BStrategy:    types.Int64Null(),
			Refs:         types.Int64Null(),
			Coder:        types.Int64Value(1),
			ScThreshold:  types.Int64Null(),
		},
	}

	body, err := json.Marshal(toAPIModel(plan))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	payload := struct {
		CustomProfiles []map[string]any `json:"custom_profiles"`
	}{}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(payload.CustomProfiles) != 1 {
		t.Fatalf("custom_profiles: want one element, got %v", payload.CustomProfiles)
	}

	hls := payload.CustomProfiles[0]

	// The API does not store a position for hls, so it is not sent.
	if _, present := hls["position"]; present {
		t.Fatalf("hls payload carries position: %v", hls)
	}

	// framerate and coder are integers in the API.
	if hls["framerate"] != float64(25) || hls["coder"] != float64(1) {
		t.Fatalf("framerate/coder: want numbers, got %v / %v", hls["framerate"], hls["coder"])
	}

	if v, present := hls["hls_flags"]; !present || v != nil {
		t.Fatalf("hls_flags: want explicit null, got %v (present: %t)", v, present)
	}
}
