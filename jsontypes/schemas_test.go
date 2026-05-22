// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zclconf/go-cty/cty"
)

func TestProviderSchema_Validate(t *testing.T) {
	t.Run("nil schemas", func(t *testing.T) {
		var ps *ProviderSchemas
		if err := ps.Validate(); err == nil || err.Error() != "provider schema data is nil" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("no format version", func(t *testing.T) {
		ps := &ProviderSchemas{}
		if err := ps.Validate(); err == nil || err.Error() != "unexpected provider schema data, format version is missing" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("invalid format version", func(t *testing.T) {
		ps := &ProviderSchemas{
			FormatVersion: "invalid version constraint",
		}
		if err := ps.Validate(); err == nil || err.Error() != `invalid format version "invalid version constraint": Malformed version: invalid version constraint` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("format version does not satisfy constraint", func(t *testing.T) {
		ps := &ProviderSchemas{
			FormatVersion: "0.0.2",
		}
		if err := ps.Validate(); err == nil || err.Error() != `unsupported provider schema format version: "0.0.2" does not satisfy ">= 0.1, < 2.0"` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		ps := &ProviderSchemas{
			FormatVersion: "0.1.0",
		}
		if err := ps.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestProviderSchema_Unmarshal(t *testing.T) {
	resSchema := Schema{
		Version: 1,
		Block: &SchemaBlock{
			Attributes: map[string]*SchemaAttribute{
				"name": {
					Optional: true,
				},
			},
		},
	}
	t.Run("success", func(t *testing.T) {
		in := `{
	"format_version": "1.1.0",
	"provider_schemas": {
		"provider_name": {
			"resource_schemas": {
				"res1": {
					"version": 1,
					"block": {
						"attributes": {
							"name": {
								"optional": true
							}
						}
					}
				}
			},
			"data_source_schemas": {
				"data1": {
					"version": 1,
					"block": {
						"attributes": {
							"name": {
								"optional": true
							}
						}
					}
				}
			},
			"ephemeral_resource_schemas": {
				"eph1": {
					"version": 1,
					"block": {
						"attributes": {
							"name": {
								"optional": true
							}
						}
					}
				}
			}
		}
	}
}`
		want := &ProviderSchemas{
			FormatVersion: "1.1.0",
			Schemas: map[string]*ProviderSchema{
				"provider_name": {
					ResourceSchemas: map[string]*Schema{
						"res1": &resSchema,
					},
					DataSourceSchemas: map[string]*Schema{
						"data1": &resSchema,
					},
					EphemeralResourceSchemas: map[string]*Schema{
						"eph1": &resSchema,
					},
				},
			},
		}
		var got *ProviderSchemas
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{})
		if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		ps := &ProviderSchemas{}
		if err := json.Unmarshal([]byte("invalid input"), &ps); err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
			t.Errorf("invalid error: %s", err)
		}
	})
}
