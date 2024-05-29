package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CustomSearchClient struct {
	APIKey string
	CX     string
}

func NewCustomSearchClient(apiKey, cx string) *CustomSearchClient {
	return &CustomSearchClient{APIKey: apiKey, CX: cx}
}

func (c *CustomSearchClient) Search(query string) ([]string, error) {
	escapedQuery := url.QueryEscape(query)
	requestURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s", escapedQuery, c.CX, c.APIKey)

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Custom Search API error: %s, body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	var links []string
	for _, item := range items {
		if entry, ok := item.(map[string]interface{}); ok {
			if link, ok := entry["link"].(string); ok {
				links = append(links, link)
			}
		}
	}

	return links, nil
}
