# Import the active VCL configuration, the ID value doesn't matter.
terraform import 'transparentedge_vclconf.prod' 0

# Importing VCL code has its quirks, since the API parses the code and
# add/removes newlines and spaces the diff won't be equal, it's usually
# better to just copy the last configuration from the dashboard and
# apply again.
