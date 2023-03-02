# Sites can be imported specifying the domain (single site)
terraform import 'transparentedge_site.www_example3_com' 'www.example3.com'

# If the sites are in a list/set, import a single site from the list with:
terraform import 'transparentedge_site.all["www.example2.com"]' 'www.example2.com'
