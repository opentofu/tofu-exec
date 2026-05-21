// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

// ResourceMode is a string representation of the resource type found
// in certain fields in the plan.
type ResourceMode string

const (
	// DataResourceMode is the resource mode for data sources.
	DataResourceMode ResourceMode = "data"

	// ManagedResourceMode is the resource mode for managed resources.
	ManagedResourceMode ResourceMode = "managed"

	// EphemeralResourceMode is the resource mode for ephemeral resources.
	EphemeralResourceMode ResourceMode = "ephemeral"
)
