package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetCertificates() ([]SSLCertificate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslconfig/", c.HostURL, c.CompanyId), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the list of certificates: %s", string(body))
	}

	certs := []SSLCertificate{}
	if err := json.Unmarshal(body, &certs); err != nil {
		return nil, err
	}

	return certs, nil
}
