package teclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetSiteVerifyString(siteDomain string) string {
	data := SiteVerifyStringAPIModelRequest{Domain: siteDomain}

	req, err := c.prepareJSONRequest(data, "POST", fmt.Sprintf("%s/v1/companies/%d/siteverification/", c.HostURL, c.CompanyID))
	if err != nil {
		return ""
	}

	body, sc, err := c.doRequest(req)
	if err != nil || sc != 200 {
		return ""
	}

	svsResp := SiteVerifyStringAPIModelResponse{}

	err = json.Unmarshal(body, &svsResp)
	if err != nil {
		return ""
	}

	return svsResp.Txt
}

func (c *Client) GetSites() ([]SiteAPIModel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/companies/%d/sites/", c.HostURL, c.CompanyID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure while retrieving the list of sites: %s", c.parseAPIError(body))
	}

	sites := []SiteAPIModel{}

	err = json.Unmarshal(body, &sites)
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func (c *Client) GetSite(siteID int) (*SiteAPIModel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/companies/%d/sites/%d/", c.HostURL, c.CompanyID, siteID), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	if sc != http.StatusOK {
		return nil, fmt.Errorf("failure while retrieving the site with ID %d: %s", siteID, c.parseAPIError(body))
	}

	site := SiteAPIModel{}

	err = json.Unmarshal(body, &site)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (c *Client) CreateSite(site SiteNewAPIModel) (*SiteAPIModel, bool, error) {
	// returns model, error, verify_error
	req, err := c.prepareJSONRequest(site, "POST", fmt.Sprintf("%s/v1/companies/%d/sites/", c.HostURL, c.CompanyID))
	if err != nil {
		return nil, false, err
	}

	body, sc, err := c.doRequest(req)
	if sc == http.StatusBadGateway {
		return nil, false, errors.New("failed to create the site - API error")
	}

	if sc == http.StatusForbidden {
		// Verification error
		msg := "Please ensure that the site can be verified with one of the following two options:\n" +
			"  * Option 1: A tcdn.txt file in the root of your site with the verification string\n" +
			"  * Option 2: A TXT record: _tcdn_challenge." + site.URL + " with the verification string\n"

		// Best effor here to show the user the verification string
		verifyString := c.GetSiteVerifyString(site.URL)
		if verifyString != "" {
			msg = msg + "\nThe verification string for this site is: " + verifyString + "\n"
		}

		msg = msg + "\nAPI Response: " + c.parseAPIError(body) + "\n\n" +
			"If you need to get the verification string again, run 'terraform plan' and 'terraform show'" +
			" with the datasource 'siteverify'.\nIn case of doubts please contact with support."

		return nil, true, fmt.Errorf("validation error:\n%s", msg)
	}

	if sc == http.StatusBadRequest {
		if strings.Contains(c.parseAPIError(body), "Site ownership denied") {
			// site belongs to another company
			return nil, false, fmt.Errorf("site not owned: %s", c.parseAPIError(body))
		}
		// check if the site already exists
		if existingSite := c.GetIfExists(body, site.URL); existingSite != nil {
			return existingSite, false, nil
		}
	}

	if err != nil {
		return nil, false, fmt.Errorf("%d - %s", sc, err.Error())
	}

	if sc != http.StatusOK && sc != http.StatusCreated { // 200 = new, 201 = activated again
		return nil, false, fmt.Errorf("%d - %s", sc, c.parseAPIError(body))
	}

	newSite := SiteAPIModel{}

	err = json.Unmarshal(body, &newSite)
	if err != nil {
		return nil, false, err
	}

	return &newSite, false, nil
}

func (c *Client) GetIfExists(body []byte, siteDomain string) *SiteAPIModel {
	errorMessage := c.parseAPIError(body)
	if strings.Contains(errorMessage, "already exists") {
		// Try to find the site
		sites, err := c.GetSites()
		if err == nil {
			for _, site := range sites {
				if site.URL == siteDomain {
					return &site
				}
			}
		}
	}

	return nil
}

func (c *Client) DeleteSite(siteID int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/companies/%d/sites/%d/", c.HostURL, c.CompanyID, siteID), nil)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}

	if sc != http.StatusNoContent {
		return fmt.Errorf("%d - API request failed trying to DELETE the site ID %d: %s", sc, siteID, c.parseAPIError(body))
	}

	return nil
}
