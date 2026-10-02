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
		TranscodingProfileSummary: TranscodingProfileSummary{
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
			Segmentation:    types.StringNull(),
			RestrictBitrate: types.BoolValue(false),
		},
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
		t.Fatalf("promoted fields lost: %+v", out.TranscodingProfileSummary)
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

	in := TranscodingProfiles{Profiles: []TranscodingProfileSummary{{
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
		Segmentation:    types.StringNull(),
		RestrictBitrate: types.BoolValue(false),
	}}}

	if diags := state.Set(ctx, &in); diags.HasError() {
		t.Fatalf("set: %v", diags)
	}

	var out TranscodingProfiles
	if diags := state.Get(ctx, &out); diags.HasError() {
		t.Fatalf("get: %v", diags)
	}

	if len(out.Profiles) != 1 || out.Profiles[0].Name != types.StringValue("hd") {
		t.Fatalf("summary lost: %+v", out.Profiles)
	}
}

func TestPayloadSendsExplicitNulls(t *testing.T) {
	t.Parallel()

	plan := TranscodingProfile{
		TranscodingProfileSummary: TranscodingProfileSummary{
			Name:            types.StringValue("hd"),
			VideoFormat:     types.StringValue("mp4"),
			VideoCodec:      types.StringValue("h264"),
			VideoWidth:      types.Int64Null(),
			VideoHeight:     types.Int64Null(),
			VideoBitrate:    types.Int64Value(1000),
			VideoAspect:     types.StringNull(),
			AudioBitrate:    types.Int64Value(96),
			AudioCodec:      types.StringNull(),
			Segmentation:    types.StringNull(),
			RestrictBitrate: types.BoolValue(false),
		},
		Overlay: &Overlay{
			URL:         types.StringValue("https://example.com/logo.png"),
			Position:    types.StringValue("after"),
			ScaleWidth:  types.Int64Value(128),
			ScaleHeight: types.Int64Value(128),
			Opacity:     types.Float64Value(0.8),
			Horizontal:  types.StringNull(),
			Vertical:    types.StringNull(),
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
	// survives the update.
	for _, key := range []string{"video_aspect", "audio_codec", "segmentation", "video_width", "video_height"} {
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

	if v, present := overlay["horizontal"]; !present || v != nil {
		t.Fatalf("overlay.horizontal: want explicit null, got %v (present: %t)", v, present)
	}

	// The overlay payload must not carry the hls fields of the response union.
	if _, present := overlay["hls_time"]; present {
		t.Fatalf("overlay payload leaks hls fields: %v", overlay)
	}
}

func TestSegmentationNormalization(t *testing.T) {
	t.Parallel()

	blank := " "
	value := "x"

	if got := normalizeSegmentation(&blank, types.StringNull()); !got.IsNull() {
		t.Fatalf("whitespace with no prior value should be null, got %v", got)
	}

	if got := normalizeSegmentation(&blank, types.StringValue(" ")); got != types.StringValue(" ") {
		t.Fatalf("whitespace asked for by the config should be kept, got %v", got)
	}

	if got := normalizeSegmentation(&blank, types.StringValue("x")); !got.IsNull() {
		t.Fatalf("whitespace clearing a real prior value should be null, got %v", got)
	}

	if got := normalizeSegmentation(&value, types.StringNull()); got != types.StringValue("x") {
		t.Fatalf("real value should be kept, got %v", got)
	}

	if got := normalizeSegmentation(nil, types.StringValue(" ")); !got.IsNull() {
		t.Fatalf("nil should be null, got %v", got)
	}
}
