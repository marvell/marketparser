package marketparser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const baseUrl = "http://cp.marketparser.ru/api/v2"

type client struct {
	apiKey string

	httpClient *http.Client
	logger     *log.Logger
}

func NewClient(apiKey string) (*client, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("you must specified API-key")
	}

	return &client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

func (c *client) DebugMode() {
	c.logger = log.New(os.Stdout, "[marketparser] ", log.LstdFlags)
}

func (c *client) debug(f string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(f, args...)
	}
}

func (c *client) get(urlPath string, pageNumber int) ([]byte, error) {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("wrong base URL: %s", err)
	}

	requestUrl.Path += urlPath

	q := requestUrl.Query()
	if pageNumber > 1 {
		q.Set("page", fmt.Sprint(pageNumber))
	}
	requestUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %s", err)
	}

	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	c.debug("URL: %s", req.URL.String())
	c.debug("HEADERS: %v", req.Header)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while doing request: %s", err)
	}

	c.debug("RES CODE: %d", res.StatusCode)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response: %s", err)
	}

	if res.StatusCode != 200 {
		var errorResponse struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}

		c.debug("BODY: %s", body)

		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf("got error while trying to unmarshaling error message: %s", err)
		}

		return nil, fmt.Errorf("api error: %d: %s", errorResponse.Code, errorResponse.Message)
	}

	return body, nil
}
