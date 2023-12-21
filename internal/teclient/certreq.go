package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetCRDNSProviders() ([]CRDNSProvider, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/dnshook/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Failure retrieving DNS Providers: %s", c.parseAPIError(body))
	}

	providers := []CRDNSProvider{}
	if err := json.Unmarshal(body, &providers); err != nil {
		return nil, err
	}

	return providers, nil
}

func (c *Client) GetDNSCNAMEVerification() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/autoprovisioning/%d/ssldnsverificationcname/", c.HostURL, c.CompanyId), nil)
	if err != nil {
		return "", err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return "", err
	}
	if sc != 200 {
		return "", fmt.Errorf("Failure retrieving CNAME verification: %s", c.parseAPIError(body))
	}

	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	cnameValue, exists := data["cname"].(string)
	if !exists {
		return "", fmt.Errorf("[API ERROR] Cname not found in the payload.")
	}

	return cnameValue, nil
}
