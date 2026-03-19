package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetBackend(backendID int, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyID, backendID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve the backend with ID %d: %s", backendID, c.parseAPIError(body))
	}

	backend := BackendAPIModel{}

	err = json.Unmarshal(body, &backend)
	if err != nil {
		return nil, err
	}

	return &backend, nil
}

func (c *Client) GetBackends(environment APIEnvironment) ([]BackendAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/%s/%d/backends/", c.HostURL, envpath, c.CompanyID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve the list of backends: %s", c.parseAPIError(body))
	}

	backends := []BackendAPIModel{}

	err = json.Unmarshal(body, &backends)
	if err != nil {
		return nil, err
	}

	return backends, nil
}

func (c *Client) GetBackendByName(name string, environment APIEnvironment) (*BackendAPIModel, error) {
	backends, err := c.GetBackends(environment)
	if err != nil {
		return nil, err
	}

	for _, backend := range backends {
		if backend.Name == name {
			return &backend, nil
		}
	}

	return nil, fmt.Errorf("no backend named '%s' found", name)
}

func (c *Client) CreateBackend(backend NewBackendAPIModel, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := c.prepareJSONRequest(backend, http.MethodPost, fmt.Sprintf("%s/v1/%s/%d/backends/", c.HostURL, envpath, c.CompanyID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newBackend := BackendAPIModel{}

	err = json.Unmarshal(body, &newBackend)
	if err != nil {
		return nil, err
	}

	return &newBackend, nil
}

func (c *Client) UpdateBackend(backend BackendAPIModel, environment APIEnvironment) (*BackendAPIModel, error) {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := c.prepareJSONRequest(backend, http.MethodPut, fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyID, backend.ID))
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc != http.StatusOK && sc != http.StatusCreated {
		return nil, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newBackend := BackendAPIModel{}

	err = json.Unmarshal(body, &newBackend)
	if err != nil {
		return nil, err
	}

	return &newBackend, nil
}

func (c *Client) DeleteBackend(backendID int, environment APIEnvironment) error {
	envpath := c.MustGetAPIEnvironmentPath(environment)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/%s/%d/backends/%d/", c.HostURL, envpath, c.CompanyID, backendID), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if sc == http.StatusForbidden {
		if strings.Contains(c.parseAPIError(body), "references in active config") {
			return errors.New("cannot delete a backend with references in the active autoprovisioning configuration, please remove all the references from the configuration first")
		}
	}

	if sc != http.StatusNoContent {
		return fmt.Errorf("%d - API request failed trying to DELETE the backend ID %d: %s", sc, backendID, c.parseAPIError(body))
	}

	return nil
}
