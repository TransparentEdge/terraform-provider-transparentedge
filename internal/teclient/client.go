package teclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultAPIHTTPTimeout time.Duration = 50 * time.Second
)

func NewClient(host *string, companyid *int, clientid *string, clientsecret *string, insecure *bool, auth *bool, useragent *string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	token := TokenStruct{}

	c := Client{
		HTTPClient: &http.Client{Timeout: defaultAPIHTTPTimeout, Transport: tr},
		Token:      token,

		HostURL:      *host,
		CompanyId:    *companyid,
		ClientId:     *clientid,
		ClientSecret: *clientsecret,
		VerifySSL:    *insecure,
		UserAgent:    *useragent,
	}

	if *auth {
		if err := c.getToken(); err != nil {
			return nil, err
		}
	} else {
		c.Token.Token = "noauth"
	}

	return &c, nil
}

func (c *Client) getToken() error {
	req_body := []byte(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", c.ClientId, c.ClientSecret))

	req, err := http.NewRequest("POST", c.HostURL+"/v1/oauth2/access_token/", bytes.NewBuffer(req_body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if resp.StatusCode == 401 {
		return fmt.Errorf("%d - Could not create API client, please ensure credentials are correct.", resp.StatusCode)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		resp_body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%d - Failure creating API Token: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("%d - Unable to create API Token: %s", resp.StatusCode, c.parseAPIError(resp_body))
	}

	err_decode := json.NewDecoder(resp.Body).Decode(&c.Token)
	if err_decode != nil || c.Token.Token == "" {
		return fmt.Errorf("Fatal error creating API Token: %s", err_decode.Error())
	}

	return err
}

func (c *Client) doRequest(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.Token))
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, err
}

func (c *Client) prepareJSONRequest(jdata interface{}, method string, url string) (*http.Request, error) {
	data, err := json.Marshal(jdata)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) MustGetAPIEnvironmentPath(environment APIEnvironment) string {
	if environment == ProdEnv {
		return "autoprovisioning"
	} else if environment == StagingEnv {
		return "staging"
	}

	// Invalid env
	panic(fmt.Sprintf("Invalid environment: %+v", environment))
}

func (c *Client) parseAPIError(body []byte) string {
	apiError := APIMessage{}
	apiDetail := APIDetail{}

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&apiError); err == nil {
		return apiError.Message
	}

	decoder = json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&apiDetail); err == nil {
		return apiDetail.Detail
	}

	return string(body)
}
