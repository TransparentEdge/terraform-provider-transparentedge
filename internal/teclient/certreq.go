package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) GetCRDNSProviders() ([]CRDNSProvider, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/autoprovisioning/dnshook/", c.HostURL), nil) // nolint: perfsprint
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure retrieving DNS Providers: %s", c.parseAPIError(body))
	}

	providers := []CRDNSProvider{}

	err = json.Unmarshal(body, &providers)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func (c *Client) GetDNSCNAMEVerification() (string, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/autoprovisioning/%d/ssldnsverificationcname/", c.HostURL, c.CompanyID), nil)
	if err != nil {
		return "", err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	if sc != http.StatusOK {
		return "", fmt.Errorf("failure retrieving CNAME verification: %s", c.parseAPIError(body))
	}

	var data map[string]any

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	cnameValue, exists := data["cname"].(string)
	if !exists {
		return "", errors.New("fatal: CNAME not found in the API response")
	}

	return cnameValue, nil
}

func (c *Client) GetCRDNSCredential(id int) (CRDNSCredential, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d/", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return CRDNSCredential{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CRDNSCredential{}, err
	}

	if sc != http.StatusOK {
		return CRDNSCredential{}, fmt.Errorf("failure retrieving DNS Credential with id %d. %s", id, c.parseAPIError(body))
	}

	data := CRDNSCredential{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return CRDNSCredential{}, err
	}

	return data, nil
}

func (c *Client) CreateDNSCredential(dnsCredential NewCRDNSCredential) (*CRDNSCredential, error) {
	req, err := c.prepareJSONRequest(dnsCredential, http.MethodPost, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/", c.HostURL, c.CompanyID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc == http.StatusBadRequest {
		apiError := c.parseAPIError(body)

		err := json.Unmarshal(body, &apiError)
		if err == nil {
			return nil, errors.New(apiError)
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newData := CRDNSCredential{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}

func (c *Client) UpdateDNSCredential(dnsCredential NewCRDNSCredential, id int) (*CRDNSCredential, error) {
	req, err := c.prepareJSONRequest(dnsCredential, http.MethodPut, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d", c.HostURL, c.CompanyID, id))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc == http.StatusBadRequest {
		apiError := c.parseAPIError(body)

		err := json.Unmarshal(body, &apiError)
		if err == nil {
			return nil, errors.New(apiError)
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newData := CRDNSCredential{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}

func (c *Client) DeleteCRCredential(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if sc != 204 {
		return fmt.Errorf("%d - DELETE Failed for ID %d: %s", sc, id, c.parseAPIError(body))
	}

	return nil
}

func (c *Client) GetCertReqDNS(id int) (CertReqDNS, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return CertReqDNS{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CertReqDNS{}, err
	}

	if sc != http.StatusOK {
		return CertReqDNS{}, fmt.Errorf("failure retrieving DNS Certificate Request with id %d. %s", id, c.parseAPIError(body))
	}

	data := CertReqDNS{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return CertReqDNS{}, err
	}

	return data, nil
}

func (c *Client) CreateDNSCertReq(certreq any) (*CertReqDNS, error) {
	req, err := c.prepareJSONRequest(certreq, http.MethodPost, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/", c.HostURL, c.CompanyID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc == http.StatusBadRequest {
		apiError := c.parseAPIError(body)

		err := json.Unmarshal(body, &apiError)
		if err == nil {
			return nil, errors.New(apiError)
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newData := CertReqDNS{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}

func (c *Client) DeleteDNSCertReq(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if sc != http.StatusNoContent {
		return fmt.Errorf("%d - DELETE Failed for ID %d: %s", sc, id, c.parseAPIError(body))
	}

	return nil
}

func (c *Client) UpdateDNSCertReq(certReqID int, credID int) error {
	data := map[string]any{
		"credential": credID,
	}

	req, err := c.prepareJSONRequest(data, http.MethodPut, fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyID, certReqID))
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc == http.StatusBadRequest {
		apiError := c.parseAPIError(body)

		err := json.Unmarshal(body, &apiError)
		if err == nil {
			return errors.New(apiError)
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	return nil
}

func (c *Client) GetCertReqHTTP(id int) (CertReqHTTP, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/autoprovisioning/%d/sslcertificaterequest/%d", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return CertReqHTTP{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CertReqHTTP{}, err
	}

	if sc != http.StatusOK {
		return CertReqHTTP{}, fmt.Errorf("failure retrieving HTTP Certificate Request with id %d. %s", id, c.parseAPIError(body))
	}

	data := CertReqHTTP{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return CertReqHTTP{}, err
	}

	return data, nil
}

func (c *Client) CreateHTTPCertReq(certreq any) (*CertReqHTTP, error) {
	req, err := c.prepareJSONRequest(certreq, http.MethodPost, fmt.Sprintf("%s/v1/autoprovisioning/%d/sslcertificaterequest/", c.HostURL, c.CompanyID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc == http.StatusBadRequest {
		apiError := c.parseAPIError(body)

		err := json.Unmarshal(body, &apiError)
		if err == nil {
			return nil, errors.New(apiError)
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newData := CertReqHTTP{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}
