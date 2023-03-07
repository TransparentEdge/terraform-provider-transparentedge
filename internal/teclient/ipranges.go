package teclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetIPRanges() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/companies/ipranges", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Couldn't retrieve IP Ranges: %s", c.parseAPIError(body))
	}

	ranges := []string{}
	if err := json.Unmarshal(body, &ranges); err != nil {
		return nil, err
	}

	return ranges, nil
}
