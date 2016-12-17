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

func TestUpdatePrice(t *testing.T) {
	campaigns, err := createClient(t).GetCampaigns()
	if err != nil {
		t.Error(err)
	}

	if len(campaigns) > 0 {
		products := []*PriceProduct{
			{
				ID: "product-1",
				Name: "Product 1",
				Cost: 1234,
				YandexModelID: 11929063,
			},
			{
				ID: "product-2",
				Name: "Product 2",
				Cost: 4321,
				YandexModelID: 11929061,
			},
		}

		err = createClient(t).UpdatePrice(campaigns[0].ID, products)
		if err != nil {
			t.Error(err)
		}
	}
}
