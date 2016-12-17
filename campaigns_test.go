package marketparser

import (
	"testing"
)

func TestGetPriceDetails(t *testing.T) {
	campaigns, err := createClient(t).GetCampaigns()
	if err != nil {
		t.Fatal(err)
	}

	if len(campaigns) > 0 {
		priceDetails, err := createClient(t).GetPriceDetails(campaigns[0].ID)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("PRICE DETAILS: %#v", priceDetails)
	}
}
