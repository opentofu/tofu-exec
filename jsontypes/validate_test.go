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

func TestValidate_Validate(t *testing.T) {
	t.Run("nil function", func(t *testing.T) {
		var vo *ValidateOutput
		if err := vo.Validate(); err == nil || err.Error() != "validation output is nil" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("no format version", func(t *testing.T) {
		vo := &ValidateOutput{}
		if err := vo.Validate(); err != nil {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("invalid format version", func(t *testing.T) {
		vo := &ValidateOutput{
			FormatVersion: "invalid version constraint",
		}
		if err := vo.Validate(); err == nil || err.Error() != `invalid format version "invalid version constraint": Malformed version: invalid version constraint` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("format version does not satisfy constraint", func(t *testing.T) {
		vo := &ValidateOutput{
			FormatVersion: "0.0.2",
		}
		if err := vo.Validate(); err == nil || err.Error() != `unsupported validation output format version: "0.0.2" does not satisfy ">= 0.1, < 2.0"` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		vo := &ValidateOutput{
			FormatVersion: "1.0.0",
		}
		if err := vo.Validate(); err != nil {
			t.Errorf("invalid error: %s", err)
		}
	})
}

func TestValidate_Unmarshal(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		in := `{
	"valid": true,
	"error_count": 1,
	"warning_count": 1,
	"diagnostics": [
		{
			"severity": "error",
			"summary": "error diagnostic",
			"detail": "details of the error diagnostic"
		},
		{
			"severity": "warning",
			"summary": "warning diagnostic",
			"detail": "details of the warning diagnostic"
		},
		{
			"severity": "unknown",
			"summary": "unknown diagnostic",
			"detail": "details of the unknown diagnostic"
		}
	]
}`
		want := &ValidateOutput{
			Valid:        true,
			ErrorCount:   1,
			WarningCount: 1,
			Diagnostics: []Diagnostic{
				{
					Severity: "error",
					Summary:  "error diagnostic",
					Detail:   "details of the error diagnostic",
				},
				{
					Severity: "warning",
					Summary:  "warning diagnostic",
					Detail:   "details of the warning diagnostic",
				},
				{
					Severity: "unknown",
					Summary:  "unknown diagnostic",
					Detail:   "details of the unknown diagnostic",
				},
			},
		}
		var got *ValidateOutput
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		opts := cmpopts.IgnoreUnexported(cty.Type{})
		if diff := cmp.Diff(want, got, opts); diff != "" {
			t.Fatalf("unexpected result (-want,+got):\n%s", diff)
		}
	})
	t.Run("invalid input", func(t *testing.T) {
		in := `invalid input`
		var got *ValidateOutput
		if err := json.Unmarshal([]byte(in), &got); err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
			t.Fatalf("unexpected error: %s", err)
		}
	})
}
