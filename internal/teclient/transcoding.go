package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// /v1/media/ is not scoped by API environment (staging/production is separated by
// company_id instead), so there is no staging twin of these methods.

// ErrTranscodingProfileNotFound is returned by GetTranscodingProfile when the API
// answers 404, so that the resource can drop the profile from the state instead of
// failing every subsequent plan.
var ErrTranscodingProfileNotFound = errors.New("transcoding profile not found")

func (c *Client) GetTranscodingProfiles() ([]TranscodingProfileAPIModel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/media/%d/transcoding_profiles/", c.HostURL, c.CompanyID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure retrieving Transcoding Profiles: %s", c.parseAPIError(body))
	}

	profiles := []TranscodingProfileAPIModel{}

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (c *Client) GetTranscodingProfile(id int) (*TranscodingProfileAPIModel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/media/%d/transcoding_profile/%d/", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc == http.StatusNotFound {
		return nil, fmt.Errorf("%w: id %d", ErrTranscodingProfileNotFound, id)
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure retrieving Transcoding Profile with id %d. %s", id, c.parseAPIError(body))
	}

	data := TranscodingProfileAPIModel{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) CreateTranscodingProfile(p NewTranscodingProfileAPIModel) (*TranscodingProfileAPIModel, error) {
	req, err := c.prepareJSONRequest(p, http.MethodPost, fmt.Sprintf("%s/v1/media/%d/transcoding_profiles/", c.HostURL, c.CompanyID))
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

	newData := TranscodingProfileAPIModel{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}

func (c *Client) UpdateTranscodingProfile(p NewTranscodingProfileAPIModel, id int) (*TranscodingProfileAPIModel, error) {
	req, err := c.prepareJSONRequest(p, http.MethodPut, fmt.Sprintf("%s/v1/media/%d/transcoding_profile/%d/", c.HostURL, c.CompanyID, id))
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

	newData := TranscodingProfileAPIModel{}

	err = json.Unmarshal(body, &newData)
	if err != nil {
		return nil, err
	}

	return &newData, nil
}

func (c *Client) DeleteTranscodingProfile(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/media/%d/transcoding_profile/%d/", c.HostURL, c.CompanyID, id), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	// A 404 means the profile is already gone (for example, deleted from the dashboard
	// between the refresh and the apply), which is the desired outcome of a delete.
	if sc != http.StatusNoContent && sc != http.StatusNotFound {
		return fmt.Errorf("%d - DELETE Failed for ID %d: %s", sc, id, c.parseAPIError(body))
	}

	return nil
}

// GetTranscodingAllowedValues is cached for the lifetime of the client (one Terraform
// command), since every transcoding profile in a plan checks its values against it. The
// lock also collapses the concurrent first calls into a single request. Errors are not
// cached, so a transient failure is retried by the next caller.
func (c *Client) GetTranscodingAllowedValues() (*TranscodingAllowedValuesAPIModel, error) {
	c.allowedValuesMu.Lock()
	defer c.allowedValuesMu.Unlock()

	if c.allowedValues != nil {
		return c.allowedValues, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/media/%d/allowed_values/", c.HostURL, c.CompanyID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure retrieving Transcoding Allowed Values: %s", c.parseAPIError(body))
	}

	data := TranscodingAllowedValuesAPIModel{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	c.allowedValues = &data

	return c.allowedValues, nil
}
