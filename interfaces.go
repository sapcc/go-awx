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
)

// Client represents the client for the AWX API.
type Client interface {
	Ping(ctx context.Context) (output GetPingOutput, err error)
	Reader
	Writer
}

// Object represents an object in AWX.
type Object interface {
}

// ObjectList represents a list of objects in AWX.
type ObjectList interface {
}

// ListOption represents an option for list operation.
type ListOption interface {
}

// GetOption represents an option for get operation.
type GetOption interface {
}

// Reader knows how to read and list AWX objects.
type Reader interface {
	// Get retrieves an obj for the given object key from AWX.
	// obj must be a struct pointer so that obj can be updated with the response
	// returned by the Server.
	Get(ctx context.Context, key ObjectKey, obj Object, httpStatus []int) error

	// List retrieves list of objects for a given resource and list options. On a
	// successful call, Items field in the list will be populated with the
	// result returned from the server.
	List(ctx context.Context, key ObjectKey, list ObjectList, opts ListOption, httpStatus []int) error
}

// Writer knows how to create, delete, and update Kubernetes objects.
type Writer interface {
	// Create saves the object obj in AWX. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Create(ctx context.Context, key ObjectKey, obj Object, httpStatus []int) error

	// Delete deletes the given obj from AWXAWX.
	Delete(ctx context.Context, key ObjectKey, httpStatus []int) error

	// Update updates the given obj in the AWX. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Update(ctx context.Context, key ObjectKey, obj Object, httpStatus []int) error
}
