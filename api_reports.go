package marketparser

import (
	"encoding/json"
	"fmt"
	"time"
)

// ReportStatus is a type for report status
type ReportStatus string

var (
	ReportStatusWaiting                    ReportStatus = "WAITING"
	ReportStatusWaitingForRedownloadErrors ReportStatus = "WAITING_FOR_REDOWNLOAD_ERRORS"
	ReportStatusInProgress                 ReportStatus = "IN_PROGRESS"
	ReportStatusOk                         ReportStatus = "OK"
	ReportStatusError                      ReportStatus = "ERROR"
)

// Report is a struct contains report details
type Report struct {
	ID                     int          `json:"id"`
	Status                 ReportStatus `json:"status"`
	CreatedAt              time.Time    `json:"createdAt"`
	IsSuccessfullyFinished bool         `json:"isSuccessfullyFinished"`
	StartedAt              time.Time    `json:"startedAt"`
	FinishedAt             time.Time    `json:"finishedAt"`
	CountErrorProducts     int          `json:"countErrorProducts"`
	CountOkProducts        int          `json:"countOkProducts"`
}

func (c *Client) GetReports(campaignId int) ([]*Report, error) {
	if campaignId < 1 {
		return nil, fmt.Errorf("passed empty campaign ID: %d", campaignId)
	}

	reports := make([]*Report, 0)

	for page := 0; ; page++ {
		body, err := c.get(fmt.Sprintf("/campaigns/%d/reports.json", campaignId), page)
		if err != nil {
			return nil, err
		}

		var response struct {
			Response struct {
				Total   int       `json:"total"`
				Reports []*Report `json:"reports"`
			} `json:"response"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		reports = append(reports, response.Response.Reports...)

		if response.Response.Total <= len(reports) {
			break
		}
	}

	return reports, nil
}

func (c *Client) GetReportDetails(campaignId, reportId int) (*Report, error) {
	if campaignId < 1 || reportId < 1 {
		return nil, fmt.Errorf("either campaign ID or report ID is empty")
	}

	resBody, err := c.get(fmt.Sprintf("/campaigns/%d/reports/%d.json", campaignId, reportId), 0)
	if err != nil {
		return nil, err
	}

	var res struct{
		Response *Report `json:"response"`
	}

	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling error: %s", err)
	}

	return res.Response, nil
}

// CreateReport creates a new report and returns its id or error.
func (c *Client) CreateReport(campaignId int) (int, error) {
	if campaignId < 1 {
		return 0, fmt.Errorf("passed empty campaign ID: %d", campaignId)
	}

	resBody, err := c.post(fmt.Sprintf("/campaigns/%d/reports.json", campaignId), nil)
	if err != nil {
		return 0, err
	}

	var res struct {
		Response struct {
			ID int `json:"id"`
		} `json:"response"`
	}

	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return 0, fmt.Errorf("marshaling error: %s", err)
	}

	return res.Response.ID, nil
}

type ReportProductOffer struct {
	ShopName         string `json:"shopName"`
	Price            int    `json:"price"`
	DeliveryPrice    int    `json:"deliveryPrice"`
	InStock          bool   `json:"inStock"`
	Pickup           bool   `json:"pickup"`
	ProducerWarranty bool   `json:"producerWarranty"`
	LinkToOffer      string `json:"linkToOffer"`
}

type ReportProduct struct {
	Name             string                `json:"name"`
	YandexModelId    int                   `json:"yandexModelId"`
	YandexRegionId   int                   `json:"yandexRegionId"`
	YandexRegionName string                `json:"yandexRegionName"`
	OurId            string                `json:"ourId"`
	OurCost          int                   `json:"ourCost"`
	MedianPrice      int                   `json:"medianPrice"`
	MaxPrice         int                   `json:"maxPrice"`
	MinPrice         int                   `json:"minPrice"`
	AveragePrice     int                   `json:"averagePrice"`
	CountOffers      int                   `json:"countOffers"`
	Offers           []*ReportProductOffer `json:"offers"`
}

// GetReportResults returns list of products from the report.
func (c *Client) GetReportResults(campaignId, reportId int) ([]*ReportProduct, error) {
	if campaignId < 1 || reportId < 1 {
		return nil, fmt.Errorf("either campaign ID or report ID is empty")
	}

	reportProducts := make([]*ReportProduct, 0)

	for page := 1; ; page++ {
		resBody, err := c.get(fmt.Sprintf("/campaigns/%d/reports/%d/results.json", campaignId, reportId), page)
		if err != nil {
			return nil, err
		}

		var res struct {
			Response struct {
				Total    int              `json:"total"`
				Products []*ReportProduct `json:"products"`
			} `json:"response"`
		}

		err = json.Unmarshal(resBody, &res)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling error: %s", err)
		}

		reportProducts = append(reportProducts, res.Response.Products...)

		if res.Response.Total <= len(reportProducts) {
			break
		}
	}

	return reportProducts, nil
}
