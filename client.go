package marketparser

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"bytes"
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

func (c *client) prepareUrl(urlPath string, pageNumber int) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("wrong base URL: %s", err)
	}

	u.Path += urlPath

	q := u.Query()
	if pageNumber > 1 {
		q.Set("page", fmt.Sprint(pageNumber))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *client) makeRequest(method string, u string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %s", err)
	}

	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *client) parseError(body []byte) error {
	var errorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	err := json.Unmarshal(body, &errorResponse)
	if err != nil {
		return fmt.Errorf("got error while trying to unmarshaling error message: %s", err)
	}

	return fmt.Errorf("api error: %d: %s", errorResponse.Code, errorResponse.Message)
}

func (c *client) get(urlPath string, pageNumber int) ([]byte, error) {
	u, err := c.prepareUrl(urlPath, pageNumber)
	if err != nil {
		return nil, err
	}

	req, err := c.makeRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	c.debug("URL: GET %s", req.URL.String())
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
		return nil, c.parseError(body)
	}

	return body, nil
}

func (c *client) post(urlPath string, requestBody []byte) ([]byte, error) {
	u, err := c.prepareUrl(urlPath, 1)
	if err != nil {
		return nil, err
	}

	req, err := c.makeRequest("POST", u, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	c.debug("URL: POST %s", req.URL.String())
	c.debug("HEADERS: %v", req.Header)
	c.debug("BODY: %q", requestBody)

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

	c.debug("RES BODY: %q", body)

	if res.StatusCode != 200 {
		return nil, c.parseError(body)
	}

	return body, nil
}