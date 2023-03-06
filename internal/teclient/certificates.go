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
		return nil, fmt.Errorf("Couldn't retrieve the list of certificates: %s", c.parseAPIError(body))
	}

	certs := []SSLCertificate{}
	if err := json.Unmarshal(body, &certs); err != nil {
		return nil, err
	}

	return certs, nil
}

func (c *Client) GetCertificate(certID int) (*SSLCertificate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslconfig/%d/", c.HostURL, c.CompanyId, certID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the Custom Certificate with ID: %d. %s", certID, c.parseAPIError(body))
	}

	cert := SSLCertificate{}
	if err := json.Unmarshal(body, &cert); err != nil {
		return nil, err
	}

	return &cert, nil
}

func (c *Client) CreateCustomCertificate(cert SSLCustomCertificate) (*SSLCertificate, error) {
	req, err := c.prepareJSONRequest(cert, "POST", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslconfig/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newCustomCertificate := SSLCertificate{}
	if err := json.Unmarshal(body, &newCustomCertificate); err != nil {
		return nil, err
	}

	return &newCustomCertificate, nil
}

func (c *Client) UpdateCustomCertificate(cert SSLCustomCertificate) (*SSLCertificate, error) {
	req, err := c.prepareJSONRequest(cert, "PUT", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslconfig/%d/", c.HostURL, c.CompanyId, cert.ID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	// API doesnt return the new model at this moment, try to get the updated certificate
	newCustomCertificate, err := c.GetCertificate(cert.ID)
	if err != nil {
		return nil, fmt.Errorf("Certificate was updated but couldn't retrieve the new data from API, an import is required. " + err.Error())
	}

	return newCustomCertificate, nil
}

func (c *Client) DeleteCustomCertificate(certID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslconfig/%d/", c.HostURL, c.CompanyId, certID), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if sc != 204 {
		return fmt.Errorf("%d - API request failed trying to DELETE the Certificate ID %d: %s", sc, certID, c.parseAPIError(body))
	}

	return nil
}
