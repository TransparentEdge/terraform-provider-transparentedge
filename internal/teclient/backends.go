package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetBackend(backendID int, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.GetAPIEnvironmentPath(environment)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyId, backendID), nil)
	req.Header.Set("User-Agent", USERAGENT)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the backend with ID %d: %s", backendID, string(body))
	}

	backend := BackendAPIModel{}
	if err := json.Unmarshal(body, &backend); err != nil {
		return nil, err
	}

	return &backend, nil
}

func (c *Client) GetBackends(environment APIEnvironment) ([]BackendAPIModel, error) {
	envpath := c.GetAPIEnvironmentPath(environment)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s/%d/backends/", c.HostURL, envpath, c.CompanyId), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the list of backends: %s", string(body))
	}

	backends := []BackendAPIModel{}
	if err := json.Unmarshal(body, &backends); err != nil {
		return nil, err
	}

	return backends, nil
}

func (c *Client) CreateBackend(backend NewBackendAPIModel, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.GetAPIEnvironmentPath(environment)
	req, err := c.preparePostRequest(backend, fmt.Sprintf("%s/v1/%s/%d/backends/", c.HostURL, envpath, c.CompanyId))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, string(body))
	}

	newBackend := BackendAPIModel{}
	if err := json.Unmarshal(body, &newBackend); err != nil {
		return nil, err
	}

	return &newBackend, nil
}

func (c *Client) UpdateBackend(backend BackendAPIModel, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.GetAPIEnvironmentPath(environment)
	req, err := c.preparePutRequest(backend, fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyId, backend.ID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if !(sc == 200 || sc == 201) {
		return nil, fmt.Errorf("%d - %s", sc, string(body))
	}

	newBackend := BackendAPIModel{}
	if err := json.Unmarshal(body, &newBackend); err != nil {
		return nil, err
	}

	return &newBackend, nil
}

func (c *Client) DeleteBackend(backendID int, environment APIEnvironment) error {
	envpath := c.GetAPIEnvironmentPath(environment)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyId, backendID), nil)
	req.Header.Set("User-Agent", USERAGENT)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}
	if sc == 403 {
		if strings.Contains(string(body), "references in active config") {
			return fmt.Errorf("Cannot delete a backend with references in the active autoprovisioning configuration, please remove all the references from the configuration first.")
		}
	}
	if sc != 204 {
		return fmt.Errorf("%d - API request failed trying to DELETE the backend ID %d: %s", sc, backendID, string(body))
	}

	return nil
}
