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
	"strconv"
	"strings"
)

// ListGetResponse represents the output of the ListGet method.
type ListGetResponse struct {
	Count    int    `json:"count,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

// ErrorResponse represents the awx error response.
type ErrorResponse struct {
	All []string `json:"__all__"`
}

// Error represents the output of the Error method.
type Error struct {
	StatusCode int
	Msg        []string
}

// Error returns the error message.
func (e *Error) Error() string {
	return "status: " + strconv.Itoa(e.StatusCode) + ", messages: " + strings.Join(e.Msg, ", ")
}

// NewError creates a new error.
func NewError(statusCode int, msg []string) *Error {
	return &Error{statusCode, msg}
}

// ObjectKey represents the key of an object.
type ObjectKey struct {
	Resource   string
	ResourceID string
	Action     string
}

func (k ObjectKey) String() string {
	path := k.Resource
	if k.ResourceID != "" {
		path = path + "/" + k.ResourceID
	}
	if k.Action != "" {
		path = path + "/" + k.Action + "/"
	}
	return path
}
