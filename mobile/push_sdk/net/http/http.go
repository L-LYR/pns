package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Client struct {
	base string
	c    *http.Client
}

type Payload map[string]interface{}

func MustNewHTTPClient(base string) *Client {
	return &Client{
		base: base,
		c: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
			},
		},
	}
}

func (c *Client) POST(path string, payload Payload) (Payload, error) {
	return c.Do(c.base+path, http.MethodPost, payload)
}

func (c *Client) PATCH(path string, payload Payload) (Payload, error) {
	return c.Do(c.base+path, http.MethodPatch, payload)
}

func (c *Client) PUT(path string, payload Payload) (Payload, error) {
	return c.Do(c.base+path, http.MethodPut, payload)
}

func (c *Client) GET(path string, payload Payload) (Payload, error) {
	return c.Do(c.base+path, http.MethodGet, payload)
}

func (c *Client) DELETE(path string, payload Payload) (Payload, error) {
	return c.Do(c.base+path, http.MethodDelete, payload)
}

func (c *Client) Do(URL string, Method string, payload Payload) (Payload, error) {
	type responseBody struct {
		ErrorCode    int     `json:"errorCode"`
		ErrorMessage string  `json:"errorMessage"`
		Payload      Payload `json:"payload"`
	}
	result := &responseBody{}
	body, err := jsoniter.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(Method, URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	response, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("network error")
	}
	rawResult, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := jsoniter.Unmarshal(rawResult, result); err != nil {
		return nil, err
	}
	if result.ErrorCode != 0 {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Payload, nil
}
