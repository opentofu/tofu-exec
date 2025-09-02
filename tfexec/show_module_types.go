// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// This file is essentially a copy of the following in the opentofu core codebase.
// https://github.com/opentofu/opentofu/blob/28493bc63f83aaa5fb2ff5063f050d80f9c51f4f/internal/command/jsonconfig/config.go
// 2 things are modified.
// - Types are exported, since this is intended to be public API.
// - And expressions are removed, since those aren't marshaled during module config generation

package tfexec

import "encoding/json"

type ModuleRoot struct {
	Module Module `json:"root_module"`
}
type Module struct {
	Outputs map[string]Output `json:"outputs,omitempty"`
	// Resources are sorted in a user-friendly order that is undefined at this
	// time, but consistent.
	Resources   []Resource            `json:"resources,omitempty"`
	ModuleCalls map[string]ModuleCall `json:"module_calls,omitempty"`
	Variables   Variables             `json:"variables,omitempty"`
}

type ModuleCall struct {
	Source            string   `json:"source,omitempty"`
	Module            *Module  `json:"module,omitempty"`
	VersionConstraint string   `json:"version_constraint,omitempty"`
	DependsOn         []string `json:"depends_on,omitempty"`
}

// Variables is the json representation of the Variables provided to the current
// plan.
type Variables map[string]*Variable

type Variable struct {
	Type        json.RawMessage `json:"type,omitempty"`
	Default     json.RawMessage `json:"default,omitempty"`
	Description string          `json:"description,omitempty"`
	Required    bool            `json:"required,omitempty"`
	Sensitive   bool            `json:"sensitive,omitempty"`
	Deprecated  string          `json:"deprecated,omitempty"`
}

// Resource is the representation of a Resource in the config
type Resource struct {
	// Address is the absolute resource address
	Address string `json:"address,omitempty"`

	// Mode can be "managed" or "data"
	Mode string `json:"mode,omitempty"`

	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`

	// ProviderConfigKey is the key into "provider_configs" (shown above) for
	// the provider configuration that this resource is associated with.
	//
	// NOTE: If a given resource is in a ModuleCall, and the provider was
	// configured outside of the module (in a higher level configuration file),
	// the ProviderConfigKey will not match a key in the ProviderConfigs map.
	ProviderConfigKey string `json:"provider_config_key,omitempty"`

	// Provisioners is an optional field which describes any provisioners.
	// Connection info will not be included here.
	Provisioners []Provisioner `json:"provisioners,omitempty"`

	// Expressions" describes the resource-type-specific  content of the
	// configuration block.
	Expressions map[string]any `json:"expressions,omitempty"`

	// SchemaVersion indicates which version of the resource type schema the
	// "values" property conforms to.
	SchemaVersion *uint64 `json:"schema_version,omitempty"`

	DependsOn []string `json:"depends_on,omitempty"`
}

type Provisioner struct {
	Type        string         `json:"type,omitempty"`
	Expressions map[string]any `json:"expressions,omitempty"`
}

type Output struct {
	Sensitive   bool     `json:"sensitive,omitempty"`
	Deprecated  string   `json:"deprecated,omitempty"`
	DependsOn   []string `json:"depends_on,omitempty"`
	Description string   `json:"description,omitempty"`
}
