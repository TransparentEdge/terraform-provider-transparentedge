resource "transparentedge_transcoding_profile" "hd_with_logo" {
  name          = "hd-with-logo"
  video_format  = "mp4"
  video_codec   = "h264"
  video_bitrate = 2500
  audio_codec   = "libfdk_aac"
  audio_bitrate = 128

  # Logo/watermark, positioned on the top-right corner
  overlay = {
    url        = "https://static.example.com/logo.png"
    position   = "after"
    horizontal = "right"
    vertical   = "top"
    offset_x   = 44
    offset_y   = 44
    opacity    = 0.8
  }
}
