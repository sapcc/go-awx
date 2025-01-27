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

// JobTemplateList represents the output of the ListJobTemplates method.
type JobTemplateList struct {
	ListGetResponse
	Results []*JobTemplate `json:"results,omitempty"`
}

// JobTemplate represents the output of the GetJobTemplates method.
type JobTemplate struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Type     string `json:"type"`
	Modified string `json:"modified"`
	Created  string `json:"created"`
	Status   string `json:"status"`
}

// ListJobTemplateInput represents the input of the ListJobTemplates method.
type ListJobTemplateInput struct {
	ID    string `schema:"id,omitempty"`
	Name  string `schema:"name,omitempty"`
	Query string `schema:"-,omitempty"`
}

// CreateJobTemplatesSchedule represents the input of the PostJobTemplateSchedule method.
type CreateJobTemplatesSchedule struct {
	Name               string `json:"name"`
	RRULE              string `json:"rrule"`
	UnifiedJobTemplate int    `json:"unified_job_template"`
	ExtraData          string `json:"extra_data"`
	Inventory          string `json:"inventory"`
	Limit              string `json:"limit"`
}
