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
