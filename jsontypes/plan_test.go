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

func TestPlan_Validate(t *testing.T) {
	t.Run("nil plan", func(t *testing.T) {
		var p *Plan
		if err := p.Validate(); err == nil || err.Error() != "plan is nil" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("no format version", func(t *testing.T) {
		p := &Plan{}
		if err := p.Validate(); err == nil || err.Error() != "unexpected plan input, format version is missing" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("invalid format version", func(t *testing.T) {
		p := &Plan{
			FormatVersion: "invalid version constraint",
		}
		if err := p.Validate(); err == nil || err.Error() != `invalid format version "invalid version constraint": Malformed version: invalid version constraint` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("format version does not satisfy constraint", func(t *testing.T) {
		p := &Plan{
			FormatVersion: "0.0.2",
		}
		if err := p.Validate(); err == nil || err.Error() != `unsupported plan format version: "0.0.2" does not satisfy ">= 0.1, < 2.0"` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		p := &Plan{
			FormatVersion: "0.1.0",
		}
		if err := p.Validate(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}

func TestPlan_Unmarshal(t *testing.T) {
	in := `{
	"format_version": "1.1.0",
	"terraform_version": "1.12.0",
	"configuration": {
		"root_module": {
			"variables": {
				"in": {
					"default": 1.2
				}
			}
		}
	}
}`
	t.Run("without decoder number usage", func(t *testing.T) {
		p := &Plan{}
		if err := json.Unmarshal([]byte(in), &p); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		want := &Plan{
			FormatVersion:    "1.1.0",
			TerraformVersion: "1.12.0",
			Config: &Config{
				RootModule: &ConfigModule{
					Variables: map[string]*ConfigVariable{
						"in": {
							Default: 1.2,
						},
					},
				},
			},
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{}, Plan{})
		if diff := cmp.Diff(want, p, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})
	t.Run("with decoder number usage", func(t *testing.T) {
		p := &Plan{}
		p.UseJSONNumber(true)
		if err := json.Unmarshal([]byte(in), &p); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		want := &Plan{
			FormatVersion:    "1.1.0",
			TerraformVersion: "1.12.0",
			Config: &Config{
				RootModule: &ConfigModule{
					Variables: map[string]*ConfigVariable{
						"in": {
							Default: 1.2,
						},
					},
				},
			},
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{}, Plan{})
		if diff := cmp.Diff(want, p, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		p := &Plan{}
		if err := json.Unmarshal([]byte("invalid input"), &p); err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
			t.Errorf("invalid error: %s", err)
		}
	})
}
