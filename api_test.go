package marketparser

import (
	"testing"
)

func TestGetCampaigns(t *testing.T) {
	campaigns, err := createClient(t).GetCampaigns()
	if err != nil {
		t.Error(err)
	}

	for i := range campaigns {
		t.Logf("%d) %#v", i, campaigns[i])
	}
}