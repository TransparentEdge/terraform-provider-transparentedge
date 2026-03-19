package teclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetIPRanges() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/companies/ipranges", c.HostURL), nil) // nolint: perfsprint
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve IP Ranges: %s", c.parseAPIError(body))
	}

	ranges := []string{}

	err = json.Unmarshal(body, &ranges)
	if err != nil {
		return nil, err
	}

	return ranges, nil
}
