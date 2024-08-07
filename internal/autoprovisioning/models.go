package autoprovisioning

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Site struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	ID       types.Int64    `tfsdk:"id"`
	Domain   types.String   `tfsdk:"domain"`
	Active   types.Bool     `tfsdk:"active"`
}

type SiteDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Company types.Int64  `tfsdk:"company"`
	Domain  types.String `tfsdk:"domain"`
	Active  types.Bool   `tfsdk:"active"`
	Ssl     types.Bool   `tfsdk:"ssl"`
}

type Sites struct {
	Sites []SiteDataSourceModel `tfsdk:"sites"`
}

type SiteVerify struct {
	Domain              types.String `tfsdk:"domain"`
	VerificantionString types.String `tfsdk:"verification_string"`
}

type Backend struct {
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
	HCInterval   types.Int64  `tfsdk:"hcinterval"`
	HCDisabled   types.Bool   `tfsdk:"hcdisabled"`
}

type Backends struct {
	Backends []Backend `tfsdk:"backends"`
}

type VCLConf struct {
	ID             types.Int64  `tfsdk:"id"`
	Company        types.Int64  `tfsdk:"company"`
	VCLCode        types.String `tfsdk:"vclcode"`
	UploadDate     types.String `tfsdk:"uploaddate"`
	ProductionDate types.String `tfsdk:"productiondate"`
	User           types.String `tfsdk:"user"`
}

type Certificates struct {
	Certificates []Certificate `tfsdk:"certificates"`
}

type Certificate struct {
	ID            types.Int64  `tfsdk:"id"`
	Company       types.Int64  `tfsdk:"company"`
	CommonName    types.String `tfsdk:"commonname"`
	Domains       types.String `tfsdk:"domains"`
	Expiration    types.String `tfsdk:"expiration"`
	Autogenerated types.Bool   `tfsdk:"autogenerated"`
	Standalone    types.Bool   `tfsdk:"standalone"`
	DNSChallenge  types.Bool   `tfsdk:"dnschallenge"`
	PublicKey     types.String `tfsdk:"publickey"`
	PrivateKey    types.String `tfsdk:"privatekey"`
}

type CustomCertificate struct {
	ID         types.Int64  `tfsdk:"id"`
	CommonName types.String `tfsdk:"commonname"`
	Domains    types.String `tfsdk:"domains"`
	Expiration types.String `tfsdk:"expiration"`
	PublicKey  types.String `tfsdk:"publickey"`
	PrivateKey types.String `tfsdk:"privatekey"`
}

// DNS Certificate Requests
type CertReqDNSProvider struct {
	DNSProvider types.String `tfsdk:"dns_provider"`
	Parameters  types.List   `tfsdk:"parameters"`
}

type CertReqDNSProviders struct {
	Providers []CertReqDNSProvider `tfsdk:"providers"`
}

type CertReqCNAMEVerify struct {
	CNAME types.String `tfsdk:"cname"`
}

type CertReqDNSCredential struct {
	ID          types.Int64  `tfsdk:"id"`
	Alias       types.String `tfsdk:"alias"`
	DNSProvider types.String `tfsdk:"dns_provider"`
	Parameters  types.Map    `tfsdk:"parameters"`
}

type CertReqDNS struct {
	ID            types.Int64  `tfsdk:"id"`
	Domains       types.Set    `tfsdk:"domains"`
	Credential    types.Int64  `tfsdk:"credential"`
	CertificateID types.Int64  `tfsdk:"certificate_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	StatusMessage types.String `tfsdk:"status_message"`
}

// HTTP Certificate requests
type CertReqHTTP struct {
	ID            types.Int64  `tfsdk:"id"`
	Domains       types.Set    `tfsdk:"domains"`
	Standalone    types.Bool   `tfsdk:"standalone"`
	CertificateID types.Int64  `tfsdk:"certificate_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	StatusMessage types.String `tfsdk:"status_message"`
}
