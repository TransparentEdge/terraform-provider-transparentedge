package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetVclConfs(offset int, environment APIEnvironment) ([]VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s/%d/config/?offset=%d", c.HostURL, envpath, c.CompanyId, offset), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the list of configurations: %s", c.parseAPIError(body))
	}

	vclconfs := []VCLConfAPIModel{}
	if err := json.Unmarshal(body, &vclconfs); err != nil {
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
		return nil, fmt.Errorf("No VCL configurations found.")
	}

	return &topVclConf, nil
}

func (c *Client) CreateVclconf(vclconf NewVCLConfAPIModel, environment APIEnvironment) (*VCLConfAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	// Add a fixed comment to annotate that this configuration is being managed by terraform
	// TODO: consider porting this to the state in the future
	vclconf.Comment = fmt.Sprintf("Managed with %s", c.UserAgent)

	req, err := c.prepareJSONRequest(vclconf, "POST", fmt.Sprintf("%s/v1/%s/%d/config/", c.HostURL, envpath, c.CompanyId))
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
			return nil, fmt.Errorf("VCL COMPILATION ERROR\n\n%s\n", apiError)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newVclConf := VCLConfAPIModel{}
	if err := json.Unmarshal(body, &newVclConf); err != nil {
		return nil, err
	}

	return &newVclConf, nil
}
