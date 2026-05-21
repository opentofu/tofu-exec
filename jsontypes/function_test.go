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

func TestFunctions_Validate(t *testing.T) {
	t.Run("nil function", func(t *testing.T) {
		var f *Functions
		if err := f.Validate(); err == nil || err.Error() != "metadata functions data is nil" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("no format version", func(t *testing.T) {
		f := &Functions{}
		if err := f.Validate(); err == nil || err.Error() != "unexpected metadata functions data, format version is missing" {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("invalid format version", func(t *testing.T) {
		f := &Functions{
			FormatVersion: "invalid version constraint",
		}
		if err := f.Validate(); err == nil || err.Error() != `invalid format version "invalid version constraint": Malformed version: invalid version constraint` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("format version does not satisfy constraint", func(t *testing.T) {
		f := &Functions{
			FormatVersion: "0.0.2",
		}
		if err := f.Validate(); err == nil || err.Error() != `unsupported metadata functions format version: "0.0.2" does not satisfy "~> 1.0"` {
			t.Errorf("invalid error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		f := &Functions{
			FormatVersion: "1.0.0",
		}
		if err := f.Validate(); err != nil {
			t.Errorf("invalid error: %s", err)
		}
	})
}

func TestFunctions_Unmarshal(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected *Functions
	}{
		{
			name: "multiple functions",
			in: `
{
	"format_version": "1.1.0",
	"function_signatures": {
		"empty": {
			"description": "This is the empty func description",
			"summary": "checks if string is empty",
			"deprecation_message": "deprecated. use empty_unicode instead",
			"return_type": "bool",
			"parameters": [
				{"name": "in", "type": "string"}
			]
		},
		"join": {
			"description": "Joins all the elements by using the separator",
			"return_type": "string",
			"parameters": [
				{"name": "separator", "type": "string"}
			],
			"variadic_parameter": {
				"name": "parts",
				"type": "string"
			}
		}
	}
}
`,
			expected: &Functions{
				FormatVersion: "1.1.0",
				Signatures: map[string]*FunctionSignature{
					"empty": {
						Description:        "This is the empty func description",
						Summary:            "checks if string is empty",
						DeprecationMessage: "deprecated. use empty_unicode instead",
						ReturnType:         cty.Bool,
						Parameters: []*FunctionParameter{
							{
								Name: "in",
								Type: cty.String,
							},
						},
					},
					"join": {
						Description: "Joins all the elements by using the separator",
						ReturnType:  cty.String,
						Parameters: []*FunctionParameter{
							{
								Name: "separator",
								Type: cty.String,
							},
						},
						VariadicParameter: &FunctionParameter{
							Name: "parts",
							Type: cty.String,
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got *Functions
			if err := json.Unmarshal([]byte(tc.in), &got); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			opts := cmpopts.IgnoreUnexported(cty.Type{})
			if diff := cmp.Diff(tc.expected, got, opts); diff != "" {
				t.Fatalf("unexpected result (-want,+got):\n%s", diff)
			}
		})
	}
}
