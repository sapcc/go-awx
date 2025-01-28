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

func TestJobTemplates(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /job_templates", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer 12345" {
			http.Error(w, "Auth header was incorrect", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"results": [
				{
					"id": 1,
					"name": "Test Job",
					"url": "/api/v2/job_templates/1/",
					"type": "job",
					"modified": "2025-01-01T00:00:00Z",
					"created": "2025-01-01T00:00:00Z",
					"status": "successful"
				}
			]
		}`))
		assert.NoError(t, err)
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
	result := JobTemplateList{}
	err = client.List(context.Background(), ObjectKey{
		Resource: "job_templates",
	}, &result, input, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Results))
	assert.Equal(t, 1, result.Results[0].ID)
	assert.Equal(t, "Test Job", result.Results[0].Name)
}

func TestCreateJobTemplatesSchedule(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /job_templates/1/schedules/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer 12345" {
			http.Error(w, "Auth header was incorrect", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{
			"id": 1,
			"name": "Test Schedule",
			"rrule": "FREQ=DAILY;COUNT=1",
			"unified_job_template": 1,
			"extra_data": "",
			"inventory": 1,
			"limit": ""
		}`))
		assert.NoError(t, err)
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
	result := Schedule{
		Name:  "Test Schedule",
		RRULE: "FREQ=DAILY;COUNT=1",
	}
	err = client.Create(context.Background(), ObjectKey{
		Resource:   "job_templates",
		ResourceID: "1",
		Action:     "schedules",
	}, &result, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Test Schedule", result.Name)
}
