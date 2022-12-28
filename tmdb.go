// Copyright (c) 2019 Cyro Dubeux. License MIT.

// Package tmdb is a wrapper for working with TMDb API.
package tmdb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

// TMDb constants
const (
	baseURL           = "https://api.themoviedb.org/3"
	permissionURL     = "https://www.themoviedb.org/authenticate/"
	authenticationURL = "/authentication/"
	movieURL          = "/movie/"
	tvURL             = "/tv/"
	tvSeasonURL       = "/season/"
	tvEpisodeURL      = "/episode/"
	personURL         = "/person/"
	searchURL         = "/search/"
	collectionURL     = "/collection/"
	companyURL        = "/company/"
	configurationURL  = "/configuration/"
	creditURL         = "/credit/"
	discoverURL       = "/discover/"
	networkURL        = "/network/"
	keywordURL        = "/keyword/"
	genreURL          = "/genre/"
	guestSessionURL   = "/guest_session/"
	listURL           = "/list/"
	accountURL        = "/account/"
	watchProvidersURL = "/watch/providers/"
)

// Client type is a struct to instantiate this pkg.
type Client struct {
	// TMDb apiKey to use the client.
	apiKey string
	// sessionId to use the client.
	sessionID string
	// Auto retry flag to indicates if the client
	// should retry the previous operation.
	autoRetry bool
	// resty.Client for custom configuration.
	http *resty.Client
}

// Response type is a struct for http responses.
type Response struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// Init setups the resty Client with an apiKey.
func Init(apiKey string, client *resty.Client) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key is empty")
	}
	if client == nil {
		client = resty.New()
	}
	return &Client{apiKey: apiKey, http: client}, nil
}

// InitWithRawClient setups the http Client with an apiKey
func InitWithRawClient(apiKey string, hc *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("api key is empty")
	}
	var client *resty.Client
	if client == nil {
		client = resty.NewWithClient(hc)
	}
	return &Client{apiKey: apiKey, http: client}, nil
}

// SetSessionID will set the session id.
func (c *Client) SetSessionID(sid string) error {
	if sid == "" {
		return errors.New("the session id is empty")
	}
	c.sessionID = sid
	return nil
}

// SetClientAutoRetry sets autoRetry flag to true.
func (c *Client) SetClientAutoRetry() {
	c.autoRetry = true
}

// GetHttpClient get resty client
func (c *Client) GetHttpClient() *resty.Client {
	return c.http
}

// Auto retry default duration.
const defaultRetryDuration = time.Second * 5

// retryDuration calculates the retry duration time.
func retryDuration(resp *http.Response) time.Duration {
	retryTime := resp.Header.Get("Retry-After")
	if retryTime == "" {
		return defaultRetryDuration
	}
	seconds, err := strconv.ParseInt(retryTime, 10, 32)
	if err != nil {
		return defaultRetryDuration
	}
	return time.Duration(seconds) * time.Second
}

// shouldRetry determines whether the status code indicates that the
// previous operation should be retried at a later time.
func shouldRetry(status int) bool {
	return status == http.StatusAccepted || status == http.StatusTooManyRequests
}

func (c *Client) get(url string, data interface{}) error {
	if url == "" {
		return errors.New("url field is empty")
	}

	if c.http.GetClient().Timeout == 0 {
		c.http.SetTimeout(time.Second * 10)
	}

	for {
		resp, err := c.http.R().SetHeaders(map[string]string{
			"content-type": "application/json;charset=utf-8",
		}).SetResult(data).Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode() == http.StatusTooManyRequests && c.autoRetry {
			time.Sleep(retryDuration(resp.RawResponse))
			continue
		}
		if resp.StatusCode() == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode() != http.StatusOK {
			return c.decodeError(resp.RawResponse)
		}
		break
	}
	return nil
}

func (c *Client) request(
	url string,
	body interface{},
	method string,
	data interface{},
) error {
	if url == "" {
		return errors.New("url field is empty")
	}
	if c.http.GetClient().Timeout == 0 {
		c.http.SetTimeout(time.Second * 10)
	}
	for {
		request := c.http.R().SetHeaders(map[string]string{
			"content-type": "application/json;charset=utf-8",
		}).SetBody(body).SetResult(data)
		resp, err := request.Execute(method, url)
		if err != nil {
			return errors.New(err.Error())
		}
		if c.autoRetry && shouldRetry(resp.StatusCode()) {
			time.Sleep(retryDuration(resp.RawResponse))
			continue
		}
		// Checking if the response is greater or equal
		// to 300 or less than 200.
		if resp.StatusCode() >= http.StatusMultipleChoices ||
			resp.StatusCode() < http.StatusOK ||
			resp.StatusCode() == http.StatusNoContent {
			return c.decodeError(resp.RawResponse)
		}
		break
	}
	return nil
}

func (c *Client) fmtOptions(
	urlOptions map[string]string,
) string {
	options := ""
	if len(urlOptions) > 0 {
		for key, value := range urlOptions {
			options += fmt.Sprintf(
				"&%s=%s",
				key,
				url.QueryEscape(value),
			)
		}
	}
	return options
}

// Error type represents an error returned by the TMDB API.
type Error struct {
	StatusMessage string `json:"status_message,omitempty"`
	Success       bool   `json:"success,omitempty"`
	StatusCode    int    `json:"status_code,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"code: %d | success: %t | message: %s",
		e.StatusCode,
		e.Success,
		e.StatusMessage,
	)
}

func (c *Client) decodeError(r *http.Response) error {
	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("could not read body response: %s", err)
	}
	if len(resBody) == 0 {
		return fmt.Errorf(
			"[%d]: empty body %s",
			r.StatusCode,
			http.StatusText(r.StatusCode),
		)
	}
	buf := bytes.NewBuffer(resBody)
	var clientError Error
	if err := json.NewDecoder(buf).Decode(&clientError); err != nil {
		return fmt.Errorf(
			"couldn't decode error: (%d) [%s]",
			len(resBody),
			resBody,
		)
	}
	return clientError
}
