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
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockHTTPClient is a mock implementation of the HTTPClient interface.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	options := ClientOptions{
		Endpoint: "http://example.com",
	}

	client, err := NewClient(options)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestPing(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"version":"1.0","active_node":"node1","install_uuid":"uuid"}`)),
			}, nil
		},
	}
	client, err := NewClient(ClientOptions{
		Endpoint:   "http://example.com",
		HTTPClient: mockClient,
		Token:      "test-token",
		Agent:      "test-agent",
	})
	assert.NoError(t, err)
	ctx := context.Background()
	output, err := client.Ping(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "1.0", output.Version)
	assert.Equal(t, "node1", output.ActiveNode)
	assert.Equal(t, "uuid", output.InstallUUID)
}
