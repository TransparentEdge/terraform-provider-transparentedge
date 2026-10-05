package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// /v1/media/ is not scoped by API environment (staging/production is separated by
// company_id instead), so there is no staging twin of these methods.

var (
	// ErrTranscodingProfileNotFound is returned by GetTranscodingProfile when the API
	// answers 404, so that the resource can drop the profile from the state instead of
	// failing every subsequent plan.
	ErrTranscodingProfileNotFound = errors.New("transcoding profile not found")

	// ErrTranscodingProfileNotRead is returned by CreateTranscodingProfile when the profile
	// was created but the follow-up GET that hydrates custom_profiles failed. The profile
	// returned alongside it is the create response, and the caller must still write it to
	// the state: reporting a plain error would leave the profile orphaned in the API and the
	// next apply would create a duplicate.
	ErrTranscodingProfileNotRead = errors.New("transcoding profile created but could not be read back")
)

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

	// The create response (TranscodingProfile) does not include custom_profiles, unlike
	// the get/update response (TranscodingProfileDetail). Fetch it to hydrate the state.
	// If that GET fails the profile has already been created, so return the create response
	// with ErrTranscodingProfileNotRead instead of discarding it: the caller needs the id to
	// write the state, or the profile is orphaned and the next apply duplicates it.
	detail, err := c.GetTranscodingProfile(newData.ID)
	if err != nil {
		return &newData, fmt.Errorf("%w: %w", ErrTranscodingProfileNotRead, err)
	}

	return detail, nil
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

func (c *Client) GetTranscodingAllowedValues() (*TranscodingAllowedValuesAPIModel, error) {
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

	return &data, nil
}
