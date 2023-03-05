package teclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetSiteVerifyString(site_domain string) string {
	data := SiteVerifyStringAPIModelRequest{Domain: site_domain}
	req, err := c.prepareJSONRequest(data, "POST", fmt.Sprintf("%s/v1/companies/%d/siteverification/", c.HostURL, c.CompanyId))
	if err != nil {
		return ""
	}

	body, sc, err := c.doRequest(req)
	if err != nil || sc != 200 {
		return ""
	}

	svsResp := SiteVerifyStringAPIModelResponse{}
	if err := json.Unmarshal(body, &svsResp); err != nil {
		return ""
	}

	return svsResp.Txt
}

func (c *Client) GetSites() ([]SiteAPIModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/companies/%d/sites/", c.HostURL, c.CompanyId), nil)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the list of sites: %s", string(body))
	}

	sites := []SiteAPIModel{}
	if err := json.Unmarshal(body, &sites); err != nil {
		return nil, err
	}

	return sites, nil
}

func (c *Client) GetSite(siteID int) (*SiteAPIModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/companies/%d/sites/%d/", c.HostURL, c.CompanyId, siteID), nil)
	req.Header.Set("User-Agent", defaultUserAgent)
	if err != nil {
		return nil, err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	if sc != 200 {
		return nil, fmt.Errorf("Couldn't retrieve the site with ID %d: %s", siteID, string(body))
	}

	site := SiteAPIModel{}
	if err := json.Unmarshal(body, &site); err != nil {
		return nil, err
	}

	return &site, nil
}

func (c *Client) CreateSite(site SiteNewAPIModel) (*SiteAPIModel, bool, error) {
	// returns model, error, verify_error
	req, err := c.prepareJSONRequest(site, "POST", fmt.Sprintf("%s/v1/companies/%d/sites/", c.HostURL, c.CompanyId))
	if err != nil {
		return nil, false, err
	}

	body, sc, err := c.doRequest(req)
	if sc == 502 {
		return nil, false, fmt.Errorf("502 - Error creating the site.")
	}
	if sc == 403 {
		// Verification error
		msg := "Please ensure that the site can be verified with one of the following two options:\n" +
			"  * Option 1: A tcdn.txt file in the root of your site with the verification string\n" +
			"  * Option 2: A TXT record: _tcdn_challenge." + site.Url + " with the verification string\n"

		// Best effor here to show the user the verification string
		verify_string := c.GetSiteVerifyString(site.Url)
		if verify_string != "" {
			msg = msg + "\nThe verification string for this site is: " + verify_string + "\n"
		}

		msg = msg + "\nAPI Response: " + string(body) + "\n\n" +
			"If you need to get the verification string again, run 'terraform plan' and 'terraform show'" +
			" with the datasource 'siteverify'.\nIn case of doubts please contact with support."

		return nil, true, fmt.Errorf("Validation error:\n%s", msg)
	}
	if sc == 400 {
		if strings.Contains(string(body), "Site ownership denied") {
			// site belongs to another company
			return nil, false, fmt.Errorf("Site not owned: %s", string(body))
		}
		// check if the site already exists
		if existingSite := c.GetIfExists(body, site.Url); existingSite != nil {
			return existingSite, false, nil
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("%d - %s", sc, err.Error())
	}
	if !(sc == 200 || sc == 201) { // 200 = new, 201 = activated again
		return nil, false, fmt.Errorf("%d - %s", sc, string(body))
	}

	newSite := SiteAPIModel{}
	if err := json.Unmarshal(body, &newSite); err != nil {
		return nil, false, err
	}

	return &newSite, false, nil
}

func (c *Client) GetIfExists(body []byte, site_domain string) *SiteAPIModel {
	errorMessage := ErrorAPIMessage{}
	if err := json.Unmarshal(body, &errorMessage); err == nil {
		if strings.Contains(errorMessage.Message, "already exists") {
			// Try to find the site
			if sites, err := c.GetSites(); err == nil {
				for _, site := range sites {
					if site.Url == site_domain {
						return &site
					}
				}
			}
		}
	}
	return nil
}

func (c *Client) DeleteSite(siteID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/companies/%d/sites/%d/", c.HostURL, c.CompanyId, siteID), nil)
	req.Header.Set("User-Agent", defaultUserAgent)
	if err != nil {
		return err
	}

	body, sc, err := c.doRequest(req)
	if err != nil {
		return err
	}
	if sc != 204 {
		return fmt.Errorf("%d - API request failed trying to DELETE the site ID %d: %s", sc, siteID, string(body))
	}

	return nil
}
