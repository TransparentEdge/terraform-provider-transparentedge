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
