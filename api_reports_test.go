package marketparser

import (
	"testing"
)

func TestGetReports(t *testing.T) {
	campaigns, err := createClient(t).GetCampaigns()
	if err != nil {
		t.Error(err)
	}

	if len(campaigns) > 0 {
		reports, err := createClient(t).GetReports(campaigns[0].ID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("%#v", reports[0])
	}
}
