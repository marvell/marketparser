package marketparser

import (
	"encoding/json"
	"time"
)

// Campaign describes an item in campaign response
type Campaign struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	CreatedAt            time.Time `json:"createdAt"`
	ReadyToCreateReports bool      `json:"readyToCreateReports"`
}

// GetCampaigns returns list of campaigns or error. Uses /campaigns.json handler.
func (c *client) GetCampaigns() ([]*Campaign, error) {
	body, err := c.get("/campaigns.json", 1)
	if err != nil {
		return nil, err
	}

	var response struct {
		Response struct {
			Total     int         `json:"total"`
			Campaigns []*Campaign `json:"campaigns"`
		} `json:"response"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Response.Campaigns, nil
}

