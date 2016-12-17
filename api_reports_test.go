package marketparser

import (
	"testing"
)

// TestGetReports trying to get list of reports.
func TestGetReports(t *testing.T) {
	campaigns := getCampaigns(t)
	if len(campaigns) > 0 {
		reports, err := createClient(t).GetReports(campaigns[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Returnes %d reports", len(reports))
	}
}

// TestCreateReport trying to create report for first campaign.
func TestCreateReport(t *testing.T) {
	campaigns := getCampaigns(t)
	if len(campaigns) > 0 {
		reportId, err := createClient(t).CreateReport(campaigns[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Report ID: %d", reportId)
	}
}

func TestGetReportResults(t *testing.T) {
	campaigns := getCampaigns(t)
	if len(campaigns) > 0 {
		reports, err := createClient(t).GetReports(campaigns[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		for i := range reports {
			if reports[i].Status == ReportStatusOk {
				reportProducts, err := createClient(t).GetReportResults(campaigns[0].ID, reports[i].ID)
				if err != nil {
					t.Fatal(err)
				}

				t.Logf("Returns %d products", len(reportProducts))

				break
			}
		}
	}
}