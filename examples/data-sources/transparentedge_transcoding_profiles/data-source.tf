data "transparentedge_transcoding_profiles" "all" {}

output "transcoding_profile_ids" {
  value = { for p in data.transparentedge_transcoding_profiles.all.profiles : p.name => p.id }
}
