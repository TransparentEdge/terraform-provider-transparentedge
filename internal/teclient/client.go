package teclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultAPIHTTPTimeout time.Duration = 50 * time.Second
)

type TokenStruct struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
	Scope     string `json:"scope"`
}

type Client struct {
	HTTPClient *http.Client

	HostURL      string
	CompanyId    int
	ClientId     string
	ClientSecret string
	VerifySSL    bool
	UserAgent    string

	Token TokenStruct
}

func NewClient(host *string, companyid *int, clientid *string, clientsecret *string, verifyssl bool, useragent *string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: verifyssl},
		Proxy:           http.ProxyFromEnvironment,
	}

	c := Client{
		HTTPClient: &http.Client{Timeout: defaultAPIHTTPTimeout, Transport: tr},

		HostURL:      *host,
		CompanyId:    *companyid,
		ClientId:     *clientid,
		ClientSecret: *clientsecret,
		VerifySSL:    verifyssl,
		UserAgent:    *useragent,
	}

	if err := c.GetToken(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) GetToken() error {
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
			return fmt.Errorf("%d - Couldn't create an API Token: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("%d - Couldn't create an API token: %s", resp.StatusCode, string(resp_body))
	}

	if err_decode := json.NewDecoder(resp.Body).Decode(&c.Token); err_decode != nil || c.Token.Token == "" {
		return fmt.Errorf("Couldn't process the API response in order to create the token.")
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

func (c *Client) GetAPIEnvironmentPath(environment APIEnvironment) string {
	if environment == ProdEnv {
		return "autoprovisioning"
	} else if environment == StagingEnv {
		return "staging"
	}

	// requests will fail
	return "invalid_env"
}
