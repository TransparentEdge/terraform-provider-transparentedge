package teclient

import "net/http"

// Environment
type APIEnvironment int

const (
	ProdEnv    APIEnvironment = 0
	StagingEnv APIEnvironment = 1
)

// API
type APIMessage struct {
	Message string `json:"message"`
}

type APIDetail struct {
	Detail string `json:"detail"`
}

type TokenStruct struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
	Scope     string `json:"scope"`
}

type Client struct {
	HTTPClient *http.Client
	Token      TokenStruct

	HostURL      string
	CompanyId    int
	ClientId     string
	ClientSecret string
	VerifySSL    bool
	UserAgent    string
}

// Sites
type SiteAPIModel struct {
	ID      int    `json:"id"`
	Company int    `json:"company"`
	Url     string `json:"url"`
	Active  bool   `json:"active"`
	Ssl     bool   `json:"ssl"`
}

type SiteNewAPIModel struct {
	Url string `json:"url"`
}

type SiteVerifyStringAPIModelRequest struct {
	Domain string `json:"domain"`
}

type SiteVerifyStringAPIModelResponse struct {
	Txt string `json:"txt"`
}

// Backends
type BackendAPIModel struct {
	ID           int    `json:"id"`
	Company      int    `json:"company"`
	Name         string `json:"name"`
	Origin       string `json:"origin"`
	Ssl          bool   `json:"ssl"`
	Port         int    `json:"port"`
	HCHost       string `json:"host"`
	HCPath       string `json:"health_check"`
	HCStatusCode int    `json:"status_code"`
	HCInterval   int    `json:"interval"`
	HCDisabled   bool   `json:"no_probe"`
}

type NewBackendAPIModel struct {
	Name         string `json:"name"`
	Origin       string `json:"origin"`
	Ssl          bool   `json:"ssl"`
	Port         int    `json:"port"`
	HCHost       string `json:"host"`
	HCPath       string `json:"health_check"`
	HCStatusCode int    `json:"status_code"`
	HCInterval   int    `json:"interval"`
	HCDisabled   bool   `json:"no_probe"`
}

// VCL Configs
type VCLConfCreator struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type VCLConfAPIModel struct {
	ID             int            `json:"id"`
	Company        int            `json:"company"`
	VCLCode        string         `json:"config_body"`
	UploadDate     string         `json:"upload_dt"`
	ProductionDate string         `json:"production_dt"`
	Validated      bool           `json:"validated"`
	Active         bool           `json:"active"`
	Deployed       bool           `json:"deployed"`
	CreatorUser    VCLConfCreator `json:"creator_user"`
}

type NewVCLConfAPIModel struct {
	VCLCode string `json:"config_body"`
	Comment string `json:"comment"`
}

// Certificates
type SSLCertificate struct {
	ID            int      `json:"id"`
	Company       int      `json:"company"`
	CommonName    string   `json:"name"`
	Domains       []string `json:"domains"`
	Expiration    string   `json:"expiration"`
	Autogenerated bool     `json:"autogenerated"`
	Standalone    bool     `json:"standalone"`
	DNSChallenge  bool     `json:"dns_challenge"`
	PublicKey     string   `json:"cert"`
	PrivateKey    string   `json:"key"`
}

type SSLCustomCertificate struct {
	ID            int    `json:"id"`
	Autogenerated bool   `json:"autogenerated"`
	DNSChallenge  bool   `json:"dns_challenge"`
	PublicKey     string `json:"cert"`
	PrivateKey    string `json:"key"`
}

// DNS Certificate Requests
type CRDNSKeys struct {
	KeyNameID int    `json:"key_name_id"`
	KeyName   string `json:"key_name"`
}

type CRDNSProvider struct {
	ID       int         `json:"id"`
	HookName string      `json:"hook_name"`
	Provider string      `json:"provider"`
	Keys     []CRDNSKeys `json:"keys"`
}

type CRDNSCNAMEVerification struct {
	CNAME string `json:"cname"`
}

type CRDNSCreds struct {
	Provider string `json:"provider"`
	KeyName  string `json:"dns_key_name"`
	KeyValue string `json:"dns_key_value"`
}

type CRDNSCredential struct {
	ID    int          `json:"id"`
	Alias string       `json:"alias"`
	Creds []CRDNSCreds `json:"creds"`
}

type NewCRDNSCreds struct {
	KeyName  string `json:"hook_key_name"`
	KeyValue string `json:"hook_key_value"`
}

type NewCRDNSCredential struct {
	Alias string          `json:"alias"`
	Creds []NewCRDNSCreds `json:"creds"`
}

type CertReqDNS struct {
	ID         int    `json:"id"`
	Credential int    `json:"credential"`
	Domains    string `json:"domains"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	// as pointers since they can be null
	CertificateID *int    `json:"certificate"`
	Log           *string `json:"log"`
}

// HTTP Certificate Requests
type CertReqHTTP struct {
	ID         int    `json:"id"`
	CommonName string `json:"cn"`
	SAN        string `json:"san"`
	Standalone bool   `json:"standalone"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"last_modified_at"`
	// as pointers since they can be null
	CertificateID *int    `json:"certificate"`
	Log           *string `json:"error_text"`
}
