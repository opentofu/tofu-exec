// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExpressions_Unmarshal(t *testing.T) {
	cases := []struct {
		name     string
		in       string
		expected *Expression
	}{
		{
			name: "constant value",
			in: `
{
  "constant_value": "ami-foobar"
}
`,
			expected: &Expression{
				ExpressionData: &ExpressionData{
					ConstantValue: "ami-foobar",
				},
			},
		},
		{
			name: "null in constant value",
			in: `
{
  "constant_value": {
	"foo": null
  }
}
`,
			expected: &Expression{
				ExpressionData: &ExpressionData{
					ConstantValue: map[string]interface{}{
						"foo": nil,
					},
				},
			},
		},
		{
			name: "references",
			in: `
{
	"references": ["var.in"]
}
`,
			expected: &Expression{
				ExpressionData: &ExpressionData{
					ConstantValue: UnknownConstantValue,
					References:    []string{"var.in"},
				},
			},
		},
		{
			name: "nested blocks",
			in: `
[
	{
		"test_nested_block": {
			"references": ["var.other"]
		}
	}
]
`,
			expected: &Expression{
				ExpressionData: &ExpressionData{
					NestedBlocks: []map[string]*Expression{
						{
							"test_nested_block": {
								ExpressionData: &ExpressionData{
									References:    []string{"var.other"},
									ConstantValue: UnknownConstantValue,
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got *Expression
			if err := json.Unmarshal([]byte(tc.in), &got); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Fatalf("unexpected result (-want,+got):\n%s", diff)
			}
		})
	}
}
