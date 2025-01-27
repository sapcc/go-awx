/* ******************************************************************************
*
* Copyright 2025 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/
package awx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateSchedule(t *testing.T) {
	schedule := Schedule{
		Name:               "Test Schedule",
		RRULE:              "FREQ=DAILY;INTERVAL=1",
		UnifiedJobTemplate: 1,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/schedules", r.URL.Path)

		var received Schedule
		err := json.NewDecoder(r.Body).Decode(&received)
		assert.NoError(t, err)
		assert.Equal(t, schedule, received)

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(Schedule{
			ID:          1,
			Name:        schedule.Name,
			RRULE:       schedule.RRULE,
			Description: "Test Schedule",
		})
		assert.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(ClientOptions{
		Endpoint: server.URL + "/",
		Token:    "12345",
	},
	)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Create(ctx, ObjectKey{
		Resource: "schedules",
	}, schedule, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, schedule.ID)
	assert.Equal(t, schedule.Name, schedule.Name)
	assert.Equal(t, schedule.RRULE, schedule.RRULE)
}

func TestClient_DeleteSchedule(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/schedules/1", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client, err := NewClient(ClientOptions{
		Endpoint: server.URL + "/",
		Token:    "12345",
	},
	)
	assert.NoError(t, err)

	ctx := context.Background()
	err = client.Delete(ctx, ObjectKey{
		Resource:   "schedules",
		ResourceID: "1",
	}, nil)
	assert.NoError(t, err)
}
