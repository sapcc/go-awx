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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListJobs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer 12345" {
			http.Error(w, "Auth header was incorrect", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"results": [
				{
					"id": 1,
					"name": "Test Job",
					"url": "/api/v2/jobs/1/",
					"type": "job",
					"modified": "2025-01-01T00:00:00Z",
					"created": "2025-01-01T00:00:00Z",
					"unified_job_template": 1,
					"status": "successful",
					"failed": false,
					"started": "2025-01-01T00:00:00Z",
					"finished": "2025-01-01T00:00:00Z"
				}
			]
		}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	client, err := NewClient(
		ClientOptions{
			Endpoint:           server.URL + "/",
			Token:              "12345",
			InsecureSkipVerify: true,
		},
	)
	assert.NoError(t, err)
	input := ListJobsInput{ID: "1"}
	result := JobList{}
	err = client.List(context.Background(), ObjectKey{
		Resource: "jobs",
	}, &result, input, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, 1, result.Results[0].ID)
}

func TestGetJobs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs/1", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer 12345" {
			http.Error(w, "Auth header was incorrect", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{

			"id": 1,
			"name": "Test Job",
			"url": "/api/v2/jobs/1/",
			"type": "job",
			"modified": "2025-01-01T00:00:00Z",
			"created": "2025-01-01T00:00:00Z",
			"unified_job_template": 1,
			"status": "successful",
			"failed": false,
			"started": "2025-01-01T00:00:00Z",
			"finished": "2025-01-01T00:00:00Z"

		}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClient(
		ClientOptions{
			Endpoint:           server.URL + "/",
			Token:              "12345",
			InsecureSkipVerify: true,
		},
	)
	assert.NoError(t, err)

	result := Job{}
	err = client.Get(context.Background(), ObjectKey{
		Resource:   "jobs",
		ResourceID: "1",
	}, &result, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test Job", result.Name)
}

func TestCanCancelJob(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs/1/cancel", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		if auth !=
			"Bearer 12345" {
			http.Error(w, "Auth header was incorrect", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"can_cancel": true
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	client, err := NewClient(
		ClientOptions{
			Endpoint:           server.URL + "/",
			Token:              "12345",
			InsecureSkipVerify: true,
		},
	)
	assert.NoError(t, err)
	result := CanCancelJob{}
	err = client.Get(context.Background(), ObjectKey{
		Resource:   "jobs",
		ResourceID: "1",
		Action:     "cancel",
	}, &result, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, result.CanCancel)
}
