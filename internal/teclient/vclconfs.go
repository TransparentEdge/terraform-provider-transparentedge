package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetVclConfs() ([]VCLConfAPIModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/autoprovisioning/%d/config/", c.HostURL, c.CompanyId), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the list of configurations: %s", string(body))
	}

	vclconfs := []VCLConfAPIModel{}
	if err := json.Unmarshal(body, &vclconfs); err != nil {
		return nil, err
	}

	return vclconfs, nil
}

func (c *Client) GetActiveVCLConf() (*VCLConfAPIModel, error) {
	vclconfs, err := c.GetVclConfs()
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

func (c *Client) CreateVclconf(vclconf NewVCLConfAPIModel) (*VCLConfAPIModel, error) {
	req, err := c.preparePostRequest(vclconf, fmt.Sprintf("%s/v1/autoprovisioning/%d/config/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if sc == 400 {
		apiError := ErrorAPIMessage{}
		if err := json.Unmarshal(body, &apiError); err == nil {
			return nil, fmt.Errorf("VCL COMPILATION ERROR\n\n%s\n", apiError.Message)
		}
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, string(body))
	}

	newVclConf := VCLConfAPIModel{}
	if err := json.Unmarshal(body, &newVclConf); err != nil {
		return nil, err
	}

	return &newVclConf, nil
}
