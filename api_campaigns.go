package marketparser

import (
	"encoding/json"
	"fmt"
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
	var campaigns []*Campaign

	for page := 1; ; page++ {
		body, err := c.get("/campaigns.json", page)
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

		campaigns = append(campaigns, response.Response.Campaigns...)

		if response.Response.Total <= len(campaigns) {
			break
		}
	}

	return campaigns, nil
}

// PriceProduct describes product details in campaing price
type PriceProduct struct {
	ID            string            `json:"id"`
	Name          string            `json:"name,omitempty"`
	Cost          int               `json:"cost,omitempty"`
	YandexModelID int               `json:"yandex_model_id,omitempty"`
	Custom        map[string]string `json:"custom,omitempty"`
}

// UpdatePrice updates price for the campaign
func (c *client) UpdatePrice(campaignId int, products []*PriceProduct) error {
	if campaignId < 1 {
		return fmt.Errorf("passed empty campaign ID: %d", campaignId)
	}

	if len(products) < 1 {
		return fmt.Errorf("must be passed al least one product")
	}

	requestBody, err := json.Marshal(struct {
		Products []*PriceProduct `json:"products"`
	}{Products: products})
	if err != nil {
		return err
	}

	body, err := c.post(fmt.Sprintf("/campaigns/%d/price.json", campaignId), requestBody)
	if err != nil {
		return err
	}

	var response struct {
		Response struct {
			Success bool `json:"success"`
		} `json:"response"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("wrong response after updated price: %s", err)
	}

	if response.Response.Success != true {
		return fmt.Errorf("unsuccessfully attemtp to update price")
	}

	return nil
}

// PriceStatus is a type for price status
type PriceStatus string

var (
	// PriceStatusError is an error status of the price
	PriceStatusError PriceStatus = "ERROR"

	// PriceStatusParsed is a status shows that parsing already done
	PriceStatusParsed PriceStatus = "PARSED"

	PriceStatusReadyToBeParsed  PriceStatus = "READY_TO_BE_PARSED"
	PriceStatusSearchInProgress PriceStatus = "SEARCH_IN_PROGRESS"

	// PriceStatusProcessed is a status shows that all of work under the price already done
	PriceStatusProcessed     PriceStatus = "PROCESSED"
	PriceStatusTooLowBalance PriceStatus = "NOT_ENOUGH_BALANCE_TO_PROCESS"
)

type PriceDetails struct {
	Id                       int         `json:"id"`
	CreatedAt                time.Time   `json:"createdAt"`
	Status                   PriceStatus `json:"status"`
	CountNotEmptyRows        int         `json:"countNotEmptyRows"`
	CountFoundDuplicatedRows int         `json:"countFoundDuplicatedRows"`
	IsSuccessfullyProcessed  bool        `json:"isSuccessfullyProcessed"`
}

// GetPriceDetails returns price details or error
func (c *client) GetPriceDetails(campaignId int) (*PriceDetails, error) {
	if campaignId < 1 {
		return nil, fmt.Errorf("passed empty campaign ID: %d", campaignId)
	}

	body, err := c.get(fmt.Sprintf("/campaigns/%d/price.json", campaignId), 0)
	if err != nil {
		return nil, err
	}

	var response struct {
		Response *PriceDetails `json:"response"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Response, nil
}
