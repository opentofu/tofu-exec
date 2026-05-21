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

func TestState_Validate(t *testing.T) {
	t.Run("nil state", func(t *testing.T) {
		var s *State
		if err := s.Validate(); err == nil || err.Error() != "state is nil" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("no format version", func(t *testing.T) {
		s := &State{}
		if err := s.Validate(); err == nil || err.Error() != "unexpected state input, format version is missing" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("invalid format version", func(t *testing.T) {
		s := &State{
			FormatVersion: "invalid version constraint",
		}
		if err := s.Validate(); err == nil || err.Error() != `invalid format version "invalid version constraint": Malformed version: invalid version constraint` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("format version does not satisfy constraint", func(t *testing.T) {
		s := &State{
			FormatVersion: "0.0.2",
		}
		if err := s.Validate(); err == nil || err.Error() != `unsupported state format version: "0.0.2" does not satisfy ">= 0.1, < 2.0"` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		s := &State{
			FormatVersion: "0.1.0",
		}
		if err := s.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestState_Unmarshal(t *testing.T) {
	in := `{
	"format_version": "1.1.0",
	"terraform_version": "1.12.0",
	"values": {
		"root_module": {
			"resources": [
				{
					"address": "test_resource.name",
					"mode": "managed",
					"type": "test_resource",
					"name": "name",
					"provider_name": "test",
					"values": {
						"name": "test_name",
						"val": 1
					}
				}
			]
		},
		"outputs": {
			"out": {
				"type": "string",
				"value": "test"
			}
		}
	}
}`
	t.Run("without decoder number usage", func(t *testing.T) {
		p := &State{}
		if err := json.Unmarshal([]byte(in), &p); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		want := &State{
			FormatVersion:    "1.1.0",
			TerraformVersion: "1.12.0",
			Values: &StateValues{
				RootModule: &StateModule{
					Resources: []*StateResource{
						{
							Address:      "test_resource.name",
							Mode:         ManagedResourceMode,
							Type:         "test_resource",
							Name:         "name",
							ProviderName: "test",
							AttributeValues: map[string]any{
								"name": "test_name",
								"val":  float64(1),
							},
						},
					},
				},
				Outputs: map[string]*StateOutput{
					"out": {
						Type:  cty.String,
						Value: "test",
					},
				},
			},
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{}, State{})
		if diff := cmp.Diff(want, p, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})
	t.Run("with decoder number usage", func(t *testing.T) {
		s := &State{}
		s.UseJSONNumber(true)
		if err := json.Unmarshal([]byte(in), &s); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		want := &State{
			FormatVersion:    "1.1.0",
			TerraformVersion: "1.12.0",
			Values: &StateValues{
				RootModule: &StateModule{
					Resources: []*StateResource{
						{
							Address:      "test_resource.name",
							Mode:         ManagedResourceMode,
							Type:         "test_resource",
							Name:         "name",
							ProviderName: "test",
							AttributeValues: map[string]any{
								"name": "test_name",
								"val":  json.Number("1"),
							},
						},
					},
				},
				Outputs: map[string]*StateOutput{
					"out": {
						Type:  cty.String,
						Value: "test",
					},
				},
			},
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{}, State{})
		if diff := cmp.Diff(want, s, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		s := &State{}
		if err := json.Unmarshal([]byte("invalid input"), &s); err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
			t.Errorf("invalid error: %s", err)
		}
	})
}
