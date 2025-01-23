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

// Schedule represents the output of the GetSchedule method.
type Schedule struct {
	ID                 int    `json:"id"`
	RRULE              string `json:"rrule"`
	Type               string `json:"type"`
	URL                string `json:"url"`
	Created            string `json:"created"`
	Modified           string `json:"modified"`
	ExtraData          string `json:"extra_data"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	UnifiedJobTemplate int    `json:"unified_job_template"`
	Inventory          int    `json:"inventory"`
}

// ScheduleList represents the output of the ListSchedule method.
type ScheduleList struct {
	ListGetResponse
	Results []*Schedule `json:"results,omitempty"`
}

// ListSchedulesInput represents the input of the ListSchedules method.
type ListSchedulesInput struct {
	ID   string `schema:"id,omitempty"`
	Name string `schema:"name,omitempty"`
}
