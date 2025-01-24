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
	"github.com/gorilla/schema"
)

var encode = schema.NewEncoder()

// InventoryList represents the output of the ListInventories method.
type InventoryList struct {
	ListGetResponse
	Results []*Inventory `json:"results,omitempty"`
}

// Inventory represents the output of the GetInventories method.
type Inventory struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	URL               string `json:"url"`
	Type              string `json:"type"`
	Modified          string `json:"modified"`
	Created           string `json:"created"`
	HasActiveFailures bool   `json:"has_active_failures"`
	TotalHosts        int    `json:"total_hosts"`
	PendingDeletion   bool   `json:"pending_deletion"`
}

// InventoryListInput represents the input of the ListInventories method.
type InventoryListInput struct {
	Query string `schema:"-,omitempty"`
	Name  string `schema:"name,omitempty"`
}
