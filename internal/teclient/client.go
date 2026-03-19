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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure}, // nolint: gosec
		Proxy:           http.ProxyFromEnvironment,
	}

	token := TokenStruct{}

	c := Client{
		HTTPClient: &http.Client{Timeout: defaultAPIHTTPTimeout, Transport: tr},
		Token:      token,

		HostURL:      *host,
		CompanyID:    *companyid,
		ClientID:     *clientid,
		ClientSecret: *clientsecret,
		VerifySSL:    *insecure,
		UserAgent:    *useragent,
	}

	if *auth {
		err := c.getToken()
		if err != nil {
			return nil, err
		}
	} else {
		c.Token.Token = "noauth"
	}

	return &c, nil
}

func (c *Client) getToken() error {
	reqBody := fmt.Appendf(nil, "client_id=%s&client_secret=%s&grant_type=client_credentials", c.ClientID, c.ClientSecret)

	req, err := http.NewRequest(http.MethodPost, c.HostURL+"/v1/oauth2/access_token/", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("received status code %d, could not create API client, please ensure credentials are correct", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%d - Failure creating API Token: %w", resp.StatusCode, err)
		}

		return fmt.Errorf("received status code %d, unable to create API Token: %s", resp.StatusCode, c.parseAPIError(respBody))
	}

	errDecode := json.NewDecoder(resp.Body).Decode(&c.Token)
	if errDecode != nil || c.Token.Token == "" {
		return fmt.Errorf("fatal error creating API Token: %s", errDecode.Error())
	}

	return err
}

func (c *Client) doRequest(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token.Token)
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

func (*Client) prepareJSONRequest(jdata any, method string, url string) (*http.Request, error) {
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

func (*Client) MustGetAPIEnvironmentPath(environment APIEnvironment) string {
	switch environment {
	case ProdEnv:
		return "autoprovisioning"
	case StagingEnv:
		return "staging"
	}

	// Invalid env
	panic(fmt.Sprintf("Invalid environment: %+v", environment))
}

func (*Client) parseAPIError(body []byte) string {
	apiError := APIMessage{}
	apiDetail := APIDetail{}

	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&apiError)
	if err == nil {
		return apiError.Message
	}

	decoder = json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&apiDetail)
	if err == nil {
		return apiDetail.Detail
	}

	return string(body)
}
