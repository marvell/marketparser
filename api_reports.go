package marketparser

import (
	"time"
	"fmt"
	"encoding/json"
)

// ReportStatus is a type for report status
type ReportStatus string

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

func (c *client) GetReports(campaignId int) ([]*Report, error) {
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
