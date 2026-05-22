// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

// VersionOutput represents output from the version -json command
// added in v0.13
type VersionOutput struct {
	Version            string            `json:"terraform_version"`
	Revision           string            `json:"terraform_revision"`
	Platform           string            `json:"platform,omitempty"`
	ProviderSelections map[string]string `json:"provider_selections"`
	Outdated           bool              `json:"terraform_outdated"`
}
