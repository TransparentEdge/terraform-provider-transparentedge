data "transparentedge_transcoding_allowed_values" "all" {}

output "allowed_video_formats" {
  value = data.transparentedge_transcoding_allowed_values.all.video_formats
}
