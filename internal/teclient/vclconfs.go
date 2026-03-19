package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

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
	vclconfs, err := c.GetVclConfs(1, environment)
	if err != nil {
		return nil, err
	}

	topVclConf := VCLConfAPIModel{ID: -1}
	for _, vclconf := range vclconfs {
		if vclconf.ID > topVclConf.ID {
			topVclConf = vclconf
		}
	}

	if topVclConf.ID <= 0 {
		return nil, errors.New("no VCL configurations found")
	}

	return &topVclConf, nil
}

func (c *Client) CreateVclconf(vclconf NewVCLConfAPIModel, environment APIEnvironment) (*VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	// Add a fixed comment to annotate that this configuration is being managed by terraform
	// consider porting this to the state in the future
	vclconf.Comment = "Managed with " + c.UserAgent

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

	return &newVclConf, nil
}
