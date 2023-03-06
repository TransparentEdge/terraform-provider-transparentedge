package staging

import "github.com/hashicorp/terraform-plugin-framework/types"

type StagingBackend struct {
	ID           types.Int64  `tfsdk:"id"`
	Company      types.Int64  `tfsdk:"company"`
	Name         types.String `tfsdk:"name"`
	VclName      types.String `tfsdk:"vclname"`
	Origin       types.String `tfsdk:"origin"`
	Ssl          types.Bool   `tfsdk:"ssl"`
	Port         types.Int64  `tfsdk:"port"`
	HCHost       types.String `tfsdk:"hchost"`
	HCPath       types.String `tfsdk:"hcpath"`
	HCStatusCode types.Int64  `tfsdk:"hcstatuscode"`
}

type StagingVCLConf struct {
	ID             types.Int64  `tfsdk:"id"`
	Company        types.Int64  `tfsdk:"company"`
	VCLCode        types.String `tfsdk:"vclcode"`
	UploadDate     types.String `tfsdk:"uploaddate"`
	ProductionDate types.String `tfsdk:"productiondate"`
	User           types.String `tfsdk:"user"`
}
