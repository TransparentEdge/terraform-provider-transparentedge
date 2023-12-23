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

func (c *Client) GetCRDNSCredential(id int) (CRDNSCredential, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d/", c.HostURL, c.CompanyId, id), nil)
	if err != nil {
		return CRDNSCredential{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CRDNSCredential{}, err
	}
	if sc != 200 {
		return CRDNSCredential{}, fmt.Errorf("Failure retrieving DNS Credential with id %d. %s", id, c.parseAPIError(body))
	}

	data := CRDNSCredential{}
	if err := json.Unmarshal(body, &data); err != nil {
		return CRDNSCredential{}, err
	}

	return data, nil
}

func (c *Client) CreateDNSCredential(dns_credential NewCRDNSCredential) (*CRDNSCredential, error) {
	req, err := c.prepareJSONRequest(dns_credential, "POST", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := c.parseAPIError(body)
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, fmt.Errorf(apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	new_data := CRDNSCredential{}
	if err := json.Unmarshal(body, &new_data); err != nil {
		return nil, err
	}

	return &new_data, nil
}

func (c *Client) UpdateDNSCredential(dns_credential NewCRDNSCredential, id int) (*CRDNSCredential, error) {
	req, err := c.prepareJSONRequest(dns_credential, "PUT", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d", c.HostURL, c.CompanyId, id))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := c.parseAPIError(body)
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, fmt.Errorf(apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	new_data := CRDNSCredential{}
	if err := json.Unmarshal(body, &new_data); err != nil {
		return nil, err
	}

	return &new_data, nil
}

func (c *Client) DeleteCRCredential(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscredential/%d", c.HostURL, c.CompanyId, id), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyId, id), nil)
	if err != nil {
		return CertReqDNS{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CertReqDNS{}, err
	}
	if sc != 200 {
		return CertReqDNS{}, fmt.Errorf("Failure retrieving DNS Certificate Request with id %d. %s", id, c.parseAPIError(body))
	}

	data := CertReqDNS{}
	if err := json.Unmarshal(body, &data); err != nil {
		return CertReqDNS{}, err
	}

	return data, nil
}

func (c *Client) CreateDNSCertReq(certreq interface{}) (*CertReqDNS, error) {
	req, err := c.prepareJSONRequest(certreq, "POST", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := c.parseAPIError(body)
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, fmt.Errorf(apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	new_data := CertReqDNS{}
	if err := json.Unmarshal(body, &new_data); err != nil {
		return nil, err
	}

	return &new_data, nil
}

func (c *Client) DeteleDNSCertReq(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyId, id), nil)
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

func (c *Client) UpdateDNSCertReq(certreq_id int, credential_id int) error {
	data := map[string]interface{}{
		"credential": credential_id,
	}
	req, err := c.prepareJSONRequest(data, "PUT", fmt.Sprintf("%s/v1/autoprovisioning/%d/dnscertrequest/%d", c.HostURL, c.CompanyId, certreq_id))
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := c.parseAPIError(body)
		if err := json.Unmarshal(body, &apiError); err == nil {
			return fmt.Errorf(apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	return nil
}

func (c *Client) GetCertReqHTTP(id int) (CertReqHTTP, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslcertificaterequest/%d", c.HostURL, c.CompanyId, id), nil)
	if err != nil {
		return CertReqHTTP{}, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return CertReqHTTP{}, err
	}
	if sc != 200 {
		return CertReqHTTP{}, fmt.Errorf("Failure retrieving HTTP Certificate Request with id %d. %s", id, c.parseAPIError(body))
	}

	data := CertReqHTTP{}
	if err := json.Unmarshal(body, &data); err != nil {
		return CertReqHTTP{}, err
	}

	return data, nil
}

func (c *Client) CreateHTTPCertReq(certreq interface{}) (*CertReqHTTP, error) {
	req, err := c.prepareJSONRequest(certreq, "POST", fmt.Sprintf("%s/v1/autoprovisioning/%d/sslcertificaterequest/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := c.parseAPIError(body)
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, fmt.Errorf(apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	new_data := CertReqHTTP{}
	if err := json.Unmarshal(body, &new_data); err != nil {
		return nil, err
	}

	return &new_data, nil
}
