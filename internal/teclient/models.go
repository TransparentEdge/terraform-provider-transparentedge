package teclient

type ErrorAPIMessage struct {
	Message string `json:"message"`
}

type SiteAPIModel struct {
	ID      int    `json:"id"`
	Url     string `json:"url"`
	Company int    `json:"company"`
	Ssl     bool   `json:"ssl"`
	Active  bool   `json:"active"`
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
}

type NewBackendAPIModel struct {
	Name         string `json:"name"`
	Origin       string `json:"origin"`
	Ssl          bool   `json:"ssl"`
	Port         int    `json:"port"`
	HCHost       string `json:"host"`
	HCPath       string `json:"health_check"`
	HCStatusCode int    `json:"status_code"`
}

type VCLConfCreator struct {
	ID        int `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type VCLConfAPIModel struct {
	ID             int         `json:"id"`
	Company        int            `json:"company"`
	VCLCode        string         `json:"config_body"`
	UploadDate     string         `json:"upload_dt"`
	ProductionDate string         `json:"production_dt"`
	Validated      bool           `json:"validated"`
	Active         bool           `json:"active"`
	Deployed       bool           `json:"deployed"`
	CreatorUser    VCLConfCreator `json:"creator_user"`
}
