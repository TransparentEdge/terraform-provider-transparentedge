package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var tfProviderSuffixRe = regexp.MustCompile(`\s{0,1}\[Terraform/[^\]]+\]$`)

func (c *Client) GetVclConfs(offset int, environment APIEnvironment) ([]VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/%s/%d/config/?offset=%d", c.HostURL, envpath, c.CompanyID, offset), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure while retrieving the list of configurations: %s", c.parseAPIError(body))
	}

	vclconfs := []VCLConfAPIModel{}

	err = json.Unmarshal(body, &vclconfs)
	if err != nil {
		return nil, err
	}

	return vclconfs, nil
}

func (c *Client) GetActiveVCLConf(environment APIEnvironment) (*VCLConfAPIModel, error) {
	confs, err := c.GetVclConfs(1, environment)
	if err != nil {
		return nil, err
	}

	top := VCLConfAPIModel{ID: -1}
	for _, vc := range confs {
		if vc.ID > top.ID {
			top = vc
		}
	}

	if top.ID <= 0 {
		return nil, errors.New("no VCL configurations found")
	}

	// remove version suffix
	top.Comment = stripProviderSuffix(top.Comment)

	return &top, nil
}

func (c *Client) GetVCLConfByID(environment APIEnvironment, id int) (*VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/%s/%d/config/%d", c.HostURL, envpath, c.CompanyID, id), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure while retrieving the list of configurations: %s", c.parseAPIError(body))
	}

	conf := VCLConfAPIModel{}

	err = json.Unmarshal(body, &conf)
	if err != nil {
		return nil, err
	}

	// remove version suffix
	conf.Comment = stripProviderSuffix(conf.Comment)

	return &conf, nil
}

func (c *Client) CreateVclconf(vclconf NewVCLConfAPIModel, environment APIEnvironment) (*VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	// append version suffix
	vclconf.Comment = appendProviderSuffix(vclconf.Comment, c.ProviderVersion)

	req, err := c.prepareJSONRequest(vclconf, "POST", fmt.Sprintf("%s/v1/%s/%d/config/", c.HostURL, envpath, c.CompanyID))
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
			return nil, fmt.Errorf("VCL COMPILATION ERROR\n\n%s\n", apiError) // nolint
		}
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newVclConf := VCLConfAPIModel{}

	err = json.Unmarshal(body, &newVclConf)
	if err != nil {
		return nil, err
	}

	// remove version suffix
	newVclConf.Comment = stripProviderSuffix(newVclConf.Comment)

	return &newVclConf, nil
}

// appendProviderSuffix adds the provider version suffix to a comment.
func appendProviderSuffix(comment, version string) string {
	clean := tfProviderSuffixRe.ReplaceAllString(comment, "")
	suffix := " [Terraform/" + version + "]"

	return clean + suffix
}

// stripProviderSuffix removes the provider suffix from a comment if present.
func stripProviderSuffix(comment string) string {
	return tfProviderSuffixRe.ReplaceAllString(comment, "")
}
