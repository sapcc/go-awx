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

import "time"

// JobList represents the output of the ListJobs method.
type JobList struct {
	ListGetResponse
	Results []*Job `json:"results,omitempty"`
}

// Job represents the output of the GetJobs method.
type Job struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	URL                string    `json:"url"`
	Type               string    `json:"type"`
	Modified           string    `json:"modified"`
	Created            string    `json:"created"`
	UnifiedJobTemplate int       `json:"unified_job_template"`
	Status             string    `json:"status"`
	Failed             bool      `json:"failed"`
	Started            time.Time `json:"started"`
	Finished           time.Time `json:"finished"`
}

// ListJobsInput represents the input of the ListJobs method.
type ListJobsInput struct {
	ID         string `schema:"id,omitempty"`
	LaunchType string `schema:"launch_type,omitempty"`
	ScheduleID int    `schema:"schedule__id,omitempty"`
	Query      string `schema:"query,omitempty"`
}

// CanCancelJob represents the output of the GetCancelJob method.
type CanCancelJob struct {
	CanCancel bool `json:"can_cancel"`
}
