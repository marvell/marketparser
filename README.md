# Go-client library for MarketParser API

[DOC](https://godoc.org/github.com/marvell/marketparser)

## Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/marvell/marketparser"
)

var campaignId = <CAMPAING_ID>

var insterstedProducts = map[int]string{
	12912675: "Google Chromecast Audio",
	11919264: "Audio-Technica ATH-MSR7",
	14207067: "Apple Watch Series 2 42mm",
	12407907: "Cambridge Audio CXA60",
}

func main() {
	client, err := marketparser.NewClient("YOUR_API_KEY")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	priceProducts := make([]*marketparser.PriceProduct, 0)
	for yandexId, productName := range insterstedProducts {
		priceProducts = append(priceProducts, &marketparser.PriceProduct{
			ID:            fmt.Sprint(yandexId),
			Name:          productName,
			YandexModelID: yandexId,
		})
	}

	err = client.UpdatePrice(campaignId, priceProducts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		priceDetails, err := client.GetPriceDetails(campaignId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if priceDetails.Status == marketparser.PriceStatusProcessed {
			break
		}

		fmt.Printf("price isn't ready yet: %s\n", priceDetails.Status)
		time.Sleep(2 * time.Second)
	}

	reportId, err := client.CreateReport(campaignId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		report, err := client.GetReportDetails(campaignId, reportId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if report.Status == marketparser.ReportStatusOk {
			break
		}

		fmt.Printf("report isn't ready yet: %s\n", report.Status)
		time.Sleep(5 * time.Second)
	}

	reportResults, err := client.GetReportResults(campaignId, reportId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, r := range reportResults {
		fmt.Printf("%d: %s\t%d\t%d\t%d\n", i, r.Name, r.MinPrice, r.AveragePrice, r.MaxPrice)
	}
}
```