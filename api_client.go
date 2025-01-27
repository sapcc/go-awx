/******************************************************************************
*
*  Copyright 2025 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

package awx

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents the client for the AWX API.
type client struct {
	parsedURL  *url.URL
	httpClient HTTPClient
	token      string
	agent      string
	username   string
	password   string
	version    string
}

// ClientOptions represents the options for the client.
type ClientOptions struct {
	Endpoint           string
	HTTPClient         HTTPClient
	InsecureSkipVerify bool
	Username           string
	Password           string
	Token              string
	Agent              string
	Version            string
}

// GetAuthTokenInput represents the input of the GetAuthToken method.
type GetAuthTokenInput struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// GetAuthTokenOutput represents the output of the GetAuthToken method.
type GetAuthTokenOutput struct {
	Token   string    `json:"token,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

// HTTPClient provides the interface for a client making HTTP requests with the
// API.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient creates a new client for the AWX API.
func NewClient(options ClientOptions) (c Client, err error) {
	// Check the URL:
	if options.Endpoint == "" {
		return c, errors.New("the baseURL is mandatory")
	}
	iCl := &client{}
	iCl.parsedURL, err = url.Parse(options.Endpoint)
	if err != nil {
		err = fmt.Errorf("the URL '%s' isn't valid: %s", options.Endpoint, err.Error())
		return nil, err
	}
	if options.HTTPClient == nil {
		options.HTTPClient = createClient(options)
	}
	if options.Agent == "" {
		options.Agent = "go-awx-client"
	}
	iCl.version = options.Version
	iCl.agent = options.Agent
	iCl.httpClient = options.HTTPClient
	iCl.token = options.Token
	iCl.password = options.Password
	iCl.username = options.Username
	return iCl, nil
}

func createClient(options ClientOptions) (client *http.Client) {
	// Create the HTTP client:
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: options.InsecureSkipVerify, // #nosec G402
			},
		},
	}
	return
}

// DoRequest performs an HTTP request to the AWX API.
func (c *client) DoRequest(req *http.Request, okCodes []int) ([]byte, error) {
	if c.httpClient == nil {
		return nil, errors.New("the HTTP client is mandatory")
	}
	if c.parsedURL == nil {
		return nil, errors.New("the URL is mandatory")
	}
	if c.token == "" {
		if err := c.getAuthToken(); err != nil {
			return nil, err
		}
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer " + c.token},
		"User-Agent":    []string{c.agent},
		"Content-Type":  []string{"application/json"},
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	// Validate the HTTP response status.
	ok := false
	for _, code := range okCodes {
		if res.StatusCode == code {
			ok = true
			break
		}
	}
	if !ok {
		output := ErrorResponse{}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected error response: %s", err.Error())
		}
		dec := json.NewDecoder(bytes.NewReader(body)).Decode(&output)
		if dec != nil {
			return nil, NewError(res.StatusCode, []string{
				"unexpected error response: " + string(body),
			})
		}
		if output.All == nil {
			return nil, NewError(res.StatusCode, []string{
				"unexpected error response: " + string(body),
			})
		}
		return nil, NewError(res.StatusCode, output.All)
	}
	body, err := io.ReadAll(res.Body)
	return body, err
}

// GetPingOutput represents the output of the Ping method.
type GetPingOutput struct {
	Version     string `json:"version"`
	ActiveNode  string `json:"active_node"`
	InstallUUID string `json:"install_uuid"`
}

// Ping retrieves awx information.
func (c *client) Ping(ctx context.Context) (output GetPingOutput, err error) {
	if c.parsedURL == nil {
		return output, errors.New("the URL is mandatory")
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    c.parsedURL.JoinPath("ping/"),
	}

	body, err := c.DoRequest(&req, []int{200})
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &output); err != nil {
		return
	}
	return
}

func (c *client) Get(ctx context.Context, key ObjectKey, output Object, httpStatus []int) error {
	req := http.Request{
		Method: http.MethodGet,
		URL:    c.parsedURL.JoinPath(key.String()),
	}
	fmt.Println(req.URL)
	if httpStatus == nil {
		httpStatus = []int{http.StatusOK}
	}
	body, err := c.DoRequest(&req, httpStatus)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, output)
	return err
}

func (c *client) Create(ctx context.Context, key ObjectKey, obj Object, status []int) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(obj); err != nil {
		return err
	}
	req := http.Request{
		Method: http.MethodPost,
		URL:    c.parsedURL.JoinPath(key.String()),
		Body:   io.NopCloser(&buf),
	}
	fmt.Println(req.URL)
	if len(status) == 0 {
		status = []int{http.StatusCreated}
	}
	body, err := c.DoRequest(&req, status)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &obj)
	return err
}

func (c *client) List(ctx context.Context, key ObjectKey, obj ObjectList, options ListOption, httpStatus []int) error {
	values := url.Values{}
	err := encode.Encode(options, values)
	if err != nil {
		return err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    c.parsedURL.JoinPath(key.String()),
	}
	req.Host = c.parsedURL.Host
	req.URL.RawQuery = values.Encode()

	if httpStatus == nil {
		httpStatus = []int{http.StatusOK}
	}
	body, err := c.DoRequest(&req, httpStatus)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &obj)
}

func (c *client) Update(ctx context.Context, key ObjectKey, obj Object, httpStatus []int) error {
	return nil
}

func (c *client) Delete(ctx context.Context, key ObjectKey, httpStatus []int) error {
	req := http.Request{
		Method: http.MethodDelete,
		URL:    c.parsedURL.JoinPath(key.String()),
	}
	if httpStatus == nil {
		httpStatus = []int{http.StatusNoContent}
	}
	_, err := c.DoRequest(&req, httpStatus)
	return err
}

func (c *client) Patch(ctx context.Context, path string, obj Object) error {
	return nil
}

func (c *client) getAuthToken() error {
	var input GetAuthTokenInput
	var output GetAuthTokenOutput
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(input); err != nil {
		return err
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.parsedURL.JoinPath("tokens/").String(),
		http.NoBody,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error obtaining auth token: %s", string(body))
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	err = json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return fmt.Errorf("error obtaining auth token: %s", string(body))
	}
	if output.Token == "" {
		return fmt.Errorf("Error obtaining auth token: %s", string(body))
	}
	c.token = output.Token
	return nil
}
